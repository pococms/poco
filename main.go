package main

import (
	"bufio"
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	ytembed "github.com/13rac1/goldmark-embed"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

// IMPORTANT: This is the same name as used in the go:embed directive
const pocoDir = ".poco"

// This directory gets embedded into the executable. It's
// then copied into every new project.
//
//go:embed .poco
var pocoFiles embed.FS

// If javascript files are included, need this to avoid
// starting before doc has loaded.
// It's saved as "loading.js"
var loading = `
if (document.readyState !== 'loading') {
    docReady();
} else {
    document.addEventListener('DOMContentLoaded', function () {
        docReady();
    });
}
`

// Required begininng for a valid HTML document
var docType = `<!DOCTYPE html>
<html lang=`

// If a page lacks a title tag, it fails validation.
// Insert this if none is found.
var poweredBy = `Powered by PocoCMS`

// Adds Javascript after the body, just before the closing </body> tag
func (c *config) scriptAfter() string {
	// If javascript files are included, they should be
	// called from inside this function.
	// NOTE: Make sure the final } gets inserted
	// before the closing </code> tag

	slice := fmStrSlice("script-after", c.fm)
	if slice == nil {
		return ""
	}
	// Return value
	scripts := "<script>\n" + "\t" + "function docReady() {\n\t"
	var s string

	for _, value := range slice {
		filename := value
		s = s + c.getWebOrLocalFileStr(filename)
	}
	scripts += s + "}\n" + "</script>" + "\n"
	return scripts

}

// assemble takes the raw converted HTML in article,
// uses it to generate finished HTML document, and returns
// that document as a string.
func (c *config) assemble(filename string, article string) string {
	// This will contain the completed document as a string.
	htmlFile := ""
	// Execute templates. That way {{ .Title }} will be converted into
	// whatever frontMatter["title"] is set to, etc.
	if parsedArticle, err := doTemplate("", article, c); err != nil {
		quit(1, err, c, "%v: template error", filename)
	} else {
		article = parsedArticle
	}

	// If it's the home page, and a timestamp was requested,
	// insert it in a paragraph at the top of the article.
	timestamp := ""
	if c.timestamp && c.currentFilename == c.homePage {
		timestamp = "\n<p>" + theTime() + "</p>\n"
	}
	// Get Javascript that goes after the body
	scriptAfterStr := c.scriptAfter()
	if scriptAfterStr != "" {
		//hasScript = true
	}
	// Build the completed HTML document from the component pieces.
	htmlFile = docType + "\"" + c.lang + "\">" + "\n" +
		"<head>" +
		"\n\t<meta charset=\"utf-8\">" +
		"\n\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
		c.titleTag() +
		c.metatags() +
		c.linktags() +
		c.stylesheets() +
		c.styleTags() +
		"</head>\n<body>" +
		"\n<div id=\"page-container\">" +
		"\n<div id=\"content-wrap\">\n" +
		"\t" + c.header() +
		"\n\t" + c.nav() +
		"\n\t" + c.aside() +
		"\n<article id=\"article-poco\">\n" + timestamp + article + "\n" + "</article>" + "\n" +
		"</div><!-- content-wrap -->\n" +
		c.footer() + "\n" +
		"</div><!-- page-container -->\n" +
		"<script> {" + "\n" +
		c.documentReady() +
		"}\n</script>" + "\n" +
		"</body>\n</html>\n"
	return htmlFile
} //   assemble

// THEME

// getFm takes Markdown filename passed in, opens that file,
// and returns its front matter. The way Goldmark works it
// means it was necessary to parse and convert the Markdown
// too, but that just gets discarded.
// c.fm is left untouched.
func (c *config) getFm(filename string) map[string]interface{} {
	var rawHTML string
	var err error
	newC := newConfig()

	// Convert Markdown file, possibly with front matter, to HTML
	if rawHTML, err = mdYAMLFileToHTMLString(newC, filename); err != nil {
		quit(1, err, c, "Problem converting %s to HTML", filename)
	}

	// Execute its templates.
	if _, err = doTemplate("", rawHTML, newC); err != nil {
		quit(1, err, c, "%v: Unable to execute ", filename)
	}

	// And return a new front matter object
	return newC.fm
}

// HTML UTILITIES
// documentReady() inserts Javascript code to ensure that
// user-defined Javascript executes only after the full
// HTML DOM has been loaded.
func (c *config) documentReady() string {
	return c.fileToString(".poco/js/docready.js")
}

// defaultHomePage() Generates a simple home page as an HTML string
// Uses the file segment of dir as the the H1 title.
// Uses current directory if "." or "" are passed
func defaultHomePage(dir string) string {

	var indexMdFront = `---
stylesheets:
    - "https://cdn.jsdelivr.net/gh/pococms/poco/.poco/css/poquito.css"
---
`
	var indexMdBody = `
hello, world.

Learn more at [PocoCMS tutorials](https://pococms.com/docs/tutorials.html) 
`
	if dir == "" || dir == "." {
		dir, _ = os.Getwd()
	}
	h1 := "Welcome to " + filepath.Base(dir)
	page := indexMdFront +
		"# " + h1 + "\n" +
		indexMdBody
	return page
}

// tagSurround takes text and surrounds it with
// opening and closing tags, so
// tagSurround("header","WELCOME","\n") returns "<header>WELCOME</header>\n"
// You can optionally include text after, because sometimes it
// makes sense to include a newline after the closing tag.
func tagSurround(tag string, txt string, extra ...string) string {
	add := ""
	if len(extra) < 1 {
		add = ""
	} else {
		add = extra[0]
	}
	return "<" + tag + ">" + txt + "</" + tag + ">" + add
}

// StyleTags takes a list of tags and inserts them into right before the
// closing head tag, so they can override anything that came before.
// These are literal tags, not filenases.
// They're listed under "style-tags" in the front matter
// Returns them as a string. For clarity each tag is indented
// and ends with a newline.
// Example:
//
// style-tags:
//   - "h1{color:blue;}"
//   - "p{color:darkgray;}"
//
// Would yield:
//
//	"{color:blue;}\n\t\tp{color:darkgray;}\n"
//
// TODO: Document
// theme contains all the (lightweight) files needed for a theme:
// header.md, style sheets, etc.
type theme struct {

	// READ ONLY: Full pathname to theme directory
	dir string

	// Who created it, natch
	author string

	// Name for the theme with spaces and other characters allowed.
	// If the directory name is my-great-theme you might
	// want this to be "My Great! Theme"
	branding string

	// One or more sentences selling the theme.
	description string

	// Holds converted and template-parsed markdown source
	// for the <header> tag.
	header string
	// Filename for header specified in front matter.
	headerFilename string

	// Holds converted and template-parsed markdown source
	// for the <nav> tag.
	nav string
	// Filename for nav specified in front matter.
	navFilename string

	// Holds converted and template-parsed markdown source
	// for the <aside> tag.
	aside string
	// Filename for aside specified in front matter.
	asideFilename string

	// Holds converted and template-parsed markdown source
	// for the <footer> tag.
	footer string
	// Filename for footer specified in front matter.
	footerFilename string

	// Contents of LICENSE file. Can't be empty
	license string

	// Contents of README.md for this theme.
	readme string

	// Name of the theme is the same as directory with no pathname
	name string

	// True of a theme is named on this page
	present bool

	// Names of stylesheets
	stylesheetFilenames []string
	// The stylesheets for each theme are concantenated, then read
	// into this string. It's injected straight into the HTML for
	// each file using this theme.
	stylesheets string

	// Extra tags added right there on the Markdown page
	styleTagNames []string
	// The style tags for each theme are concantenated, then read
	// into this string. It's injected straight into the HTML for
	// each file using this theme.
	styleTags string
}

// there are no configuration files (yet) but this holds
// configuration info for the project, for example, template
// stylesheets and current file being processed.
// That stuff lives in the front matter of the home
// page (first checks for README.md, then checks for index.md)
type config struct {
	// Command-line -cleanup flag determines
	// whether or not the publish (aka WWW) directory gets deleted on start.
	cleanup bool

	// # of files copied to webroot

	copied int
	// mdCopied tracks # of Markdown files converted and copied to webroot
	mdCopied int

	// Name of Markdown file being processed. NOTE: read it with currentFile() method.
	currentFilename string

	// dumpfm command-line option shows the front matter of each page
	dumpFm bool

	// linkStyles true means stylesheets will not be inlined.
	linkStyles bool

	// List of all files being processed
	files []string

	// Front matter
	// front matter for current theme
	fm map[string]interface{}

	// front matter for global theme
	globalFm map[string]interface{}

	// front matter for current page
	pageFm map[string]interface{}

	// Full pathname of the root index file Markdown in the root directory.
	// If present, it's either "README.md" or "index.md"
	homePage string
	// The finished home page has to be preserved here because it's generated
	// before there's webroot directory.
	homePageStr string

	// Command-line flag -lang sets the language of the HTML files
	lang string

	// markdownExtensions are how PocoCMS figures out whether
	// a file is Markdown. If it ends in any one of these then
	// it gets converted to HTML.
	markdownExtensions searchInfo

	// Command-line flag -new generates a new project by this name
	newProjectStr string

	// Port localhost server runs on
	port string

	// Home directory for source code
	root string

	// Command-line flag -serve determing if running as
	// a localhost web server
	runServe bool

	// Command line flag -settings shows configuration values
	// instead of processing files
	settings bool

	// Command line flag -settings-after shows configuration values
	// after processing files
	settingsAfter bool

	// Command-line flag -skip lets you skip
	// the named files from being processed
	skip string

	// String slice list of items not to
	// process or send to webroot
	// that will contain everything
	// from both SkipPublish in the front matter and
	// everything in the -skip command-line flag.
	skipPublish searchInfo

	// Contents of a theme directory: the theme for the current page,
	// plus the global (default) theme directory.
	pageTheme   theme
	globalTheme theme

	// Command-line flag -themes shows installed themes
	themeList bool

	// Command-line flag -timestamp inserts a timestamp at the
	// top of the article when true
	timestamp bool

	// The --verbose flag. It shows progress as the site is created.
	// Required by the verbose() function.
	verboseFlag bool

	// Output directory for published files
	webroot string
}

// findHomePage() returns the source file used for the root
// index page in the root directory. Since README.md is
// commonly used, it takes priority. Next priority is index.md
// Set c.currentFilename to the home page when found
// Pre: c.root must be a fully qualified pathname
func (c *config) findHomePage() {
	// Look for "README.md" or "index.md" in that order.
	c.homePage = indexFile(c.root)
	if c.homePage != "" {
		c.currentFilename = c.homePage
		return
	}

	if !dirEmpty(c.root) {
		// No home page.
		// Directory has files.
		// User may not wish to create a new project.
		if promptYes(c.root + " isn't a PocoCMS project, but the directory has files in it.\nCreate a home page?") {
			c.newSite()
		} else {
			quit(1, nil, c, "Can't build a project without an index.md or README.md")
		}
	} else {
		// Empty dir, so create home page
		c.newSite()
		//writeDefaultHomePage(c, c.root)
	}
	c.currentFilename = c.homePage
}

// currentFile() returns the name of the file
// being processed, since it's displayed in
// two different places
func (c *config) currentFile() string {
	return c.currentFilename
}

// setRoot() obtains a fully qualified pathname for the home page source filename
// and its root directory.
// Pre: parseComandLine()
func (c *config) setRoot() {
	var err error
	// Determine home page, which may have been passed on command line.
	if c.root == "." || c.root == "" {
		// Handle most common case: no params, just process this directory.
		c.root = currDir()
	}
	// Something's left on the command line. It's presumed to
	// be a directory. Exit if that dir doesn't exist.
	if !filepath.IsAbs(c.root) {
		c.root, err = filepath.Abs(c.root)
		if err != nil {
			quit(1, err, nil, "Can't get absolute path for home page")
		}
	}
	// c.root finally established. Does it even exist?
	if !dirExists(c.root) {
		quit(1, nil, c, "Can't find the directory %v. You may want to create a new site by running %s -new %s.", c.root, filepath.Base(os.Args[0]), flag.Arg(0))
	}
}

// setWebroot() obtains a fully qualified pathname for the webroot, where all HTML output files
// and assets go.
// Creates webroot if it doesn't exist
// Pre: parseComandLine(), setRoot()
func (c *config) setWebroot() {
	// Webroot either defaulted to WWW or was given a new location from command line.
	// Don't know if it's valid.
	// Make sure there's a webroot directory
	// First job is to expand it completely.
	var err error

	if !filepath.IsAbs(c.webroot) {
		if c.root != currDir() {
			// Handle case where user has specified a different dir for the root
			// but not an absolute path for the webroot. In other words:
			//   poco ~/foo/bar
			// When not in the ~/foo/bar directory. The webroot is then
			// presumed to be a subdirectory of that root, not the current dir.
			c.webroot = filepath.Join(c.root, c.webroot)
		} else {
			c.webroot, err = filepath.Abs(c.webroot)
		}
		if err != nil {
			quit(1, err, nil, "Can't get absolute path for webroot")
		}
	}
}

// regularize() is given a root directory and a filename.
// filename may be a URL. It may be a full pathname.
// Or it may be relative to the directory.
// Deal with those cases.
// What gets returned is the name of a file that,
// if it exists, cant be downloaded.
// Returns "" on eror, not an error, which may be a mistake
func regularize(dir string, filename string) string {

	// Do nothing if the filename is a fully
	// qualified pathname.
	if filepath.IsAbs(filename) {
		return filename
	}

	// Do nothing if the filename is a URL
	if strings.HasPrefix(filename, "http") {
		return filename
	}

	s := filepath.Join(dir, filename)
	if f, err := filepath.Abs(s); err != nil {
		quit(1, err, nil, "Unable to produce absolute pathname for %s", s)
	} else {
		return f
	}
	return ""
}

// themeDescription() takes the name of a theme directory and returns
// the name of the page layout files, stylesheets, style tags
// in the theme object. It does not read any of those files in.
// So at the end you know what the assets are but haven't done
// anything with them.
// Returns an empty theme object if themeDir is empty.
func (c *config) themeDescription(themeDir string, possibleGlobalTheme bool) theme {
	debug("\t\t\tthemeDescription('%s')", themeDir)
	// Return value
	var theme theme
	// Leave if no theme specified and no global theme specified.
	//if themeDir == "" && !c.globalTheme.present {
	if themeDir == "" && !possibleGlobalTheme {
		theme.present = false
		debug("\t\t\t\treturning with no local or global theme for  %s", c.currentFile())
		return theme
	}

	// The theme is actually just a directory name.
	theme.dir = themeDir
	// The theme's heart is its README.md file, which lists
	// assets required by the theme.
	// Get its full path.
	themeReadme := filepath.Join(theme.dir, "README.md")
	if !fileExists(themeReadme) {
		theme.present = false
		quit(1, nil, c, "%s specified for %s can't be found", themeReadme,
			c.currentFile())
	} else {
		// Found it. Read its contents.
		theme.readme = c.fileToString(themeReadme)
	}

	// Make sure there's a LICENSE file
	license := filepath.Join(theme.dir, "LICENSE")
	if !fileExists(license) {
		theme.present = false
		quit(1, nil, c, "%s theme is missing a LICENSE file", c.pageTheme.dir)
	} else {
		// Found it. Read its contents.
		theme.license = c.fileToString(license)
		// Met minimal requirements for a theme.
		theme.present = true
	}

	// Get a new config object to avoid stepping on c.config
	tmpConfig := newConfig()
	// Get the front matter for this theme.
	tmpConfig.fm = tmpConfig.getFm(themeReadme)

	// The theme's README.md file has been located.
	// A temporary config object has been created.
	// Get from the theme's front matter, author, branding,
	// description, etc.
	theme.readFm(tmpConfig.fm)
	if possibleGlobalTheme && !c.globalTheme.present {
		// If this is a global-theme declaration, read that
		// into c
		c.globalTheme = theme
		debug("\t\t\t\tGlobal theme %s found", c.globalTheme.dir)
	} else {
		// It's a local theme declaration
		c.pageTheme = theme
		debug("\t\t\t\tLocal theme %s found", c.pageTheme.dir)
	}
	return theme
} // themeDescription()

// If there's a local footer, return it.
// If not and there's a global footer, return it.
func (c *config) footer() string {
	if c.pageTheme.present {
		return c.pageTheme.footer
	}
	if c.globalTheme.present {
		return c.globalTheme.footer
	}
	return ""
}

// If there's a local aside, return it.
// If not and there's a global aside, return it.
func (c *config) aside() string {
	if c.pageTheme.present {
		return c.pageTheme.aside
	}
	if c.globalTheme.present {
		return c.globalTheme.aside
	}
	return ""
}

// If there's a local nav, return it.
// If not and there's a global nav, return it.
func (c *config) nav() string {
	if c.pageTheme.present {
		return c.pageTheme.nav
	}
	if c.globalTheme.present {
		return c.globalTheme.nav
	}
	return ""
}

// If there's a local header, return it.
// If not and there's a global header, return it.
func (c *config) header() string {
	if c.pageTheme.present {
		return c.pageTheme.header
	}
	if c.globalTheme.present {
		return c.globalTheme.header
	}
	return ""
}

// layoutElement() takes a layout element file named in the front matter
// and generates HTML, but it executes templates also.
// A layout element is one of the HTML tags such
// as header, nav, aside, article, and a few others
// For more info on layout elements see:
// https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML/Document_and_website_structure#html_layout_elements_in_more_detail
//
// So, the priority order is:
//
//   - If the front matter says "SUPPRESS" in all caps then return empty string.
//     Example:
//     ---
//     aside: "SUPPRESS"
//     ---
//
//   - If there is a file named in the front matter, process and return its contents.
//     Example:
//     ---
//     aside: "mysidebar.md"
//     ---
//
// It can be a Markdown file, in which case no tags are needed,
// or an HTML file, in which the tags must be explicit.
// The easiest way is to use markdown.
// Fore example, suppose you have a header file named mdhead.md and
// it contains only the following:
//
// hello, world.
//
// The genereated HTML would be "<p><header id="header-poco">hello, world.</header></p>"
//
// Note that each layout element is given an id of {layoutelment}-poco as shown.
//
// If it's an HTML file,
// suppose you have a header file named head.html. It
// would be named in the front matter like this:
// ---
// Header: head.html
// ---
//
// The layout element file is expected to be a complete tag. For example,
// the header file could be as simple as this:
//
//	<header id="header-poco">hello, world.</header>
//
// This function would read in the head.html file (or whatever
// the file was named in the front matter) and insert it before the
// body of the document.
func (c *config) layoutElement(tag string, t *theme) {
	// Possible states.
	// GLOBAL: No layout file specified, but a global theme is present.
	//         This is the case where the home page has Theme: "foo" in the front matter.
	//         The current page doesn't have a theme specified in its front Matter.
	//         The global theme HTML is already in c.globalTheme.theme.header etc.
	// LOCAL:  No global theme was specified, but a layout file was specified.
	//         This is the case where there's a Theme: "foobar" on the current page
	//         but the home page didn't contain a theme designation.
	//         In case 2 the layout file needs to be read in from the (non-global)
	//         theme specified on this page.
	const (
		GLOBAL = 1
		HTML   = 2
	)

	state := GLOBAL
	filename := ""

	// Converted/templated HTML */
	s := ""
	switch tag {
	case "header":
		if t.headerFilename != "" && t.headerFilename != "SUPPRESS" {
			t.headerFilename = regularize(t.dir, t.headerFilename)
			filename = t.headerFilename
		}
	case "nav":
		if t.navFilename != "" && t.navFilename != "SUPPRESS" {
			t.navFilename = regularize(t.dir, t.navFilename)
			filename = t.navFilename
		}
	case "aside":
		if t.asideFilename != "" && t.asideFilename != "SUPPRESS" {
			t.asideFilename = regularize(t.dir, t.asideFilename)
			filename = t.asideFilename
		}
	case "footer":
		if t.footerFilename != "" && t.footerFilename != "SUPPRESS" {
			t.footerFilename = regularize(t.dir, t.footerFilename)
			filename = t.footerFilename
		}
	}

	// Handle case 1.
	if filename == "" {
		return
	}

	// If HTML file return unchanged
	if path.Ext(filename) == ".html" {
		if fileExists(filename) {
			s = c.fileToString(filename)
			state = HTML
		} else {
			quit(1, nil, c, "HTML theme layout file %s not found", filename)
		}
	}

	if state != HTML {
		s = convertMdYAMLFileToHTMLFragmentStr(filename, c)
		if s != "" {
			s = "<" + tag + " id=\"" + tag + "-poco" + "\"" + ">" + s + "</" + tag + ">"
		}
	}

	switch tag {
	case "header":
		t.header = s
	case "nav":
		t.nav = s
	case "aside":
		t.aside = s
	case "footer":
		t.footer = s
	}

}

// setupGlobals() sets sitewide values such as
// home page, language, underlying theme, etc.
// Pre: parseCommandline()
func (c *config) setupGlobals() { //

	// Determine home page directory and filename
	c.setRoot()

	// If a file ends in any one of these extensions then
	// it gets converted to HTML.
	c.markdownExtensions.list = []string{".md", ".mkd", ".mdwn", ".mdown", ".mdtxt", ".mdtext", ".markdown"}

	// Set defaults for files and dirs to skip
	c.skip = "node_modules .git .DS_Store .gitignore " + pocoDir

	// Determine output directory for all HTML and assets (webroot)
	c.setWebroot()

	var err error
	// Root dir exists. Now change to it.
	if err = os.Chdir(c.root); err != nil {
		quit(1, err, c, "Unable to change to new root directory %s", c.root)
	}

	// Get name of home page source file, either README.md (first
	// priority if index.md is present) or index.md
	c.findHomePage()

	// Display home page filename in verbose mode. Same as
	// elsewhere in buildSite for all the other files.
	// Only read it through c.currentFile() after this
	c.currentFilename = c.homePage

	// Create a list of files and dirs to skip when processing
	c.getSkipPublish()

	// Display name of file being processed
	c.verbose(c.currentFile())

	// Prevent the home page from being read and converted again.
	c.skipPublish.AddStr(c.currentFile())

	// Convert home page to HTML
	c.homePageStr, _ = buildFileToTemplatedString(c, c.currentFile())

} // setupGlobals

func (c *config) styleTags() string {
	t := c.getStyleTags(c.pageFm)
	if t != "" {
		t = tagSurround("style", t, "\n")
	}
	return t
}

func (c *config) getStyleTags(fm map[string]interface{}) string {
	styleTagNames := fmStrSlice("style-tags", fm)
	styleTags := ""
	for _, tag := range styleTagNames {
		s := fmt.Sprintf("\t\t%s\n", tag)
		styleTags = styleTags + s
	}
	return styleTags
}

// linkStylesheets() extracts stylesheet references and
// creates link tags for them.
// By this time loadTheme() has been called. c.fm has hardcoded
// style tags, c.pageTheme and c.globalTheme have lists of
// stylesheets. Now it's time to turn those into HTML
// elements.
//
// See also inlineStylesheets(), which inserts stylesheet
// code directly into the HTML document
func (c *config) linkStylesheets() string {

	debug("\tlinkStylesheets()")
	//pageStyles := sliceToStylesheetStr("", c.pageTheme.stylesheetFilenames)
	pageStyles := c.getStyleTags(c.fm)
	debug("Front matter: %+v", c.fm)
	debug("\t\tpage styles: %v", pageStyles)

	// If there's a page theme,
	if c.pageTheme.present {
		//wait("c.pageTheme.dir: %s", c.pageTheme.dir)
		themeStyles := sliceToStylesheetStr(c.pageTheme.dir, c.pageTheme.stylesheetFilenames)
		//pageStyles = sliceToStylesheetStr(c.pageTheme.dir,pageStyles)
		debug("\t\ttheme styles: %v", themeStyles)
		return themeStyles + pageStyles
	}
	if c.globalTheme.present {
		globalThemeStyles := sliceToStylesheetStr(c.globalTheme.dir, c.globalTheme.stylesheetFilenames)
		debug("\t\tglobal styles: %v", globalThemeStyles)
		return globalThemeStyles + pageStyles
	}
	return pageStyles
}

// sliceToStylesheetStr takes a slice of simple stylesheet names, such as
// [ "foo.css", "bar.css" ] and converts it into a string
// consisting of stylesheet link tags separated by newlines:
//
// <link rel="stylesheet" href="foo.css"/>
// <link rel="stylesheet" href="bar.css"/>
func sliceToStylesheetStr(dir string, sheets []string) string {
	if len(sheets) <= 0 {
		return ""
	}
	var tags string
	for _, sheet := range sheets {
		// TODO: Should probably use portable path delimiters
		tag := fmt.Sprintf("\t<link rel=\"stylesheet\" href=\"%s/%s\">\n", dir, sheet)
		tags += tag
	}
	return tags
}

// inlineStylesheets() directly injects stylesheet code into the
// HTML document instead of linking to it.
// See also linkStylesheets(), which links to stylesheet
// instead of inserting directly into the HTML document
func (c *config) inlineStylesheets(dir string) string {
	debug("\t\t\t\tinlineStyleSheets(%s)", dir)
	//debug("\t\tfront matter: %+v", c.fm)
	// Return value
	s := ""

	// Look for stylesheets named on this page,
	// which ave the highest priority.
	slice := fmStrSlice("stylesheetFilenames", c.fm)
	stylesheets := ""
	if len(slice) > 0 {
		debug("\t\t\t\t\tslice of on-page stylesheets: %+v", slice)
		// Collect all the stylesheets mentioned.
		// Concatenate them into a big-ass string.
		for _, filename := range slice {
			// Get full pathname or URL of file.
			fullPath := regularize(dir, filename)
			debug("\t\t\t\t\t\tc.fm: %s", fullPath)
			// If the file is local, read it in.
			// If it's at a URL, download it.
			// For debugging purposes, add commment with filename
			s = "\n\n/* " + filename + " */\n" + c.getWebOrLocalFileStr(fullPath)
			stylesheets = s + "\n"
		}
	}

	// Get list of stylesheets for the page theme, if there is one.
	// It overrides any global theme so exit afterwards.
	if c.pageTheme.present {
    debug("\t\t\t\t\tPage theme for page %s: %s", c.currentFile(), c.pageTheme.name)
		slice = c.pageTheme.stylesheetFilenames
		// Collect all the stylesheets mentioned.
		// Concatenate them into a big-ass string.
		for _, filename := range slice {
			// Get full pathname or URL of file.
			fullPath := regularize(filepath.Join(dir, c.pageTheme.dir), filename)

			// For debugging purposes, add commment with filename
			s = "\n/* " + filepath.Base(filename) + "*/\n" +
				// If the file is local, read it in.
				// If it's at a URL, download it.
				c.getWebOrLocalFileStr(fullPath)
			stylesheets = stylesheets + s + "\n"
		}
		// Page theme overrides global so exit with that.
		if s != "" {
			return "<style>" + stylesheets + "</style>" + "\n"
		}
	}

	// xxx
	if c.globalTheme.present {
    debug("\t\tGlobal theme for page %s: %s", c.currentFile(), c.globalTheme.name)
		slice = c.globalTheme.stylesheetFilenames
		// Collect all the stylesheets mentioned.
		// Concatenate them into a big-ass string.
		for _, filename := range slice {
			// Get full pathname or URL of file.
			fullPath := regularize(filepath.Join(dir, c.globalTheme.dir), filename)

			// For debugging purposes, add commment with filename
			s = "\n/* " + filepath.Base(filename) + "*/\n" +
				// If the file is local, read it in.
				// If it's at a URL, download it.
				c.getWebOrLocalFileStr(fullPath)
			stylesheets = stylesheets + s + "\n"
		}

    if s != "" {
      return "<style>" + stylesheets + "</style>" + "\n"
    }
  }
	return ""
}

// stylesheets() generates stylesheet tags required by themes
// priority.
func (c *config) stylesheets() string {
	// Return value
	s := ""
	debug("\t\t\tstylesheets()")
	if c.linkStyles {
		// Normally stylesheets are inlined.
		// This allows them to be linked to as usual.
		s = c.linkStylesheets()
	} else {
		// Inline stylesheets as usual
		s = c.inlineStylesheets(c.root)
	}
	return s
}



// themeDataStructures obtains the data structures
// for the theme located at dir.
// If possibleGlobalTheme is true, it means the hoped-for theme is
// the global theme.
// If valid theme (all that's required: README.md and LICENSE files),
// then set theme.present to true. Covers special case for the global
// theme.
func (c *config) themeDataStructures(dir string, possibleGlobalTheme bool) *theme {
	// The theme is actually just a directory name.
	var theme theme
	theme.dir = dir
	// The theme's heart is its README.md file, which lists
	// assets required by the theme.
	// Get its full path.
	themeReadme := filepath.Join(theme.dir, "README.md")
	if !fileExists(themeReadme) {
		theme.present = false
		quit(1, nil, c, "%s specified for %s can't be found", themeReadme,
			c.currentFile())
	} else {
		// Found it. Read its contents.
		theme.readme = c.fileToString(themeReadme)
	}

	// Make sure there's a LICENSE file
	license := filepath.Join(theme.dir, "LICENSE")
	if !fileExists(license) {
		theme.present = false
		quit(1, nil, c, "%s theme is missing a LICENSE file", c.pageTheme.dir)
	} else {
		// Found it. Read its contents.
		theme.license = c.fileToString(license)
		// Met minimal requirements for a theme.
		theme.present = true
		// On the home page. The request is for the global theme.
		if possibleGlobalTheme {
			c.globalTheme.present = true
		}
	}

	// Strip path off theme to get its name
	theme.name = filepath.Base(dir)

	// Get a new config object to avoid stepping on c.config
	tmpConfig := newConfig()
	// Get the front matter for this theme.
	tmpConfig.fm = tmpConfig.getFm(themeReadme)

	// The theme's README.md file has been located.
	// A temporary config object has been created.
	// Get from the theme's front matter, author, branding,
	// description, etc.
	theme.readFm(tmpConfig.fm)

	return &theme
}

// getThemeData() obtains the name and data structures for
// any theme named on this page: the pageTheme or the globalTheme
// (Possibly both if on the home page)
// filename is the name of the current Markdown source file.
func (c *config) getThemeData(filename string) {
	debug("\tgetThemeData(%s)", filename)
	pageThemeDir := fmStr("theme", c.pageFm)

	// If NOT on the home page, check this page for a theme.
	if filename != c.homePage {
		// Nothing else to do if no theme named.
		// c.pageFm already has front matter for this page.
		if pageThemeDir == "" {
			// No theme. Business is concluded.
			// c.theme.present is already false
			return
		}

    // Check for a global theme
    
	  //debug("\t\tpageThemeDir: %s", pageThemeDir)
		// A page theme has been named. Not known to be valid.
		// Again, this is not the home page.
		// getThemeData(), new version
		if dirExists(pageThemeDir) {
			c.pageTheme = *c.themeDataStructures(pageThemeDir, false)
			//debug("\tgetThemeData(): not on home page. c.pageTheme:\n%+v\n", c.pageTheme)
			return
		}

		// On the home page. Check for a global-theme but also a page theme.
	} else {
    debug("\t\ton home page")
	  //pageThemeDir := fmStr("global-theme", c.pageFm)
		// First check for a page theme. If both are named, this will
		// override the global theme.
		if pageThemeDir != "" && dirExists(pageThemeDir) {
			// A page theme has been named. Not known to be valid.
			// Again, this is not the home page.
			// getThemeData(), new version
			c.pageTheme = *c.themeDataStructures(pageThemeDir, false)
			//debug("\tgetThemeData(): on home page. c.pageTheme:\n%+v\n", c.pageTheme)
		}
		// This is the home page. Check for a global theme.
		globalThemeDir := fmStr("global-theme", c.pageFm)
		if dirExists(globalThemeDir) {
			c.globalTheme = *c.themeDataStructures(globalThemeDir, true)
			//wait("Found global theme %s: %+v", globalThemeDir, c.globalTheme)
		}
	}
}

// loadTheme() is passed the current source filename.
// If a page theme is named in the front matter, its description
// is read. If at the home page, it reads the global theme, if any.
// It is possible at the home page to have both page theme
// and global themes names. In that case the page theme takes priority
// on the home page, as it would any other page.
// TODO: docs
func (c *config) loadTheme(filename string) {
	// Obtain the front matter for this page.
	// Any values such as header, footer, etc. will override
	// their corresponding local or global themes.
  // c.pageFm = map[string]interface{}{}
	c.pageFm = c.getFm(filename)

  debug("loadTheme(%s)", filename)

	// Get the page theme, if any.
	// If on the home page, look for both global
	// and local theme names.
	// Load data structures for those themes.
	c.getThemeData(filename)

	// Handle case where no theme was named, and there's no
	// global theme named either.
	// Any theme info is just random stuff in
	// the front matter of the current Markdown file.
	if !c.pageTheme.present && !c.globalTheme.present {
		return
	}

	// If a page theme has been named, the data structures are ready.
	// Read in its style sheets, style tags, and page layout elements.
	if c.pageTheme.present {
		// Local theme takes priority
		c.addPageElements(&c.pageTheme)
		debug("\tpage theme %s found", c.pageTheme.dir)
		return
	}

	// If a global theme has been named, the data structures are ready.
	// Read in its style sheets, style tags, and page layout elements.
	if c.globalTheme.present {
		// Local theme takes priority
		c.addPageElements(&c.globalTheme)
		debug("\tglobal theme %s found", c.globalTheme.dir)
		return
	}

} // loadTheme (new version)



func (c *config) addPageElements(t *theme) {
	c.layoutElement("header", t)
	c.layoutElement("nav", t)
	c.layoutElement("aside", t)
	c.layoutElement("footer", t)
}

func (t *theme) readFm(fm map[string]interface{}) {
	t.author = fmStr("author", fm)
	t.branding = fmStr("branding", fm)
	t.description = fmStr("description", fm)
	t.headerFilename = fmStr("header", fm)
	t.navFilename = fmStr("nav", fm)
	t.asideFilename = fmStr("aside", fm)
	t.footerFilename = fmStr("footer", fm)
	t.styleTagNames = fmStrSlice("style-tags", fm)
	//debug("\t\treadFm() styleTagNames: %+v", t.styleTagNames)
	t.stylesheetFilenames = fmStrSlice("stylesheets", fm)
}

// newConfig allocates a config object.
// sitewide configuration info.
func newConfig() *config {
	config := config{}
	return &config

}

// parseCommandLine obtains command line flags and
// initializes values.
func (c *config) parseCommandLine() {
	// cleanup determines whether or not the publish (aka WWW) directory
	// gets deleted on start.
	flag.BoolVar(&c.cleanup, "cleanup", true, "Delete publish directory before converting files")

	// debugFrontmatter command-line option shows the front matter of each page
	flag.BoolVar(&c.dumpFm, "dumpfm", false, "Shows the front matter of each page")

	// linkStyles controls whether stylesheets are inlined (normally they are)
	flag.BoolVar(&c.linkStyles, "link-styles", false, "Link to stylesheets instead of inlining them")

	// lang sets HTML lang= value, such as <html lang="fr">
	// for all files
	flag.StringVar(&c.lang, "lang", "en", "HTML language designation, such as en or fr")

	// new creates a directory, sample index.md, and pocoDir
	flag.StringVar(&c.newProjectStr, "new", "", "Create a new site")

	// Port server runs on
	flag.StringVar(&c.port, "port", ":54321", "Port to use for localhost web server")

	// Directory project lives in
	flag.StringVar(&c.root, "root", ".", "Starting directory of the project")

	// -settings command-line shows configuration values
	// instead of processing files
	flag.BoolVar(&c.settings, "settings", false, "Shows configuration values instead of processing site")

	// Run as server without processing any files
	flag.BoolVar(&c.runServe, "serve", false, "Run as a web server on localhost")

	// skip lets you skip the named files from being processed
	flag.StringVar(&c.skip, "skip", "node_modules/ .git/ .DS_Store/ .gitignore", "List of files to skip when generating a site")

	// Command line flag -settings-after shows configuration values
	// after processing files
	flag.BoolVar(&c.settingsAfter, "settings-after", false, "Shows configuration values after processing site")

	// Command-line flag -themes lists themes in the poco directory
	flag.BoolVar(&c.themeList, "themes", false, "Show themes in "+pocoDir+" directory")

	// Command-line flag -timestamp inserts a timestamp at the
	// top of the article when true
	flag.BoolVar(&c.timestamp, "timestamp", false, "Insert timestamp at top of home page article")

	// Verbose shows progress as site is generated.
	flag.BoolVar(&c.verboseFlag, "verbose", false, "Display information about project as it's generated")

	// webroot flag is the directory used to house the final generated website.
	flag.StringVar(&c.webroot, "webroot", "WWW", "Subdirectory used for generated HTML files")

	// Process command line flags such as --verbose, --title and so on.
	flag.Parse()

	// Collect configuration info for this project

	// See if a directory was specified.
	c.root = flag.Arg(0)

	// TODO: Not sure this is the right place to run this, but see
	// issue #19
	if c.themeList {
		print(c.themeDirContents())
		os.Exit(0)
	}

}

func main() {
	c := newConfig()
	// No file was given on the command line.
	// Build the project in place.

	// Collect command-line flags, directory to build, etc.
	c.parseCommandLine()

	// New project requested?
	if c.newProjectStr != "" {
		c.newProject(c.newProjectStr)
	}

	// Prevent generating a project if in same directory as poco
	dir, err := os.Executable()
	if err != nil {
		quit(1, err, c, "Can't determine executable directory", "")
	}

	// If in the poco main directory and no project was specified,
	// then prevent turning the poco main directory into a project.
	// OTOH if something is specified, e.g. poco ~/mysite, it's
	// okay to generate that site.
	if currDir() == filepath.Dir(dir) && c.root == "" {
		quit(1, err, c, "%s", "Don't run poco in its own directory. Quitting.")
	}

	// Obtain README.md or index.md.
	// Read in the front matter to get its config information.
	// Set values accordingly.
	c.setupGlobals()

	// Probably not for public release.
	// Lets me search for a new session in command line
	// history or output file
	if c.verboseFlag {
		print("project: %s =============", filepath.Base(c.root))
	}

	// If -serve flag was used just run as server.
	if c.runServe {
		if dirExists(c.webroot) {
			c.serve()
		} else {
			// Or more likely it quits silently
			quit(1, nil, c, "Can't find webroot directory %s", c.webroot)
		}
	}

	// If -settings flag just show config values and quit
	if c.settings {
		c.dumpSettings()
		os.Exit(0)
	}

	// Generate the site based in c.root. Output its contens to c.webroot.
	c.buildSite()

	// If -settings-after flag just show config values and quit
	if c.settingsAfter {
		c.dumpSettings()
	}

	final := filepath.Join(c.webroot, "index.html")
	if !c.verboseFlag {
		print("Site published to %s", final)
	} else {
		print("%s Site published to %s", theTime(), final)
	}
}

// TEMPLATE FUNCTIONS

// doTemplate takes HTML in source, expects parsed front
// matter in fm, and executes Go templates
// against the source.
// Returns a string containing the HTML with the
// template values embedded.
func doTemplate(templateName string, source string, c *config) (string, error) {
	if templateName == "" {
		templateName = "PocoCMS"
	}
	tmpl, err := template.New(templateName).Parse(source)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, c.fm)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

// buildFileToFile converts a file from Markdown to HTML, generates an output file,
// and returns name of destination file
// Used for every Markdown page on the site.
func buildFileToFile(c *config, filename string, debugFrontMatter bool) (outfile string) {
	// Convert Markdown file filename to raw HTML, then assemble a complete HTML document to be published.
	// Return the document as a string.
	html, htmlFilename := buildFileToTemplatedString(c, filename)
	// Write the contents of the completed HTML document to a file.
	// Return the name of the converted file
	return stringToFile(c, htmlFilename, html)
}

// buildFileToTemplatedString converts a file from Markdown to raw HTML,
// pulls in everything required to create a complete HTML document,
// executes templates,
// generates an output file,
// and returns name of the destination HTML file
// Does not check if the input file is Markdown.
// Returns the string and the filename
func buildFileToTemplatedString(c *config, filename string) (string, string) {
	// Exit silently if not a valid file
	if filename == "" || !fileExists(filename) {
		return "", ""
	}
	c.loadTheme(filename)
	// This will be the proposed name for the completed HTML file.
	dest := ""
	// Convert the Markdown file to an HTML string
	if rawHTML, err := mdYAMLFileToHTMLString(c, filename); err != nil {
		quit(1, err, c, "Error converting Markdown file %v to HTML", filename)
		return "", ""
	} else {
		// Strip original file's Markdown extension and make
		// the destination files' extension HTML
		dest = replaceExtension(filename, "html")
		// Take the raw converted HTML and use it to generate a complete HTML document in a string
		finishedDocument := c.assemble(c.currentFile(), rawHTML)
		// Return the finishled document and its filename
		return finishedDocument, dest
	}
}

// converts all files (except those in skipPublish.List) to HTML,
// and deposits them in webroot. Attempts to create webroot if it
// doesn't exist. webroot is expected to be a subdirectory of
// projectDir.
func (c *config) buildSite() {

	var err error
	// Make sure it's a valid site. If not, create a minimal home page.

	// Change to requested directory
	if err = os.Chdir(c.root); err != nil {
		quit(1, err, c, "Unable to change to directory %s", c.root)
	}

	// Delete webroot directory unless otherwise requested
	if c.cleanup {
		if err := os.RemoveAll(c.webroot); err != nil {
			quit(1, err, c, "Unable to delete webrootdirectory %v", c.webroot)
		}
	}

	// Collect all the files required for this project.
	c.files, err = getProjectTree(".", c.skipPublish)
	// c.files is a list of files with pathnames relative to c.root
	if err != nil {
		quit(1, err, c, "Unable to get directory tree")
	}

	// Create the webroot directory
	if !dirExists(c.webroot) {
		err := os.MkdirAll(c.webroot, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			quit(1, err, c, "Unable to create webroot directory %s", c.webroot)
		}
	}

	// no processing.
	// # of non-Markdown copied
	assetsCopied := 0
	// First write out home page
	target := filepath.Join(c.webroot, "index.html")
	target = stringToFile(c, target, c.homePageStr)
	// # of Markdown files processed
	// Start at 1 because home page
	c.mdCopied = 1

	// Main loop. Traverse the list of files to be copied.
	// If a file is Markdown as determined by its file extension,
	// convert to HTML and copy to output directory.
	// If a file isn't Markdown, copy to output directory with
	for _, filename := range c.files {

		// Full pathmame of file to be copied (may be converted to HTML first)
		source := filepath.Join(c.root, filename)

		c.currentFilename = source
		c.verbose(c.currentFile())

		// Full pathname of location of copied file in webroot
    // If it's an asset (non-Markdown file), it will be 
    // copied as is.
		target := filepath.Join(c.webroot, filename)

		// Full pathname of output directory for copied files
    targetDir := c.webroot
    // Create the webroot directory if it doesn't exist
		if !dirExists(targetDir) {
			err := os.MkdirAll(targetDir, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				quit(1, err, c, "Unable to create directory %s in webroot", targetDir)
			}
		}
		// Obtain file extension.
    ext := path.Ext(c.currentFile())

		// Replace converted filename extension, from markdown to HTML.
		// Only convert to HTML if it has a Markdown extension.
		if c.markdownExtensions.Found(ext) {
			// It's a markdown file. Convert to HTML,
			// then rename with HTML extensions.
			//jjHTML, _ := buildFileToTemplatedString(c, filename)
			HTML, _ := buildFileToTemplatedString(c, c.currentFile())
			target := filepath.Join(c.webroot, filename)
			//target := filepath.Join(c.webroot, c.currentFile())
			target = replaceExtension(target, "html")
			target = stringToFile(c, target, HTML)
			c.mdCopied++

		} else {
			// It's an asset. Just pass through.
			//copyFile(c, source, target)
			copyFile(c, source, target)
			assetsCopied++
		}

		c.copied += 1
	}
	// ALL files now copied
	// This is where the files were published
	ensureIndexHTML(c.webroot, c)
	// Display all files, Markdown or not, that were processed
	c.verbose("%s converted, %s copied. %d total", fileCount("Markdown", c.mdCopied), fileCount("asset", assetsCopied), c.copied)
	//c.copied, mdCopied, assetsCopied)
} // buildSite()

// fileCount returns a string containing
// a number followed by " word" or " words".
func fileCount(filetype string, count int) string {
	s := fmt.Sprintf("%d %s files", count, filetype)
	if count == 1 {
		return fmt.Sprintf("%d %s file", count, filetype)
	}
	return s
}

// ensureIndexHTML makes sure there's an index.html file
// in the webroot directory. It's required because some existing
// projects use README.md instead of index.md.
// Web servers don't recognize README.html as
// the home page, so an existing README.html gets renamed index.html.
// If both exist, README.md takes priority over index.md.
func ensureIndexHTML(path string, c *config) {
	readmeHTML := filepath.Join(path, "README.html")
	indexHTML := filepath.Join(path, "index.html")

	// if neither index.html nor README.html, then they
	// were missing source files to begin with.
	if !fileExists(readmeHTML) && !fileExists(indexHTML) {
		quit(1, nil, c, "No README.html or index.html could be found in %v", path)
	}

	// Both README.html and index.html exist.  Or
	// README.html exists but no index.html exists.
	// Rename README.html
	if fileExists(readmeHTML) && (fileExists(indexHTML) || !fileExists(indexHTML)) {
		err := os.Rename(readmeHTML, indexHTML)
		if err != nil {
			quit(1, err, c, "Unable to rename %v ", readmeHTML)
		}
	}
}

// getSkipPublish() obtains the list of directories
// and files to ignore during creation of  a site.
// An example would be something like:
//
//	"node_modules/ .git/ .DS_Store/ .gitignore"
//
// This should be established only on the home page
// and from the -skip command-line option
func (c *config) getSkipPublish() {

	// The home page must be processed first, and
	// once only. So it should be in the list already.

	// Add anything from the -skip command line option
	list := strings.Split(c.skip, " ")
	c.skipPublish.list = append(c.skipPublish.list, list...)
	c.skipPublish.AddStr(".backup")

	// Get what's specified in the home page front matter
	localSlice := fmStrSlice("ignore", c.fm)
	c.skipPublish.list = append(c.skipPublish.list, localSlice...)
}

// isProject() looks at the structure of the specified directory
// and tries to determine if there's already a project here.
func isProject(path string) bool {
	// If the directory doesn't exist, that's easy.
	if !dirExists(path) {
		return false
	}
	if indexFile(path) == "" {
		return false
	} else {
		return true
	}
}

// indexFile looks in the specified path for either index.md
// or README.md. Returns that filename if it exists.
// If it has both, README.md takes priority.
func indexFile(path string) string {
	readmeMd := filepath.Join(path, "README.md")
	if fileExists(readmeMd) {
		return readmeMd
	}
	indexMd := filepath.Join(path, "index.md")
	if fileExists(indexMd) {
		return indexMd
	}
	return ""
}

// SYSTEM UTILITIES
// curDir() returns the current directory name.
func currDir() string {
	if path, err := os.Getwd(); err != nil {
		return "unknown directory"
	} else {
		return path
	}
}

// dirEmpty() returns true if the specified directory is empty.
// Gratefully stolen from:
// https://stackoverflow.com/questions/30697324/how-to-check-if-directory-on-path-is-empty
func dirEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
}

// Get full path where executable lives,
// minus the name of the executable itself.
func executableDir() string {
	ex, err := os.Executable()
	if err != nil {
		quit(1, err, nil, "Can't figure out PocoCMS pathname")
	}
	// Amputate the actual filename
	return filepath.Dir(ex)
}

// FILE UTILITIES
// copyFile, well, does just that. Doesn't return errors.
func copyFile(c *config, source string, target string) {
	//c.verbose("\tcopyFile(%s,%s)", source, target)
	if source == target {
		quit(1, nil, c, "copyFile: %s and %s are the same", source, target)
	}
	if source == "" {
		quit(1, nil, c, "copyFile: no source file specified")
	}
	if target == "" {
		quit(1, nil, c, "copyFile: no destination file specified for file %s", source)
	}
	var src, trgt *os.File
	var err error
	if src, err = os.Open(source); err != nil {
		quit(1, err, c, "copyFile: Unable to open file %s", source)
	}
	defer src.Close()

	if trgt, err = os.Create(target); err != nil {
		quit(1, err, c, "copyFile: Unable to create file %s", target)
	}
	if _, err := trgt.ReadFrom(src); err != nil {
		quit(1, err, c, "Error copying file %s to %s", source, target)
	}

}

// copyPocoDir copies the embedded .poco directory into
// the current directory, which is expected to be
// a new project directory.
func (c *config) copyPocoDir(f embed.FS, dir string) error {
	if dir == "" {
		dir = currDir()
	}
	return fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			quit(1, err, c, "Problem walking .poco dir")
		}
		if d.IsDir() {
			// It's a directory. Create it in the target location.
			// Easy to do because we're guaranteed inside the project dir.
			err := os.MkdirAll(path, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				quit(1, err, c, "Unable to copy embedded directory %s", c.webroot)
			}
		} else {
			// Build a full path for the source file to copy.
			// The source file is in the same directory as
			// the poco executable.
			source := filepath.Join(executableDir(), path)
			// Destination is just path, which is guaranteed to
			// be a subdirectory of the current (new project) directory.
			copyFile(c, source, path)
		}
		return nil
	})
}

// dirExists() returns true if the name passed to it is a directory.
func dirExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

// fileExists() returns true, well, if the named file exists
func fileExists(filename string) bool {
	if filename == "" {
		return false
	}
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// fileToBuf() reads the named file into a byte slice and returns
// that byte slice. In the spirit of HTML it simply returns an empty
// slice on failure.
func fileToBuf(filename string) []byte {
	if !fileExists(filename) {
		return []byte{}
	}
	var input []byte
	var err error
	// Read the whole file into memory as a byte slice.
	input, err = ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}
	}
	return input
}

// fileToString() sucks up a file and returns its contents as a string.
func (c *config) fileToString(filename string) string {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		quit(1, err, c, "")
	}
	return string(input)
}

// replaceExtension() is passed a filename and returns a filename
// with the specified extension.
func replaceExtension(filename string, newExtension string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + newExtension
}

// stringToFile creates a file called filename without checking to see if it
// exists, then writes contents to it.
// filename is a fully qualified pathname.
// contents is the string to write
// Returns filename
func stringToFile(c *config, filename, contents string) string {
	var out *os.File
	var err error
	if out, err = os.Create(filename); err != nil {
		quit(1, err, c, "stringToFile: Unable to create file %s", filename)
	}
	if _, err = out.WriteString(contents); err != nil {
		quit(1, err, c, "Error writing to file %s", filename)
	}
	return filename
}

// downloadFile() tries to read in the named URL as text and return
// its contents as a string.
func (c *config) downloadTextFile(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		quit(1, err, c, "Unable setting up to GET file %s", url)
	}
	req.Header.Set("Accept", "application/text")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		quit(1, err, c, "Unable to download file %s", url)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		quit(1, err, c, "Unable to read reponse body for %s", url)
	}

	return string(b)

}

// getWebOrLocalFileStr reads filename and returns it as a string.
// If string starts with http or https fetches it from the web.
func (c *config) getWebOrLocalFileStr(filename string) string {
	// Return value: contents of file are stored here
	s := ""

	// Handle case of URLs as opposed to local file
	if strings.HasPrefix(filename, "http") {
		// TODO: Check for redirect?
		// https://golangdocs.com/golang-download-files
		s = c.downloadTextFile(filename)
		return s
		//
	}

	// Handle case of local file with relative path
	if !filepath.IsAbs(filename) {
		// TODO: Try replacing with filepath.Abs
		fullPath := filepath.Join(c.pageTheme.dir, filename)
		s = c.fileToString(fullPath)
		return s
	}

	// Handle case of local file with absolute path
	s = c.fileToString(filename)
	return s
}

// newSite() creates a default home page
// based in the root directory name,
// then copies the .poco directory in
func (c *config) newSite() {
	writeDefaultHomePage(c, c.root)
	c.copyPocoDir(pocoFiles, "")
}

// Generates a simple home page
// and writes it to index.md in dir. Uses the file
// segment of dir as the the H1 title.
func writeDefaultHomePage(c *config, dir string) {
	html := defaultHomePage(dir)
	pathname := filepath.Join(dir, "index.md")
	c.homePage = stringToFile(c, pathname, html)
}

// SLICE UTILITIES
// Searching a sorted slice is fast.
// This tracks whether the slice has been sorted
// and sorts it on first search.
// TODO: document
type searchInfo struct {
	list   []string
	sorted bool
}

func (s *searchInfo) Sort() {
	sort.Slice(s.list, func(i, j int) bool {
		s.sorted = true
		return s.list[i] <= s.list[j]
	})
}

func (s *searchInfo) AddStr(add string) {
	if s.Found(add) {
		return
	}
	s.list = append(s.list, add)
	s.Sort()
}

func (s *searchInfo) Found(searchFor string) bool {
	if !s.sorted {
		s.Sort()
	}
	var pos int
	l := len(s.list)
	pos = sort.Search(l, func(i int) bool {
		return s.list[i] >= searchFor
	})
	return pos < l && s.list[pos] == searchFor
}

// DIRECTORY TREE

func visit(files *[]string, skipPublish searchInfo) filepath.WalkFunc {

	// Find out what directories to exclude
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Quietly fail if unable to access path.
			return err
		}
		isDir := info.IsDir()
		// Obtain just the full pathname
		name := info.Name()

		// Inside the Walk func called by getProjectTree

		// Skip any directory to be excluded, such as
		// the pub and .git directores
		if skipPublish.Found(name) && isDir {
			return filepath.SkipDir
		}

		// It may be just a filename on the exclude list.
		name = filepath.Join(currDir(), name)
		if skipPublish.Found(name) {
			return nil
		}

		// Don't add directories to this list.
		if !isDir {
			*files = append(*files, path)
		}
		return nil
	}
}

// Obtain a list of all files in the specified project tree starting
// at the root.
// Ignore items in exclude.List
func getProjectTree(path string, skipPublish searchInfo) (tree []string, err error) {
	var files []string

	err = filepath.Walk(path, visit(&files, skipPublish))
	if err != nil {
		return []string{}, err
	}
	return files, nil
}

// Generate HTML

// linkTags() obtains the list of link tags from the "LinkTags" front matter
// and inserts them into the document.
func (c *config) linktags() string {
	linkTags := fmStrSlice("linktags", c.fm)
	if len(linkTags) < 1 {
		return ""
	}
	tags := ""
	for _, tag := range linkTags {
		tags += "\t" + tag + "\n"
	}
	return tags
}

// metatag() generates a metatag such as <meta name="description"content="PocoCMS: Markdown-based CMS in 1 file, written in Go">
func metatag(tag string, content string) string {
	if content == "" || tag == "" {
		return ""
	}
	return "\t<meta name=\"" + tag + "\"" +
		" content=" + "\"" + content + "\">\n"
}

// Generate common metatags
func (c *config) metatags() string {
	return metatag("description", fmStr("description", c.fm)) +
		metatag("keywords", fmStr("keywords", c.fm)) +
		metatag("robots", fmStr("robots", c.fm)) +
		metatag("author", fmStr("author", c.fm))
}

// titleTag turns front matter "title:" value into the
// all-important HTML <title> tag.
func (c *config) titleTag() string {
	title := fmStr("title", c.fm)
	if title == "" {
		return "\t<title>" + poweredBy + "</title>\n"
	} else {
		return "\t<title>" + title + "</title>\n"
	}
}

// PRINTY utilities

// quit displays a message fmt.Printf style and exits to the OS.
// That format string must be preceded by an exit code and an
// error object (nil if an error didn't occur).
func quit(exitCode int, err error, c *config, format string, ss ...interface{}) {
	msg := fmt.Sprint(fmtMsg(format, ss...))
	errmsg := ""
	if err != nil {
		errmsg = " " + err.Error()
	}
	if c == nil || c.currentFilename != "" {
		// Error exit.
		// Prints name of source file being processed.
		if exitCode != 0 {
			if c != nil {
				fmt.Printf("PocoCMS %s:\n \t%s%s\n", c.currentFile(), msg, errmsg)
			} else {
				fmt.Printf("PocoCMS %s\n \t%s\n", msg, errmsg)
			}
		}
	} else {
		// No c object available
		if err != nil {
			fmt.Printf("%s: %s\n", msg, errmsg)
		} else {
			fmt.Printf("%s\n", msg)
		}
	}
	os.Exit(exitCode)
}

// debug displays messages to stdout using Fprintf syntax.
// Same as print, but lets you search for debug
// in source code when it's meant to be in there
// temporarily.
// Differs from warn(), which sends its text to stderr
func debug(format string, ss ...interface{}) {
	fmt.Println(fmtMsg(format, ss...))
}

// print messages to stdout using Fprintf syntax.
// Same as debug, but meant to be left in the code.
// Differs from warn(), which sends its text to stderr
func print(format string, ss ...interface{}) {
	fmt.Println(fmtMsg(format, ss...))
}

// wait displays messages to stdout using Fprintf syntax.
// Waits for you to press a key, then Enter
// Continues if it's just Enter. press 'q' to quit.
func wait(format string, ss ...interface{}) {
	fmt.Println(fmtMsg(format, ss...))
	q := inputString()
	if len(q) >= 1 && strings.ToLower(q[0:1]) == "q" {
		quit(1, nil, nil, "Quitting")
	}
}

// warn displays messages to stderr using Fprintf syntax.
func warn(format string, ss ...interface{}) {
	msg := fmt.Sprintf(format, ss...)
	fmt.Fprintln(os.Stderr, msg)
}

// If the c.verbose flag is set, use the Printf style parameters
// to format the input and return a string.
func (c *config) verbose(format string, ss ...interface{}) {
	if c.verboseFlag {
		fmt.Println(fmtMsg(format, ss...))
	}
}

// fmtMsg() takes a list of strings like Fprintf, interpolates, and writes to a string
func fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}

// DEBUG UTILITIES/DUMP UTILITIES

// dumpSettings() lists config values
func (c *config) dumpSettings() {
	print("Global theme: %s", c.globalTheme.dir)
	print("Page theme: %s", c.pageTheme.dir)
	print("Markdown extensions: %v", c.markdownExtensions.list)
	print("Ignore: %v", c.skipPublish.list)
	print("Source directory: %s", c.root)
	print("Webroot directory: %s", c.webroot)
	print("Inline stylesheets: %v", !c.linkStyles)
	print("%s directory: %s", pocoDir, filepath.Join(executableDir(), pocoDir))
	print("Home page: %s", c.homePage)
}

// dumpFm Displays the contents of the page's front matter in JSON format
func dumpFm(c *config) string {
	var s string
	b, err := json.MarshalIndent(c.fm, "", "  ")
	if err != nil {
		return ("Error marshalling front matter")
	}
	s = string(b)
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.TrimSpace(s)
	return s
}

// PARSING UTILITIES

// convertMdYAMLFileToHTMLFragmentStr converts the Markdown code fragment,
// which may have front matter, to HTML.
// It doesn't execute templates.
// It doesn't assemble the file--that its, no header, footer, etc.
// So "hello, world" should come back as "<p>hello, world</p>"
// Returns parsed file as HTML.
func convertMdYAMLFileToHTMLFragmentStr(filename string, c *config) string {
	source := c.fileToString(filename)
	mdParser := newGoldmark()
	mdParserCtx := parser.NewContext()
	// Build a syntax tree (intermediate representation)
	// for the input Markdown text.
	_ = mdParser.Parser().Parse(text.NewReader([]byte(source)))
	var buf bytes.Buffer
	// Convert syntax tree to HTML and deposit in buf.Bytes().
	if err := mdParser.Convert([]byte(source), &buf, parser.WithContext(mdParserCtx)); err != nil {
		quit(1, err, c, "Unable to convert Markdown to HTML")
	}
	return string(buf.Bytes())
}

// mdYAMLFileToHTMLString converts a Markdown document
// with YAML front matter to HTML.
// The HTML file has not yet had templates executed,
// Destructive: replaces c.fm
// Returns a byte slice containing the HTML source.
func mdYAMLFileToHTMLString(c *config, filename string) (string, error) {
	source := fileToBuf(filename)
	var err error
	var HTML []byte
	if HTML, c.fm, err = mdYAMLToHTML(source); err != nil {
		return "", err
	} else {
		return string(HTML), nil
	}
}

// newGoldmark() allocates a Goldmark parser with a
// raft of other options.
func newGoldmark() goldmark.Markdown {
	exts := []goldmark.Extender{
		meta.New(
			meta.WithStoresInDocument(),
		),
		// Support GitHub tables & other extensions
		extension.Table,
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		extension.Linkify,
		// YouTube embedding
		ytembed.New(),
		highlighting.NewHighlighting(
			highlighting.WithStyle("autumn"),
			highlighting.WithFormatOptions()),
	}

	parserOpts := []parser.Option{
		parser.WithAttribute(),
		parser.WithAutoHeadingID()}

	renderOpts := []renderer.Option{
		html.WithUnsafe(),
		// html.WithHardWraps(),
		html.WithXHTML(),
	}
	return goldmark.New(
		goldmark.WithExtensions(exts...),
		goldmark.WithParserOptions(parserOpts...),
		goldmark.WithRendererOptions(renderOpts...),
	)
}

// mdYAMLStringToTemplatedHTMLString() takes raw HTML, converts to Markdown,
// and executes templates. Returns a string of the result.
// It doesn't assemble the file--that its, no header, footer, etc.
// So "hello, world" should come back as "<p>hello, world</p>"
// This should too:
// ---
// description: "world"
// ---
// hello, {{ .description }}.
//
// Returns parsed file as HTML.
// aside, etc.
//
// Replaces contents of c.fm.
func mdYAMLStringToTemplatedHTMLString(c *config, filename string, markdown string) string {
	var parsedHTML string
	var err error
	var b []byte
	if b, c.fm, err = mdYAMLToHTML([]byte(markdown)); err != nil {
		quit(1, err, c, "Unable to convert markdown to raw HTML")
	}
	if parsedHTML, err = doTemplate(filename, string(b), c); err != nil {
		quit(1, err, c, "%v: Problem executing template code", filename)
	}
	return parsedHTML
}

// mdYAMLtoHTML converts the Markdown file, which may
// have front matter, to HTML. The  front matter
// is one of the return values.
func mdYAMLToHTML(source []byte) ([]byte, map[string]interface{}, error) {

	mdParser := newGoldmark()
	mdParserCtx := parser.NewContext()

	document := mdParser.Parser().Parse(text.NewReader([]byte(source)))
	metaData := document.OwnerDocument().Meta()
	var buf bytes.Buffer
	// Convert Markdown source to HTML and deposit in buf.Bytes().
	if err := mdParser.Convert(source, &buf, parser.WithContext(mdParserCtx)); err != nil {
		return []byte{}, nil, err
	}
	// Obtain YAML front matter from document.
	return buf.Bytes(), metaData, nil
}

// PROMPT UTILITIES
// See also wait(), which prints a message and awaits 'q' or Enter.

// inputString() gets a string from the keyboard and returns it
// See also promptString()
func inputString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// promptString() displays a prompt, then awaits for keyboard
// input and returns it on completion.
// See also inputString(), promptYes()
func promptString(format string, ss ...interface{}) string {
	fmt.Print(fmtMsg(format, ss...))
	fmt.Print(" ")
	return inputString()
}

// promptStringDefault() displays a prompt, then awaits for keyboard
// input and returns it on completion. It precedes the end of the
// prompt with a default value in brackets.
// See also inputString(), promptYes()
func promptStringDefault(prompt string, defaultValue string) string {
	answer := promptString(prompt + " [" + defaultValue + "] ")
	if answer == "" {
		return defaultValue
	} else {
		return answer
	}
}

// promptYes() displays a prompt, then awaits
// keyboard input followed by Enter.
// Forces first letter of answer to lowercase.
// If the answer starts with 'y',
// returns true. Otherwise, loop until
// 'y' or 'n' is entered.
// See also inputString(), promptString()
func promptYes(prompt string) bool {
	// See also inputString(), promptYes()
	for {
		answer := promptString(prompt)
		if strings.HasPrefix(strings.ToLower(answer), "y") ||
			strings.HasPrefix(strings.ToLower(answer), "n") {
			return strings.HasPrefix(strings.ToLower(answer), "y")
		}
	}
}

// SERVER UTILITIES

// serve is the world's simplest web server, for quick tests
// only. c.port is a string like ":12345" and c.webroot is the
// pathname of the directory to serve static files from.
func (c *config) serve() {
	if portBusy(c.port) {
		print("Port %s is already in use", c.port)
		os.Exit(1)
	}
	if err := os.Chdir(c.webroot); err != nil {
		quit(1, err, c, "Unable to change to webroot directory %s", c.root)
	}
	// Simple static webserver:
	print("\n%s Web server running at http://localhost%s\nTo stop the web server, press Ctrl+C", theTime(), c.port)
	if err := http.ListenAndServe(c.port, http.FileServer(http.Dir(c.webroot))); err != nil {
		quit(1, err, c, "Error running web server")
	}
}

// theTime returns the current time as a string.
// Nothing configurable b/c it's just used for timestamping
// every page, a dumb diagnostic tool
func theTime() string {
	t := time.Now()
	s := fmt.Sprintf("%s", t.Format("02 Jan 2006 15:04:05"))
	return s
}

// portBusy() returns true if the port
// (in the form ":12345" is already in use.
func portBusy(port string) bool {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		return true
	}
	err = ln.Close()
	if err != nil {
		quit(1, err, nil, "Problem closing port")
	}

	return false
}

// MINIFY

// METADATA UTILITIES/FRONT MATTER UTILITIES

// fmStr is passed a front matter "type" and retrievs
// the value for the value passed in as key. Value
// is case-insensitive. So if
// c.globalFm has a theme named "pages/themes/foo"
// you'd pass in "theme" and get back "pages/themes/foo"
func fmStr(key string, fm map[string]interface{}) string {
	v := fm[strings.ToLower(key)]
	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

// fmStrSlice obtains a list of string values from the supplied front matter.
// For example, if you had this code in your Markdown file:
// ---
// Stylesheets
//   - 'https://cdn.jsdelivr.net/npm/holiday.css'
//   - 'fonts.css'
//
// ---

func fmStrSlice(key string, fm map[string]interface{}) []string {
	if key == "" {
		return []string{}
	}
	v, ok := fm[strings.ToLower(key)].([]interface{})
	if !ok {
		return []string{}
	}
	s := make([]string, len(v))
	for i, value := range v {
		r := fmt.Sprintf("%s", value)
		s[i] = r
	}
	return s
}

// newProject() takes a directory name and generates a
// site there.
// Pre: parseCommandLine()

// newProject is same as newSite(), except that
// newSite() can do its thing because it
// knows it's in a new dir. newProject sets
// up that condition and helps stop you from
// putting this project where it shouldn't be.
func (c *config) newProject(dir string) {
	c.root = dir
	if !dirExists(c.root) {
		// We're good. No directory by that name. Create it.
		err := os.MkdirAll(c.root, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			quit(1, err, c, "Unable to create new project directory %s", c.root)
		}
	} else {
		// Existing directory by that name. No need to continue.
		if isProject(c.root) {
			if promptYes("Project at " + c.root + "exists. This will replace index.md, add a " + pocoDir + " directory, and replace any directory named " + c.webroot + ". Continue?") {
				c.newSite()
			} else {
				quit(1, nil, c, "Leaving project at %s intact. Quitting.", c.root)
			}
		}
	}
}

// themeDirContents() returns a list of all installed themes
// separated by newlines
func (c *config) themeDirContents() string {
	files, err := ioutil.ReadDir(filepath.Join(pocoDir, "themes"))
	if err != nil {
		return ""
	}

	// Return value: the list of files
	var s string
	for _, f := range files {
		file := f.Name()
		s = s + fmt.Sprintf("%s\n", file)
	}
	return s
}

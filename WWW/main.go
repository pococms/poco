package main

// # Create a directory. It doesn't have to be here.
// mkdir ~/pococms
// # Navigate to that directory.
// cd ~/pococms
// # Clone the repo.
// git clone https://github.com/pococms/poco
// # The repo is now in ~/pococms/poco, so navigate there.
// cd poco
// And compile. There's only one file, so you can also use go run
// go build # OR go run main.go
//
// Example invocations
//
// If there's at least one Markdown file named index.md or README.md
// in the current directory, PocoCMS assumes it's a site and will
// create one. All you have to do is type the name of the program:
// poco
//
// Learn the command-line options:
// poco --help
//
//
// Use a style sheet from a CDN (you don't have to copy it to your project)
// poco --styles "https://cdn.jsdelivr.net/npm/holiday.css"
//
// Include the 2 css files shown
// poco --styles "theme.css light-mode.css"
//
// Compile only file template.md and deposit its HTML output in the same directory.
// Use the hack.css stylesheet from a CDN.
// Filename must come after all the options.
// ./poco --styles "https://unpkg.com/hack@0.8.1/dist/hack.css" template.md
// Use the docs subdirectory as the root of the site.
// ./poco --root "./docs"

// Get CSS file from CDN
// poco --styles "https://unpkg.com/spectre.css/dist/spectre.min.css"
// poco --styles "//writ.cmcenroe.me/1.0.4/writ.min.css" foo.md

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	//"reflect"
)

// Required begininng for a valid HTML document
var docType = `<!DOCTYPE html>
<html lang=`

// If a page lacks a title tag, it fails validation.
// Insert this if none is found.
var poweredBy = `Powered by PocoCMS`

// assemble takes the raw converted HTML in article,
// uses it to generate finished HTML document, and returns
// that document as a string.
func assemble(c *config, filename string, article string, language string, stylesheetList string) string {
	// This will contain the completed document as a string.
	htmlFile := ""
	// Execute templates. That way {{ .Title }} will be converted into
	// whatever frontMatter["Title"] is set to, etc.
	if parsedArticle, err := doTemplate("", article, c); err != nil {
		quit(1, err, c, "%v: Unable to execute ", filename)
	} else {
		article = parsedArticle
	}

	// If there are style tags in the theme's README,
	// add and enclose in a <style> tag. Otherwise
	// leave it empty.
	themeExtraTemplateTags := c.theme.styleFileTemplateTags
	if themeExtraTemplateTags != "" {
		themeExtraTemplateTags = "\t" + tagSurround("style", themeExtraTemplateTags, "\n")
	}
	// If there are style tags in the current file,
	// add and enclose in a <style> tag. Otherwise
	// leave it empty.
	extraStyleTags := c.styleTags()
	if extraStyleTags != "" {
		extraStyleTags = "\t" + tagSurround("style", extraStyleTags, "\n")
	}
	//debug("style tags: %v+\nextraStyleTags %v",c.styleTags, extraStyleTags)
	// Build the completed HTML document from the component pieces.
	htmlFile = docType + "\"" + language + "\">" + "\n" +
		"<head>\n" +
		"\t<meta charset=\"utf-8\">\n" +
		"\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
		titleTag(c) +
		metatags(c) +
		linktags(c) +
		stylesheets(stylesheetList, c) +
		themeExtraTemplateTags +
		extraStyleTags +
		"</head>\n<body>\n" +
		"<div id=\"page-container\">\n" +
		"<div id=\"content-wrap\">\n" +
		"\t" + layoutEl(c, "Header", filename) +
		"\t" + layoutEl(c, "Nav", filename) +
		"\t" + layoutEl(c, "Aside", filename) +
		"\t" + "<article id=\"article\">" + article + "\t" + "</article>" + "\n" +
		"</div><!-- content-wrap -->\n" +
		"\t" + layoutEl(c, "Footer", filename) +
		"</div><!-- page-container -->\n" +
		"</body>\n</html>\n"
	return htmlFile
} //   assemble

// THEME

// TODO: Create layoutEl() tests for
// - Nothing specified so theme header/footer/etc. are added
// - "SUPPRESS" specified when a theme is available but that element should not be displayed
// - Local file specified

// layoutEl() takes a layout element file named in the front matter
// and generates HTML, but it executes templates also.
//
// The layout element may also be a theme file.
//
// So, the priority order is:
// - If the front matter says "SUPRESS" in all caps then return empty string.
// - If there is a file named in the front matter, process and return its contents.
// - Otherwise, use a theme file.
//
// It can be a Markdown file, in which case no tags are needed,
// or an HTML file, in which the tags must be explicit.
// A layout element is one of the HTML tags such
// as header, nav, aside, article, and a few others
// For more info on layout elements see:
// https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML/Document_and_website_structure#html_layout_elements_in_more_detail
// The easiest way is to use markdown.
// Fore example, suppose you have a header file named mdhead.md and
// it contains only the following:
//
// hello, world.
//
// The genereated HTML would be "<p><header>hello, world.</header></p>"

// For example, suppose you have a header file named head.html. It
// would be named in the front matter like this:
// ---
// Header: head.html
// ---
//
// The layout element file is expected to be a complete tag. For example,
// the header file could be as simple as this:
//
//	<header>hello, world.</header>
//
// This function would read in the head.html file (or whatever
// the file was named in the front matter) and insert it before the
// body of the document.
//
// fm contains the YAML front matter.
// element is the file containing the layout element, for example, head.html.
// If element  ends in ".html" it must be a complete header tag, with
// both tags included. If element doesn't end in ".html" it is considered
// to be a Markdown file and is processed that way.
// sourcefile is the fully qualified pathname of the .md file being processed
// TODO: Code smell
func layoutEl(c *config, element string, sourcefile string) string {
	// element looks like "Header", "Footer", etc. because front matter key is capitalized.
	// Force to lowercase for use as an HTML tag.
	tag := strings.ToLower(element)

	// Get the filename for this layout element. For example,
	// if the front matter said Header: "foo.md" this would
	// return "foo.md".

	// Special case: if there's a theme using this element
	// you can suppress its output by using the special value
	// "SUPPRESS" after Header:, Nav:, Aside: or Footer: in
	// the front matter, e.g. Header: "SUPPRESS"
	filename := frontMatterStr(element, c)
	// debug("\tlayoutEl %s: %s", element, filename)
	//debug("\t%s %s (%s):  %+v", filepath.Base(sourcefile),element,filename,c.fm)

	if filename == "SUPPRESS" {
		return ""
	}

	// If no filename, then use the theme layout element, if any.
	if filename == "" {
		// No layout element specified in front matter.
		// See if there's a theme and if it has that layout element.
		// Convert to HTML and executetemplate.
		// so just return it.
		s := c.themeEl(tag)
		return s
	}
	// A filename was specified
	// Takes priority over theme.

	isMarkdown := false

	// Full path to layout element file. So the file 'layout/myheader.md'
	// woud be transformed into /Users/tom/mysite/layout/myheader.md'
	// or something similar.
	fullPath := ""

	// Get the name of the file. For example, the front
	// matter my say Header: myheader.md so
	// layoutElSource is 'myheader.md'

	layoutElSource := frontMatterStr(element, c)
	//debug("\t%s layoutEl %s", c.currentFilename, layoutElSource)
	if filepath.IsAbs(layoutElSource) {
		fullPath = layoutElSource
	} else {
		var err error
		var rel string
		// TODO: Cache current directory
		if rel, err = filepath.Rel(currDir(), sourcefile); err != nil {
			quit(1, nil, c, "Error calling filepath.Rel(%s,%s)", currDir(), sourcefile)
		}
		rel = filepath.Dir(rel)
		fullPath = filepath.Join(currDir(), rel, layoutElSource)
	}
	if filepath.Ext(fullPath) != ".html" {
		isMarkdown = true
	}

	parsedArticle := ""
	raw := ""
	var err error
	if isMarkdown {
		if !fileExists(fullPath) {
			quit(1, nil, c, "Front matter \"%s:\" specified file %s but can't find it", element, fullPath)
		}
		// TODO: This feels wasteful
		raw = convertMdYAMLFileToHTMLStr(fullPath, c)
		if parsedArticle, err = doTemplate("", raw, c); err != nil {
			quit(1, err, c, "%v: Unable to execute ", filename)
		}
		if parsedArticle != "" {
			//debug("\t\tFront matter is now: %+v", c.fm)
			wholeTag := "<" + tag + ">" + parsedArticle + "</" + tag + ">\n"
			return wholeTag
		}
		return ""
	}
	return c.fileToString(fullPath)

}

// loadTheme tries to find the named theme directory
// and load its files into c.theme
// Tests:
// - Missing README.md
// - Missing Stylesheets, StyleFileTemplates
// - Missing LICENSE file
func (c *config) loadTheme() {
	nc := getFrontMatter(c.homePage)
	themeDir := frontMatterStr("Theme", nc)
	if themeDir == "" {
		return
	}
	if !fileExists(c.homePage) {
		quit(1, nil, c, "Theme %s is missing a README.md", themeDir)
	}
	// Make sure there's a LICENSE file
	license := filepath.Join(themeDir, "LICENSE")
	if c.fileToString(license) == "" {
		quit(1, nil, c, "%s theme is missing a LICENSE file", themeDir)
	}
	c.theme.dir = themeDir
	header := filepath.Join(themeDir, "header.md")
	if fileExists(header) {
		c.theme.header = c.fileToString(header)
	}

	nav := filepath.Join(themeDir, "nav.md")
	if fileExists(nav) {
		c.theme.nav = c.fileToString(nav)
	}

	aside := filepath.Join(themeDir, "aside.md")
	if fileExists(aside) {
		c.theme.aside = c.fileToString(aside)
	}

	footer := filepath.Join(themeDir, "footer.md")
	if fileExists(footer) {
		c.theme.footer = c.fileToString(footer)
	}
	// Obtain the front matter from the README.md
	// (inside a dummy config object)
	themeReadMe := filepath.Join(themeDir, "README.md")
	nc = getFrontMatter(themeReadMe)

	// Get the list of style sheets required for this theme.
	// Remember that stylesheets not in this list won't
	// be copid in. This is different from the non-theme
	// behavior of just copying all stylesheets
	// to the webroot.

	// Read each stylesheet into a string, then appened
	// it into the theme file's styleFilesEmbedded
	// member. It will then be injected into the
	// HTML file directly, in order requested.
	styleFileList := frontMatterStrSlice("Stylesheets", nc)
	// nc.theme.dir = themeDir
	c.styleFiles(styleFileList)
	// Theme loaded. Now get additional style tags.
	c.styleTags()
	/*
	  debug("Theme: %v\nHeader: %v\nNav: %v\nAside: %v\nFooter: %v. styleGilesEmbedded: %v\n\n",
	      themeDir,
	      c.theme.header,
	      c.theme.nav,
	      c.theme.aside,
	      c.theme.footer,
	      c.theme.styleFilesEmbedded)
	*/

}

// Pre: c.theme.dir must know theme directory.
func (c *config) styleFiles(styleFileList []string) {
	//  Contents of header, nav, etc. ready to be converted from Markdown to HTML
	var s string
	for _, filename := range styleFileList {

		// Handle case of URLs as opposed to local file
		if strings.HasPrefix(filename, "http") {
			// TODO: Check for redirect?
			// https://golangdocs.com/golang-download-files
			s = c.downloadTextFile(filename)
			//

		} else if !filepath.IsAbs(filename) {
			// Handle case of local file with relative path
			fullPath := filepath.Join(c.theme.dir, filename)
			s = c.fileToString(fullPath)

		} else {
			// Handle case of local file with absolute path
			s = c.fileToString(filename)
		}
		c.theme.styleFilesEmbedded = c.theme.styleFilesEmbedded + s
	}
}

// themeEl() returns the theme layout element (header,nav
// aside, footer). Remember: this is in the case where
// no header/nav/aside/footer was specified in the Markdown
// source file's front matter. This extracts any
// such element.
func (c *config) themeEl(tag string) string {
	switch tag {
	case "header":
		if c.theme.header != "" {
			s := mdYAMLStringToTemplatedHTMLString(c, c.theme.header)
			return tagSurround(tag, s, "\n")
		}
	case "nav":
		if c.theme.nav != "" {
			s := mdYAMLStringToTemplatedHTMLString(c, c.theme.nav)
			return tagSurround(tag, s, "\n")
		}
	case "aside":
		if c.theme.aside != "" {
			s := mdYAMLStringToTemplatedHTMLString(c, c.theme.aside)
			return tagSurround(tag, s, "\n")
		}
	case "footer":
		if c.theme.footer != "" {
			s := mdYAMLStringToTemplatedHTMLString(c, c.theme.footer)
			return tagSurround(tag, s, "\n")
		}
	}
	return ""
}

// HTML UTILITIES

// tagSurround takes text and surrounds it with
// opening and closing tags, so
// tagSurround("header","WELCOME","\n") returns "<header>WELCOME</header>\n"
// You can optionally include text after, because sometimes it
// makes sense to include a newline after the closing tag.
func tagSurround(tag string, txt string, extra ...string) string {
	// TODO: Bit of a kludge. Point is I'm getting a newline at the
	// end of txt and that's what I should be focusing on.
	// It was creating tags like <header>hello\n<header>
	txt = strings.TrimSpace(txt)
	return "<" + tag + ">" + txt + "</" + tag + ">" + extra[0]
}

// sliceToStylesheetStr takes a slice of simple stylesheet names, such as
// [ "foo.css", "bar.css" ] and converts it into a string
// consisting of stylesheet link tags separated by newlines:
//
// <link rel="stylesheet" href="foo.css"/>
// <link rel="stylesheet" href="bar.css"/>
func sliceToStylesheetStr(sheets []string) string {
	var tags string
	for _, sheet := range sheets {
		tag := fmt.Sprintf("\t<link rel=\"stylesheet\" href=\"%s\">\n", sheet)
		tags += tag
	}
	return tags
}

// StyleTags takes a list of tags and inserts them into right before the
// closing head tag, so they can override anything that came before.
// These are literal tags, not filenames.
// They're listed under "StyleTags" in the front matter
// Returns them as a string. For clarity each tag is indented
// and ends with a newline.
// Example:
//
// StyleTags:
//   - "h1{color:blue;}"
//   - "p{color:darkgray;}"
//
// Would yield:
//
//	"{color:blue;}\n\t\tp{color:darkgray;}\n"
func (c *config) styleTags() string {
	tagSlice := frontMatterStrSlice("StyleTags", c)
	if tagSlice == nil {
		return ""
	}
	// Return value
	tags := ""
	for _, value := range tagSlice {
		s := fmt.Sprintf("\t\t%s\n", value)
		tags = tags + s
	}
	return tags
}

// stylesheets() takes stylesheets listed on the command line
// e.g. --styles "foo.css bar.css", and adds them to
// the head. It then generates stylesheet tags for the ones listed in
// the front matter.
// Those listed in the front matter are appended, so they take
// priority.
func stylesheets(sheets string, c *config) string {
	var globalSlice []string
	var globals string

	// Handle case of theme specified
	// This is how you tell if a theme is present
	if c.theme.dir != "" && frontMatterStr("Theme", c) != "SUPPRESS" {
		// TODO: minify these mofos
		return "<!-- EMBEDDED STYLE --><style>" + c.theme.styleFilesEmbedded + "</style>\n"
	}

	if sheets != "" {
		// Build a string from stylesheets named on the command line.
		globalSlice = strings.Split(sheets, " ")
		globals = sliceToStylesheetStr(globalSlice)
	}
	// Build a string from stylesheets named in the
	// StyleFileTemplates: front matter for the home page
	//templates := ""
	// Build a string from stylesheets named in the
	// Stylesheets: front matter for this page
	localSlice := frontMatterStrSlice("Stylesheets", c)
	locals := sliceToStylesheetStr(localSlice)

	// Stylesheets named in the front matter takes priority,
	// so they goes last. This allows you to have stylesheets
	// on the command line that act as templates, but that
	// you can override using stylesheets named in
	// the front matter.
	return globals + c.styleFileTemplates + locals
}

// The --verbose flag. It shows progress as the site is created.
// Required by the Verbose() function.
var gVerbose bool

// theme contains all the (lightweight) files needed for a theme:
// header.md, style sheets, etc.
type theme struct {

	// READ ONLY: Full pathname to theme directory
	dir string

	// Contents of LICENSE file. Can't be empty
	license string

	// Who created it, natch
	author string

	// Name for the theme with spaces and other characters allowed.
	// If the directory name is my-great-theme you might
	// want this to be "My Great! Theme"
	branding string

	// One or more sentences selling the theme.
	description string

	// Contents of header.md
	header string

	// Contents of nav.md
	nav string

	// Contents of aside.md
	aside string

	// Contents of footer.md
	footer string

	// Names of stylesheets
	styleFiles []string

	// Names of template stylesheets
	styleFileTemplates string

	// Names of style tags FROM THE CURRENT MARKDOWN FILE,
	// not the theme's README.md.
	// Scenario: You've developed a light theme.
	// You want to experiment with a dark theme.
	// So you add IN THE CURRENT MARKDOWN FILE
	// StyleTags:
	// - "article{background-color:black;color:black}"
	styleFileTemplateTags string

	// The stylesheets for each theme are concantenated, then read
	// into this string. It's injected straight into the HTML for
	// each file using this theme.
	styleFilesEmbedded string
}

// there are no configuration files (yet) but this holds
// configuration info for the project, for example, template
// stylesheets and current file being processed.
// That stuff lives in the front matter of the home
// page (first checks for README.md, then checks for index.md)
type config struct {
	// Front matter
	fm map[string]interface{}

	// Name of Markdown file being processed
	currentFilename string

	// List of all files being processed
	files []string

	// This is true only when a home page
	// (root of the directoryh tree) README.md
	// or index.md is being processed
	// TODO: Probably unnecessary
	hitHomePage bool

	// Full pathname of the root index file Markdown in the root directory.
	// If present, it's either "README.md" or "index.md"
	homePage string

	// Home directory for source code
	root string

	// List of stylesheets to apply to every page in
	// string form, ready to drop into the
	// <head>
	styleFileTemplates string

	// Contents of a theme directory
	theme theme

	// Output directory for published files
	webroot string
}

// findHomePage() returns the source file used for the root
// index page in the root directory. Since README.md is
// commonly used, it takes priority. Next priority is index.md
func (c *config) findHomePage() {
	if c.root == "." || c.root == "" {
		c.root = currDir()
	}
	// Look for "README.md" or "index.md" in that order.
	// return "" if neither found.
	c.homePage = indexFile(c.root)

	if c.homePage != "" {
		return
	}

	if !dirEmpty(c.root) {
		// No home page.
		// Directory has files.
		// User may not wish to create a new project.
		if promptYes("Create a home page?") {
			writeDefaultHomePage(c, c.root)
		} else {
			quit(1, nil, c, "Can't build a project with an index.md or README.md")
		}
	} else {
		// Empty dir, so create home page
		writeDefaultHomePage(c, c.root)
	}
}

// setup() Obtains README.md or index.md.
// Reads in the front matter to get its config information.
// Sets values accordingly.
func (c *config) setup() {
	c.findHomePage()
	// If a theme directory was named in front matter's Theme: key,
	// read it in.
	c.loadTheme()
}

// initConfig reads the home page and gets
// sitewide configuration info.
func initConfig() *config {
	config := config{}
	return &config
}

// initConfig()

func main() {
	c := initConfig()
	// cleanup determines whether or not the publish (aka WWW) directory
	// gets deleted on start.
	var cleanup bool
	flag.BoolVar(&cleanup, "cleanup", true, "Delete publish directory before converting files")

	// debugFrontmatter command-line option shows the front matter of each page
	var debugFrontMatter bool
	flag.BoolVar(&debugFrontMatter, "debug-frontmatter", false, "Shows the front matter of each page")

	// skip lets you skip the named files from being processed
	var skip string
	flag.StringVar(&skip, "skip", "node_modules .git .DS_Store .gitignore", "List of files to skip when generating a site")

	// language sets HTML lang= value, such as <html lang="fr">
	var language string
	flag.StringVar(&language, "language", "en", "HTML language designation, such as en or fr")

	//var root string
	flag.StringVar(&c.root, "root", ".", "Starting directory of the project")

	// List of stylesheets to include on each page.
	var stylesheets string
	flag.StringVar(&stylesheets, "styles", "", "One or more stylesheets (use quotes if more than one)")

	// Title tag.
	var title string
	flag.StringVar(&title, "Title", poweredBy, "Contents of the HTML title tag")

	// Verbose shows progress as site is generated.
	flag.BoolVar(&gVerbose, "verbose", false, "Display information about project as it's generated")

	// webroot is the directory used to house the final generated website.
	flag.StringVar(&c.webroot, "webroot", "WWW", "Subdirectory used for generated HTML files")

	// Process command line flags such as --verbose, --title and so on.
	flag.Parse()

	// Collect configuration info for this project

	// See if a source file was specified. Otherwise the whole directory
	// and nested subdirectories are processed.
	c.currentFilename = flag.Arg(0)

	var err error
	if c.root, err = filepath.Abs(c.root); err != nil {
		quit(1, err, c, "Unable to absolute path of %s", c.root)
	}
	if c.currentFilename != "" {
		// Something's left on the command line. It's presumed to
		// be a filename.
		if !fileExists(c.currentFilename) {
			quit(1, nil, c, "Can't find the file %v", c.currentFilename)
		} else {
			// Special case: if you name a file on the command line, it will
			// generate an HTML document from that file and pass you the new filename.
			// The output file isn't published to webroot. It's published to the
			// current directory.
			htmlFilename := buildFileToFile(c, c.currentFilename, stylesheets, language, debugFrontMatter)
			quit(0, nil, c, "Built file %s", htmlFilename)
		}
	}
	// No file was given on the command line.
	// Build the project in place.

	// Obtain README.md or index.md.
	// Read in the front matter to get its config information.
	// Set values accordingly.
	c.setup()

	// markdownExtensions are how PocoCMS figures out whether
	// a file is Markdown. If it ends in any one of these then
	// it gets converted to HTML.
	var markdownExtensions searchInfo
	markdownExtensions.list = []string{".md", ".mkd", ".mdwn", ".mdown", ".mdtxt", ".mdtext", ".markdown"}

	webrootPath := buildSite(c, c.webroot, skip, markdownExtensions, language, stylesheets, cleanup, debugFrontMatter)
	debug("End of main() after buildSite(). webrootpath: %s. c.webroot: %s", webrootPath, c.webroot)
  debug("%v",c.files)
	quit(0, nil, c, "Site published to %s", filepath.Join(webrootPath, "index.html"))

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
func buildFileToFile(c *config, filename string, stylesheets string, language string, debugFrontMatter bool) (outfile string) {
	// Convert Markdown file filename to raw HTML, then assemble a complete HTML document to be published.
	// Return the document as a string.
	html, htmlFilename := buildFileToTemplatedString(c, filename, stylesheets, language)
	// Write the contents of the completed HTML document to a file.
	writeStringToFile(c, htmlFilename, html)
	// Return the name of the converted file
	return htmlFilename
}

// buildFileToTemplatedString converts a file from Markdown to raw HTML,
// pulls in everything required to create a complete HTML document,
// executes templates,
// generates an output file,
// and returns name of the destination HTML file
// Does not check if the input file is Markdown.
// TODO: Ideally this would be called from buildSite()
// Reeturns the string and the filenlame
func buildFileToTemplatedString(c *config, filename string, stylesheets string, language string) (string, string) {
	// Exit silently if not a valid file
	if filename == "" || !fileExists(filename) {
		return "", ""
	}
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
		finishedDocument := assemble(c, c.currentFilename, rawHTML, language, stylesheets)
		// Return the finished document and its filename
		return finishedDocument, dest
	}
}

// converts all files (except those in skipPublish.List) to HTML,
// and deposits them in webroot. Attempts to create webroot if it
// doesn't exist. webroot is expected to be a subdirectory of
// projectDir.
// Return name of the root directory files are published to
func buildSite(c *config, webroot string, skip string, markdownExtensions searchInfo, language string, stylesheets string, cleanup bool, debugFrontMatter bool) string {

	var err error
	// Make sure it's a valid site. If not, create a minimal home page.

	// Change to requested directory
	if err = os.Chdir(c.root); err != nil {
		quit(1, err, c, "Unable to change to directory %s", c.root)
	}

	// Cache project's root directory
	var homeDir string
	if homeDir, err = os.Getwd(); err != nil {
		quit(1, err, c, "Unable to get name of current directory")
	}

	// Delete web root directory unless otherwise requested
	if cleanup {
		delDir := filepath.Join(homeDir, webroot)
		Verbose("Deleting directory %v", delDir)
		if err := os.RemoveAll(delDir); err != nil {
			quit(1, err, c, "Unable to delete publish directory %v", delDir)
		}
	}

	// Convert the list of exclusions into a string slice.
	// skipPublish = getSkipPublish()
	var skipPublish searchInfo
	skipPublish.list = strings.Split(skip, " ")

	// Collect all the files required for this project.
	c.files, err = getProjectTree(".", skipPublish)
	if err != nil {
		quit(1, err, c, "Unable to get directory tree")
	}

	// Full pathname of file to copy to target directory
	var source string

	// Full pathname of output directory for copied files
	var target string

	// After Markdown file is converted to HTML, it ends up in this string.
	// and eventually
	var HTML string

	// Relative directory of file. Required to determine where
	// to copy target file.
	var rel string

	// true if it was converted to HTML.
	// false if it's not a Markdown file, which means it will be copied
	// unchanged to the output directory
	var converted bool

	// Name of directory used to publish output files
	var webrootPath string

	// Main loop. Traverse the list of files to be copied.
	// If a file is Markdown as determined by its file extension,
	// convert to HTML and copy to output directory.
	// If a file isn't Markdown, copy to output directory with
	// no processing.
	for _, filename := range c.files {

		// true if it's  Markdown file converted to HTML
		converted = false

		// Get the fully qualified pathname for this file.
		c.currentFilename = filepath.Join(homeDir, filename)

		Verbose("%s", c.currentFilename)
		// Separate out the file's origin directory
		sourceDir := filepath.Dir(c.currentFilename)

		// Get the relatve directory. For example, if your directory
		// is ~/raj/blog and you're in ~/raj/blog/2023/may, then
		// the relative directory is 2023/may.
		if rel, err = filepath.Rel(homeDir, sourceDir); err != nil {
			quit(1, err, c, "Unable to get relative paths of %s and %s", homeDir, sourceDir)
		}
		debug("filepath.Rel(%s,%s) == %s", homeDir, sourceDir, rel)
		// Determine the destination directory.
		// xxx
    webrootPath = filepath.Join(homeDir, webroot, rel)
		// Obtain file extension.
		ext := path.Ext(filename)
		// Replace converted filename extension, from markdown to HTML.
		// Only convert to HTML if it has a Markdown extension.
		if markdownExtensions.Found(ext) {
			// Convert the Markdown file to an HTML string
			if HTML, err = mdYAMLFileToHTMLString(c, filename); err != nil {
				quit(1, err, c, "Error converting Markdown file to HTML")
			}
			// If asked, display the front matter
			if debugFrontMatter {
				debug("TODO: dumpFrontMatter() TODO not hit in 1 file situation")
				debug(dumpFrontMatter(c))
			}
			source = filename[0:len(filename)-len(ext)] + ".html"
			converted = true
		} else {
			// Not a Markdown file. Copy unchanged.
			source = filename
			// Insert destination (webroot) directory
			converted = false
		}
		target = filepath.Join(webrootPath, filepath.Base(source))

		// Create the target directory for this file if it
		// doesn't exist.
		if !dirExists(webrootPath) {
			err := os.MkdirAll(webrootPath, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				quit(1, err, c, "Unable to create directory %s", webrootPath)
			}
		}

		// Now have list of all files in directory tree.
		// If markdown, convert to HTML and copy that file to the HTML publication directory.
		// If not, copy to target publication directory unchanged.

		if converted {
			// Take the raw converted HTML and use it to generate a complete HTML document in a string
			h := assemble(c, c.currentFilename, HTML, language, stylesheets)
			writeStringToFile(c, target, h)
		} else {
			copyFile(c, source, target)
		}
	}
	// This is where the files were published
	ensureIndexHTML(c.webroot, c)
	// Display all files, Markdown or not, that were processed
	debug("returning webrootPath: %s. c.webroot: %s", webrootPath, c.webroot)
	return webrootPath
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

func getSkipPublish() []string {
	// var skipPublish searchInfo
	// skipPublish.list = strings.Split(skip, " ")
	// var skipPublish searchInfo
	// skipPublish.list = strings.Split(skip, " ")
	// skipPublish = getSkipPublish()
	return []string{}

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
// If it has both README.md,  takes priority.
func indexFile(path string) string {
	var indexMd, readmeMd string
	indexMd = filepath.Join(path, "index.md")
	if fileExists(indexMd) {
		return indexMd
	}
	readmeMd = filepath.Join(path, "README.md")
	if fileExists(readmeMd) {
		return readmeMd
	}
	return ""
}

// SYSTEM UTILITIES
// curDir() returns the current directory name.
func currDir() string {
	// if path, err := os.Executable(); err != nil {
	if path, err := os.Getwd(); err != nil {
		return "unknown directory"
	} else {
		return path
	}
}

// FILE UTILITIES
// copyFile, well, does just that. Doesnt' return errors.
func copyFile(c *config, source string, target string) {
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

// Generates a simple home page as an HTML string
// Uses the file segment of dir as the the H1 title.
// Uses current directory if "." or "" are passed
func defaultHomePage(dir string) string {

	var indexMdFront = `---
Stylesheets:
    - https://unpkg.com/simpledotcss/simple.min.css
---
`
	var indexMdBody = `
hello, world.

Learn more at [PocoCMS tutorials](https://pococms.com/docs/tutorials.html) 
`
	if dir == "" || dir == "." {
		dir, _ = os.Getwd()
	}
	h1 := filepath.Base(dir)
	page := indexMdFront +
		"# " + h1 + "\n" +
		indexMdBody
	return page
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

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false // Either not empty or error, suits both cases
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

// Generates a simple home page
// and writes it to index.md in dir. Uses the file
// segment of dir as the the H1 title.
// Returns the full pathname of the file.
func writeDefaultHomePage(c *config, dir string) string {
	html := defaultHomePage(dir)
	pathname := filepath.Join(dir, "index.md")
	writeStringToFile(c, pathname, html)
	c.homePage = pathname
	return pathname
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
		quit(1, err, c, "Unable to convert file %s into a string", filename)
	}
	return string(input)
}

// replaceExtension() is passed a filename and returns a filename
// with the specified extension.
func replaceExtension(filename string, newExtension string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename)) + "." + newExtension
}

// writeStringToFile creates a file called filename without checking to see if it
// exists, then writes contents to it.
// filename is afully qualified pathname.
// contents is the string to write
func writeStringToFile(c *config, filename, contents string) {
	var out *os.File
	var err error
	if out, err = os.Create(filename); err != nil {
		quit(1, err, c, "writeStringToFile: Unable to create file %s", filename)
	}
	if _, err = out.WriteString(contents); err != nil {
		// TODO: Renumber error code?
		quit(1, err, c, "Error writing to file %s", filename)
	}
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

		// Obtain just the filename.
		name := filepath.Base(info.Name())

		// Skip any directory to be excluded, such as
		// the pub and .git directores
		if skipPublish.Found(name) && isDir {
			return filepath.SkipDir
		}
		// It may be just a filename on the exclude list.
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

// Generate <link> tags

// frontMatterStr obtains a string value from the front matter. For example,
// if you had this code in your Markdown file:
// ---
// Title: yo mama
// ---
// I like {{ .Title }}
//
// It would render like this in the HTML:
// I like yo mama
func frontMatterStr(key string, c *config) string {
	v := c.fm[key]
	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

// frontMatterStrSlice obtains a list of string values from the front matter.
// For example, if you had this code in your Markdown file:
// ---
// Stylesheets
//   - 'https://cdn.jsdelivr.net/npm/holiday.css'
//   - 'fonts.css'
//
// ---
//
// I like yo mama
func frontMatterStrSlice(key string, c *config) []string {
	if key == "" {
		return []string{}
	}
	v, ok := c.fm[key].([]interface{})
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

// frontMatterStrSliceStr obtains the front matter value at
// key, which is presumed to be a string array/slice.
// Returns these values concatenated into a string
// (each string gets a newline appended for clarity)
func frontMatterStrSliceStr(key string, c *config) string {

	// Return empty string if no key provided.
	if key == "" {
		return ""
	}

	// Return empty string if requested key has no value
	// associated with it.
	v, ok := c.fm[key].([]interface{})
	if !ok {
		return ""
	}
	//s := make([]string, len(v))
	var tags string
	for _, value := range v {
		tag := fmt.Sprintf("%s\n", value)
		tags = tags + tag
	}
	return tags
}

// linkTags() obtains the list of link tags from the "LinkTags" front matter
// and inserts them into the document.
func linktags(c *config) string {
	linkTags := frontMatterStrSlice("LinkTags", c)
	if len(linkTags) < 1 {
		return ""
	}
	tags := ""
	for _, tag := range linkTags {
		tags += "\t" + tag + "\n"
	}
	return tags
}

func titleTag(c *config) string {
	title := frontMatterStr("Title", c)
	if title == "" {
		return "\t<title>" + poweredBy + "</title>\n"
	} else {
		return "\t<title>" + title + "</title>\n"
	}
}

// Generate common metatags
func metatags(c *config) string {
	return metatag("description", frontMatterStr("Description", c)) +
		metatag("keywords", frontMatterStr("Keywords", c)) +
		metatag("author", frontMatterStr("Author", c))
}

// metatag() generates a metatag such as <meta name="description"content="PocoCMS: Markdown-based CMS in 1 file, written in Go">
func metatag(tag string, content string) string {
	if content == "" || tag == "" {
		return ""
	}
	return "\t<meta name=\"" + tag + "\"" +
		" content=" + "\"" + content + "\">\n"
}

// PRINTY utilities

// If the Verbose flag is set, use the Printf style parameters
// to format the input and return a string.
func Verbose(format string, ss ...interface{}) {
	if gVerbose {
		fmt.Println(fmtMsg(format, ss...))
	}
}

// quit displays a message fmt.Printf style and exits to the OS.
// That format string must be preceded by an exit code and an
// error object (nil if an error didn't occur).
func quit(exitCode int, err error, c *config, format string, ss ...interface{}) {
	msg := fmt.Sprint(fmtMsg(format, ss...))
	errmsg := ""
	if err != nil {
		errmsg = " " + err.Error()
	}
	// fmt.Println(msg + errmsg)
	if c.currentFilename != "" {
		// Error exit
		if exitCode != 0 {
			fmt.Printf("PocoCMS %s:\n \t%s%s\n", c.currentFilename, msg, errmsg)
		} else {
			fmt.Printf("%s%s\n", msg, errmsg)
		}
	}
	os.Exit(exitCode)
}

// debug displays messages to stdout using Fprintf syntax.
// A little list printing and easier to search
func debug(format string, ss ...interface{}) {
	fmt.Println(fmtMsg(format, ss...))
}

// warn displays messages to stderr using Fprintf syntax.
func warn(format string, ss ...interface{}) {
	msg := fmt.Sprintf(format, ss...)
	fmt.Fprintln(os.Stderr, msg)
}

// fmtMsg() takes a list of strings like Fprintf, interpolates, and writes to a string
func fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}

// DEBUG UTILITIES

// dumpFrontMatter Displays the contents of the page's front matter in JSON format
func dumpFrontMatter(c *config) string {
	b, err := json.MarshalIndent(c.fm, "", "  ")
	s := string(b)
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, "\"", "")
	if err == nil {
		return s
	}
	return err.Error()
}

// PARSING UTILITIES

// convertMdYAMLToHTML converts the Markdown file, which may
// have front matter, to HTML. The front matter is passed in
// is used but not written to
func convertMdYAMLFileToHTMLStr(filename string, c *config) string {
	source := c.fileToString(filename)
	mdParser := newGoldmark()
	mdParserCtx := parser.NewContext()

	_ = mdParser.Parser().Parse(text.NewReader([]byte(source)))
	//metaData := document.OwnerDocument().Meta()
	var buf bytes.Buffer
	// Convert Markdown source to HTML and deposit in buf.Bytes().
	if err := mdParser.Convert([]byte(source), &buf, parser.WithContext(mdParserCtx)); err != nil {
		quit(1, err, c, "Unable to convert Markdown to HTML")
	}
	// Obtain YAML front matter from document.
	return string(buf.Bytes())
}

// getFrontMatter() takes a file, typically the README.md
// for a theme, and extracts its front matter. It does all
// the usual template execution, etc. It discards the
// generated HTML and returns a dummy config object with
// the front matter from the file in nc.fm.
func getFrontMatter(filename string) (newConfig *config) {
	nc := initConfig()
	var rawHTML string
	var err error

	// Convert a Markdown file, possibly with front matter, to HTML
	if rawHTML, err = mdYAMLFileToHTMLString(nc, filename); err != nil {
		quit(1, err, nc, "%v: convert %s to HTML", filename)
	}

	// Execute its templates.
	if _, err = doTemplate("", rawHTML, nc); err != nil {
		quit(1, err, nc, "%v: Unable to execute ", filename)
	}

	// And return a new config object with the front matter ready to go.
	return nc
}

// mdYAMLFiletoHTML converts a Markdown document
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
		highlighting.NewHighlighting(
			highlighting.WithStyle("github"),
			highlighting.WithFormatOptions()),
	}

	parserOpts := []parser.Option{
		parser.WithAttribute(),
		parser.WithAutoHeadingID()}

	renderOpts := []renderer.Option{
		// html.WithUnsafe(),
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
// Replaces c.fm
func mdYAMLStringToTemplatedHTMLString(c *config, markdown string) string {
	var parsedHTML string
	var err error
	var b []byte
	if b, c.fm, err = mdYAMLToHTML([]byte(markdown)); err != nil {
		quit(1, err, c, "Unable to convert markdown to raw HTML")
	}
	if parsedHTML, err = doTemplate(filepath.Base(c.currentFilename), string(b), c); err != nil {
		quit(1, err, c, "%v: Problem executing template code", c.currentFilename)
	}
	return parsedHTML
}

// mdYAMLtoHTML converts the Markdown file, which may
// have front matter, to HTML. The  front matter
// is deposited in frontMatter.
func mdYAMLToHTML(source []byte) ([]byte, map[string]interface{}, error) {

	// TODO: Does this obviate the need of some of the othe routines?
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
	// TODO Remove if not used
	fmt.Print(prompt + " [" + defaultValue + "] ")
	answer = inputString()
	if answer == "" {
		return defaultValue
	} else {
		return answer
	}
}

// promptYes() displays a prompt, then awaits
// keyboard input. If the answer starts with Y,
// returns true. Otherwise, returns false.
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
	///return strings.HasPrefix(strings.ToLower(answer), "y")
}

// PocoCMS: Markdown-based CMS in 1 file, written in Go
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
// Learn the command-line options:
// poco --help
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
// ./pococms --root "./docs"

// Get CSS file from CDN
// poco --styles "https://unpkg.com/spectre.css/dist/spectre.min.css"
// poco --styles "//writ.cmcenroe.me/1.0.4/writ.min.css" foo.md

import (
	"bytes"
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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// If you invoke poco without a filename, it creates an index file from this
// and publishes it. It also works as an informal test harness.

var indexSample = `---
Title: Try 4:37am
---
# Welcome to PocoCMS
---
Title: {{ .Title }}
`

var OLDindexSample = `---
Title: 'inserttitle'
Description: PocoCMS: Markdown-based CMS in 1 file, written in Go
Author: 'Tom Campbell'
Header: header.html
Nav: nav.html
Footer: footer.html
LinkTags:
    - <link rel="icon" href="favicon.ico">
    - <link rel="preconnect" href="https://fonts.googleapis.com">
    - <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    - <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@700&display=swap" rel="stylesheet">
Sheets: 
    - 'https://cdn.jsdelivr.net/npm/holiday.css'
SkipPublish:
    - node_modules
    - htdocs
    - public_html
    - WWW
    - .git
    - .DS_Store
    - .gitignore
---
# Welcome to PocoCMS

## To build from source:
    $ # Create a directory. It doesn't have to be here.
    $ mkdir ~/pococms
    # Navigate to that directory.
    $ cd ~/pococms
    $ # Clone the repo.
    $ git clone https://github.com/pococms/poco
    $ # The repo is now in ~/pococms/poco, so navigate there.
    $ cd poco
    $ # And compile: 
    $ go build 
    $ ### OR....
    $ # There's only one file, so you can also use go run.
    $ # That runs the go compiler, then executes the program 
    $ # if there are no compilation errors.
    $ go run main.go
    $ # This will generate an example file
    $ ./poco
    # (Then make sure poco is on your path)

## To create a website using PocoCMS
    $ mkdir ~/mysite
    $ cd ~/mysite
    $ # Create the home page
    $ nvim index.md # Replace nvim with your favorite editor
    $ poco

## Other command-line examples

Include the two css files shown:

    go run main.go -styles "theme.css light-mode.cs"

Use the docs subdirectory as the root of the site:

    ./poco -root "./docs"
`

var docType = `<!DOCTYPE html>
<html lang=`

var poweredBy = `Powered by PocoCMS`

// assemble takes the raw converted HTML and uses it to generate
// a finished HTML document. Returns it as a string.
func assemble(filename string, article string, fm map[string]interface{}, language string, stylesheetList string) string {
	// This will contain the completed document as a string.
	htmlFile := ""
	// Execute templates. That way {{ .Title }} will be converted into
	// whatever frontMatter["Title"] is set to, etc.
	if parsedArticle, err := doTemplate("", article, fm); err != nil {
		quit(1, err, "%v: Unable to execute ", filename)
	} else {
		article = parsedArticle
	}

	// Build the completed HTML document from the component pieces.
	htmlFile = docType + "\"" + language + "\">" + "\n" +
		"<head>\n" +
		"\t<meta charset=\"utf-8\">\n" +
		"\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
		titletag(fm) +
		metatags(fm) +
		linktags(fm) +
		stylesheets(stylesheetList, fm) +
		"</head>\n<body>\n" +
		layoutEl(fm, "Header", filename) +
		layoutEl(fm, "Nav", filename) +
		"<article>" + article + "</article>\n" +
		layoutEl(fm, "Aside", filename) +
		layoutEl(fm, "Footer", filename) +
		"</body>\n</html>"
	return htmlFile
} //   assemble

// HTML UTILITIES

// layoutEl() takes a layout element file named in the front matter.
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
func layoutEl(fm map[string]interface{}, element string, sourcefile string) string {
	filename := frontMatterStr(element, fm)
	if filename == "" {
		return ""
	}
	isMarkdown := false
	fullPath := ""
	tag := ""
	layoutElSource := frontMatterStr(element, fm)
	// xxx layoutEl()
	fileDir := filepath.Dir(layoutElSource)
	if filepath.IsAbs(layoutElSource) {
		fullPath = layoutElSource
	} else {
		fullPath = filepath.Join(fileDir, layoutElSource)
	}
	if filepath.Ext(fullPath) != ".html" {
		isMarkdown = true
	}

	parsedArticle := ""
	tag = strings.ToLower(element)
	raw := ""
	var err error
	if isMarkdown {
		if raw, _, err = mdYAMLFileToHTMLString(fullPath); err != nil {
			quit(1, err, "Error converting Markdown file %v to HTML", fullPath)
			return ""
		}
		// xxx layoutEl
		if parsedArticle, err = doTemplate("", raw, fm); err != nil {
			quit(1, err, "%v: Unable to execute ", filename)
		}
		wholeTag := "<" + tag + ">" + parsedArticle + "<" + tag + "/>\n"
		return wholeTag
	}
	return fileToString(fullPath) + "\n"

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

// stylesheets() takes stylesheets listed on the command line
// e.g. --styles "foo.css bar.css", and adds them to
// the head. It then generates stylesheet tags for the ones listed in
// the front matter.
// Those listed in the front matter are appended, so they take
// priority.
func stylesheets(sheets string, fm map[string]interface{}) string {
	var globalSlice []string
	var globals string
	if sheets != "" {
		// Build a string from stylesheets named on the command line.
		globalSlice = strings.Split(sheets, " ")
		globals = sliceToStylesheetStr(globalSlice)
	}
	// Build a string from stylesheets named in the front matter for this page
	localSlice := frontMatterStrSlice("Sheets", fm)
	locals := sliceToStylesheetStr(localSlice)

	// Stylesheets named in the front matter takes priority,
	// so they goes last. This allows you to have stylesheets
	// on the command line that act as templates, but that
	// you can override using stylesheets named in
	// the front matter.
	return globals + locals
}

// The --verbose flag. It shows progress as the site is created.
// Required by the Verbose() function.
var gVerbose bool

func main() {

	// cleanup determines whether or not the publish (aka WWW) directory
	// gets deleted on start.
	var cleanup bool
	flag.BoolVar(&cleanup, "cleanup", true, "Delete publish directory before converting files")

	// skip lets you skip the named files from being processed
	var skip string
	flag.StringVar(&skip, "skip", "node_modules .git .DS_Store .gitignore", "List of files to skip when generating a site")

	// language sets HTML lang= value, such as <html lang="fr">
	var language string
	flag.StringVar(&language, "language", "en", "HTML language designation, such as en or fr")

	var root string
	flag.StringVar(&root, "root", ".", "Starting directory of the project")

	// List of stylesheets to include on each page.
	var stylesheets string
	flag.StringVar(&stylesheets, "styles", "", "One or more stylesheets (use quotes if more than one)")

	// Title tag.
	var title string
	flag.StringVar(&title, "Title", poweredBy, "Contents of the HTML title tag")

	// Verbose shows progress as site is generated.
	flag.BoolVar(&gVerbose, "verbose", false, "Display information about project as it's generated")

	// webroot is the directory used to house the final generated website.
	var webroot string
	flag.StringVar(&webroot, "webroot", "WWW", "Subdirectory used for generated HTML files")

	// Process command line flags such as --verbose, --title and so on.
	flag.Parse()

	// See if asource file was specified. Otherwise the whole directory
	// and nested subdirectories are processed.
	filename := flag.Arg(0)

	if filename != "" {
		// Something's left on the command line. It's presumed to
		// be a filename.
		if !fileExists(filename) {
			quit(1, nil, "Can't find the file %v", filename)
		} else {
			// Special case: if you name a file on the command line, it will
			// generate an HTML document from that file and pass you the new filename.
			// The output file isn't published to webroot. It's published to the
			// current directory.
			htmlFilename := buildFileToFile(filename, stylesheets, language)
			quit(0, nil, "Built file %s", htmlFilename)
		}
	}
	// No file was given on the command line.
	// Build the project in place.

	// markdownExtensions are how PocoCMS figures out whether
	// a file is Markdown. If it ends in any one of these then
	// it gets converted to HTML.
	var markdownExtensions searchInfo
	markdownExtensions.list = []string{".md", ".mkd", ".mdwn", ".mdown", ".mdtxt", ".mdtext", ".markdown"}

	webrootPath := buildSite(root, webroot, skip, markdownExtensions, language, stylesheets, cleanup)
	quit(0, nil, "Site published to %s", webrootPath)

}

// TEMPLATE FUNCTIONS

// doTemplate takes HTML in source, expects parsed front
// matter in fm, and executes Go templates
// against the source.
// Returns a string containing the HTML with the
// template values embedded.
func doTemplate(templateName string, source string, fm map[string]interface{}) (string, error) {
	if templateName == "" {
		templateName = "PocoCMS"
	}
	// fmt.Printf("\tdoTemplate() fm: \n%v\n", fm)
	tmpl, err := template.New(templateName).Parse(source)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, fm)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}

// buildFileToFile converts a file from Markdown to HTML, generates an output file,
// and returns name of destination file
func buildFileToFile(filename string, stylesheets string, language string) (outfile string) {
	// Convert Markdown file filename to raw HTML, then assemble a complete HTML document to be published.
	// Return the document as a string.
	html, htmlFilename := buildFileToString(filename, stylesheets, language)
	// Write the contents of the completed HTML document to a file.
	writeStringToFile(htmlFilename, html)
	// Return the name of the converted file
	return htmlFilename
}

// buildFileToString converts a file from Markdown to raw HTML,
// pulls in everything required to create a complete HTML document,
// generates an output file,
// and returns name of the destination HTML file
// Does not check if the input file is Markdown.
// TODO: Ideally this would be called from buildSite()
// Reeturns the string and the filenlame
func buildFileToString(filename string, stylesheets string, language string) (string, string) {
	// Exit silently if not a valid file
	if filename == "" || !fileExists(filename) {
		return "", ""
	}
	// This will be the proposed name for the completed HTML file.
	dest := ""
	// Convert the Markdown file to an HTML string
	if rawHTML, fm, err := mdYAMLFileToHTMLString(filename); err != nil {
		quit(1, err, "Error converting Markdown file %v to HTML", filename)
		return "", ""
	} else {
		// Strip original file's Markdown extension and make
		// the destination files' extension HTML
		dest = replaceExtension(filename, "html")
		// Take the raw converted HTML and use it to generate a complete HTML document in a string
		finishedDocument := assemble(filename, rawHTML, fm, language, stylesheets)
		// Return the finished document and its filename
		return finishedDocument, dest
	}
}

// buildSite takes projectDir as the root directory,
// converts all files (except those in skipPublish.List) to HTML,
// and deposits them in webroot. Attempts to create webroot if it
// doesn't exist. webroot is expected to be a subdirectory of
// projectDir.
// Return name of the root directory files are published to
func buildSite(projectDir string, webroot string, skip string, markdownExtensions searchInfo, language string, stylesheets string, cleanup bool) string {

	var err error
	// Make sure it's a valid site.
	if !isProject(projectDir) {
		quit(1, err, "%s doesn't seem to be a valid site. There's no index.md or README.md", projectDir)
	}

	// Change to requested directory
	if err = os.Chdir(projectDir); err != nil {
		quit(1, err, "Unable to change to directory %s", projectDir)
	}

	// Cache project's root directory
	var homeDir string
	if homeDir, err = os.Getwd(); err != nil {
		quit(1, err, "Unable to get name of current directory")
	}

	// Delete web root directory unless otherwise requested
	if cleanup {
		delDir := filepath.Join(homeDir, webroot)
		Verbose("Deleting directory %v", delDir)
		if err := os.RemoveAll(delDir); err != nil {
			quit(1, err, "Unable to delete publish directory %v", delDir)
		}
	}

	// Convert the list of exclusions into a string slice.
	// xxx
	// skipPublish = getSkipPublish()
	var skipPublish searchInfo
	skipPublish.list = strings.Split(skip, " ")

	// Collect all the files required for this project.
	files, err := getProjectTree(".", skipPublish)
	if err != nil {
		quit(1, err, "Unable to get directory tree")
	}

	// Now have list of all files in directory tree.
	// If markdown, convert to HTML and copy that file to the HTML publication directory.
	// If not, copy to target publication directory unchanged.

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

	// Parsed front matter
	var fm map[string]interface{}

	// Main loop. Traverse the list of files to be copied.
	// If a file is Markdown as determined by its file extension,
	// convert to HTML and copy to output directory.
	// If a file isn't Markdown, copy to output directory with
	// no processing.
	for _, filename := range files {

		// true if it's  Markdown file converted to HTML
		converted = false

		// Get the fully qualified pathname for this file.
		filename = filepath.Join(homeDir, filename)

		// Separate out the file's origin directory
		sourceDir := filepath.Dir(filename)

		Verbose("%s", filename)

		// Get the relatve directory. For example, if your directory
		// is ~/raj/blog and you're in ~/raj/blog/2023/may, then
		// the relative directory is 2023/may.
		if rel, err = filepath.Rel(homeDir, sourceDir); err != nil {
			quit(1, err, "Unable to get relative paths of %s and %s", filename, webroot)
		}

		// Determine the destination directory.
		webrootPath = filepath.Join(homeDir, webroot, rel)
		// Obtain file extension.
		ext := path.Ext(filename)
		// Replace converted filename extension, from markdown to HTML.
		// Only convert to HTML if it has a Markdown extension.
		if markdownExtensions.Found(ext) {
			// Convert the Markdown file to an HTML string
			if HTML, fm, err = mdYAMLFileToHTMLString(filename); err != nil {
				quit(1, err, "Error converting Markdown file to HTML")
			}
			// Strip original file's Markdown extension and make
			// the destination files' extension HTML
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
				quit(1, err, "Unable to create directory %s", webrootPath)
			}
		}
		if converted {
			// Take the raw converted HTML and use it to generate a complete HTML document in a string
			h := assemble(filename, HTML, fm, language, stylesheets)
			writeStringToFile(target, h)
		} else {
			copyFile(source, target)
		}
	}
	// This is where the files were published
	ensureIndexHTML(webrootPath)
	return webrootPath
}

// ensureIndexHTML makes sure there's an index.html file
// in the directory. It's required because some existing
// projects use README.md instead of index.md.
// Even if a project directory had both
// an index.md and a README.md, the output README.html
// would be renamed to index.html
func ensureIndexHTML(path string) {
	readmeHTML := filepath.Join(path, "README.html")
	newIndexHTML := filepath.Join(path, "index.html")

	// if neither index.html nor README.html, then they
	// were missing source files to begin with.
	if !fileExists(readmeHTML) && !fileExists(newIndexHTML) {
		quit(1, nil, "No README.html or index.html could be found in %v", path)
	}

	// Both README.html and index.html exist.  Or
	// README.html exists but no index.html exists.
	// Rename README.html
	if fileExists(readmeHTML) && (fileExists(newIndexHTML) || !fileExists(newIndexHTML)) {
		err := os.Rename(readmeHTML, newIndexHTML)
		if err != nil {
			quit(1, err, "Unable to rename %v ", readmeHTML)
		}
		return
	}
	/*
		// Only README.html remains. Rename it.
		err := os.Rename(readmeHTML, newIndexHTML)
		if err != nil {
			quit(1, err, "Unable to create %v ", newIndexHTML)
		}
	*/

}

func getSkipPublish() []string {
	// var skipPublish searchInfo
	// skipPublish.list = strings.Split(skip, " ")
	// xxx
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
// Clear but
func copyFile(source string, target string) {
	if source == target {
		quit(1, nil, "copyFile: %s and %s are the same", source, target)
	}
	if source == "" {
		quit(1, nil, "copyFile: no source file specified")
	}
	if target == "" {
		quit(1, nil, "copyFile: no destination file specified for file %s", source)
	}
	var src, trgt *os.File
	var err error
	if src, err = os.Open(source); err != nil {
		quit(1, err, "copyFile: Unable to open file %s", source)
	}
	defer src.Close()

	if trgt, err = os.Create(target); err != nil {
		quit(1, err, "copyFile: Unable to create file %s", target)
	}
	if _, err := trgt.ReadFrom(src); err != nil {
		quit(1, err, "Error copying file %s to %s", source, target)
	}
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
// Fails quietly  if unable to open the file, since
// we're just generating HTML.
func fileToString(filename string) string {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
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
func writeStringToFile(filename, contents string) {
	var out *os.File
	var err error
	if out, err = os.Create(filename); err != nil {
		quit(1, err, "writeStringToFile: Unable to create file %s", filename)
	}
	if _, err = out.WriteString(contents); err != nil {
		// TODO: Renumber error code?
		quit(1, err, "Error writing to file %s", filename)
	}
}

// SLICE UTILITIES
// Searching a sorted slice is fast.
// This tracks whether the slice has been sorted
// and sorts it on first search.

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

// mdYAMLFiletoHTML converts a Markdown document
// with YAML front matter to HTML.
// Returns a byte slice containing the HTML source.
func mdYAMLFileToHTMLString(filename string) (string, map[string]interface{}, error) {
	source := fileToBuf(filename)
	if HTML, fm, err := mdYAMLToHTML(source); err != nil {
		return "", nil, err
	} else {
		return string(HTML), fm, nil
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
		highlighting.NewHighlighting(
			highlighting.WithStyle("github"),
			highlighting.WithFormatOptions()),
	}

	parserOpts := []parser.Option{
		parser.WithAttribute(),
		parser.WithAutoHeadingID()}

	renderOpts := []renderer.Option{
		// WithUnsafe is required for HTML templates to work properly
		html.WithUnsafe(),
		html.WithXHTML(),
	}
	return goldmark.New(
		goldmark.WithExtensions(exts...),
		goldmark.WithParserOptions(parserOpts...),
		goldmark.WithRendererOptions(renderOpts...),
	)
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
func frontMatterStr(key string, fm map[string]interface{}) string {
	v := fm[key]
	value, ok := v.(string)
	if !ok {
		return ""
	}
	return value
}

// frontMatterStrSlice obtains a list of string values from the front matter.
// For example, if you had this code in your Markdown file:
// ---
// Sheets:
//   - 'https://cdn.jsdelivr.net/npm/holiday.css'
//   - 'fonts.css'
//
// ---
//
// I like yo mama
func frontMatterStrSlice(key string, fm map[string]interface{}) []string {
	if key == "" {
		return []string{}
	}
	v, ok := fm[key].([]interface{})
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

// linkTags() obtains the list of link tags from the "LinkTags" front matter
// and inserts them into the document.
func linktags(fm map[string]interface{}) string {
	linkTags := frontMatterStrSlice("LinkTags", fm)
	if len(linkTags) < 1 {
		return ""
	}
	tags := ""
	for _, tag := range linkTags {
		tags += "\t" + tag + "\n"
	}
	return tags
}

func titletag(fm map[string]interface{}) string {
	title := frontMatterStr("Title", fm)
	if title == "" {
		return "\t<title>" + poweredBy + "</title>\n"
	} else {
		return "\t<title>" + title + "</title>\n"
	}
}

// Generate common metatags
func metatags(fm map[string]interface{}) string {
	return metatag("description", frontMatterStr("Description", fm)) +
		metatag("keywords", frontMatterStr("Keywords", fm)) +
		metatag("author", frontMatterStr("Author", fm))
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

// Printy utilities

// If the Verbose flag is set, use the Printf style parameters
// to format the input and return a string.
func Verbose(format string, ss ...interface{}) {
	if gVerbose {
		fmt.Println(fmtMsg(format, ss...))
	}
}

// / xxx
func quit(exitCode int, err error, format string, ss ...interface{}) {
	msg := fmt.Sprint(fmtMsg(format, ss...))
	errmsg := ""
	if err != nil {
		errmsg = " " + err.Error()
	}
	// fmt.Println(msg + errmsg)
	fmt.Printf("PocoCMS: %s%s\n", msg, errmsg)
	os.Exit(exitCode)
}

// debug displays messages to stdout using Fprintf syntax.
// A little list printing and easier to search
func debug(format string, ss ...interface{}) {
	fmt.Println(fmtMsg(format, ss...))
}

// fmtMsg() takes a list of strings like Fprintf, interpolates, and writes to a string
func fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}

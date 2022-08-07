// PocoCMS: Markdown-based CMS in 1 file, written in Go
package main

// git clone https://github.com/pococms/poco
// cd poco
// go mod init github.com/pococms/poco
// go mod tidy
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
// Use the docs subdirectory as the root of the site.
// ./pococms -root "./docs"

// Get CSS file from CDN
// poco -styles "https://unpkg.com/spectre.css/dist/spectre.min.css"
// poco -styles "//writ.cmcenroe.me/1.0.4/writ.min.css" foo.md

// Notes:
// - www is a subdir of project
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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var indexSample = `---
Description: 'PocoCMS: Markdown-based CMS in 1 file, written in Go'
Title: 'Powered by PocoCMS'
Author: 'Tom Campbell'
Header: header.html
Nav: nav.html
Footer: footer.html
Sheets: 
 - 'https://cdn.jsdelivr.net/npm/holiday.css'
---
# Welcome to PocoCMS

## To build from source:
    $ git clone https://github.com/pococms/poco
    $ cd poco
    $ go mod init github.com/pococms/poco
    $ go mod tidy
    $ go build # OR EVEN go run main.go
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

// Assemble takes the raw converted HTML and uses it to generated
// a finished HTML document.
func assemble(article string, frontMatter map[string]string, language string, stylesheetList string) string {
	var htmlFile string
	var sheets string
	styles := strings.Split(stylesheetList, " ")
	if stylesheetList != "" {
		for _, sheet := range styles {
			s := fmt.Sprintf("\t<link rel=\"stylesheet\" href=\"%s\"/>\n", sheet)
			sheets += s
		}
	}
	htmlFile = docType + "\"" + language + "\">" + "\n" +
		"<head>\n" +
		"\t<meta charset=\"utf-8\">\n" +
		"\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n" +
		"\t<title>" + frontMatter["Title"] + "</title>\n" +
		metatags(frontMatter) +
		stylesheets(sheets, frontMatter) +
		"</head>\n<body>\n" +
    layoutEl(frontMatter, "Header") +
    layoutEl(frontMatter, "Nav") +
		article +
    layoutEl(frontMatter, "Footer") +
		"</body>\n</html>"
	return htmlFile
}

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
//    <header>hello, world.</header>
// This function would read in the head.html file (or whatever
// the file was named in the front matter) and insert it before the
// body of the document.
func layoutEl(frontMatter map[string]string, element string) string {
	filename := frontMatter[element]
	if filename == ""{
		return ""
	}
  if !fileExists(filename) {
    return ""
  }
	return fileToString(filename) + "\n"
}


// stylesheets() takes stylesheets listed on the command line
// e.g. --styles "foo.css bar.css" and adds them to
// the head.
// It does generates stylesheet tags for the ones listed in
// the front matter. The latter are appended, so they take
// priority.
func stylesheets(sheets string, frontMatter map[string]string) string {
	s := strings.Split(frontMatter["Sheets"], " ")
	var frontStyles string
	for _, sheet := range s {
		// Why? This seems to be a Go thing.
		// More likely I'm missing something.
		sheet = strings.ReplaceAll(sheet, "[", "")
		sheet = strings.ReplaceAll(sheet, "]", "")
		tag := fmt.Sprintf("\t<link rel=\"stylesheet\" href=\"%s\"/>\n", sheet)
		if sheet != "" {
			frontStyles += tag
		}
	}
	// Stylesheets named in the front matter takes priority,
	// so they goes last. This allows you to have stylesheets
	// on the command line that act as templates, but that
	// you can override using stylesheets named in
	// the front matter.
	return sheets + frontStyles
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
	flag.StringVar(&skip, "skip", "List of files to skip when generating a site", "node_modules .git .DS_Store .gitignore")

	// language sets HTML lang= value, such as <html lang="fr">
	var language string
	flag.StringVar(&language, "language", "en", "HTML language designation, such as en or fr")

	// root is the project's root directory, where the home page is located.
	// Defaults to the current directory but other choices might be
	// "/docs" or "/_pub"
	var root string
	flag.StringVar(&root, "root", ".", "Subdirectory to use as root")

	// List of stylesheets to include on each page.
	var stylesheets string
	flag.StringVar(&stylesheets, "styles", "", "One or more stylesheets (use quotes if more than one)")

	// Title tag.
	var title string
	flag.StringVar(&title, "title", "powered by PocoCMS", "Contents of the HTML title tag")

	// Verbose shows progress as site is generated.
	flag.BoolVar(&gVerbose, "verbose", false, "Display information about project as it's generated")

	// www is the directory used to house the final generated website.
	var www string
	flag.StringVar(&www, "www", "WWW", "Subdirectory used for generated HTML files")

	// Process command line flags such as --verbose, --title and so on.
	flag.Parse()

	// See if asource file was specified. Otherwise the whole directory
	// and nested subdirectories are processed.
	filename := flag.Arg(0)
	if filename != "" && !fileExists(filename) {
		quit(fmt.Sprintf("Can't find a file named %s", filename), nil, 1)
	}

	// Convert the list of stylesheets into a string slice.
	//styles := strings.Split(stylesheets, " ")

	// markdownExtensions are how PocoCMS figures out whether
	// a file is Markdown. If it ends in any one of these then
	// it gets converted to HTML.
	var markdownExtensions searchInfo
	markdownExtensions.list = []string{".md", ".mkd", ".mdwn", ".mdown", ".mdtxt", ".mdtext", ".markdown"}

	// See if there's an index.md in the starting directory
	var rootFile string
	var err error
	if rootFile, err = filepath.Abs(root); err != nil {
		quit("Error detecting root file.", err, 1)
	}

	// See if there's an index.md at the root of the
	// project. If not, create one.
	indexMd := filepath.Join(rootFile, "index.md")
	if !fileExists(indexMd) {
		writeStringToFile(indexMd, indexSample)
	}
	targetDir := mdDirectoryTreeToHTML(root, www, skip, markdownExtensions, language, stylesheets, cleanup)
	quit(fmt.Sprintf("Files published to %s", targetDir), nil, 0)
	// TODO:
	// Preserve for the single file demo
	/*
		if HTML, err := mdFileToHTML(filename); err != nil {
			quit("Error creating Markdown file", err, 1)
		} else {
			assemble(HTML, title, language, stylesheets)
			//fmt.Println(HTML)
			quit("Complete", nil, 0)
		}
	*/

}

// mdToHTML takes Markdown source as a byte slice and converts it to HTML
// using Goldmark's default settings.
func mdToHTML(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(input, &buf); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// mdFileToHTML converts a source file to an HTML string
// using Goldmark's default settings.
func mdFileToHTML(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if HTML, err := mdToHTML(bytes); err != nil {
		return "", err
	} else {
		return string(HTML), nil
	}
}

func quit(msg string, err error, exitCode int) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err.Error())
	} else {
		fmt.Printf("%s\n", msg)
	}
	os.Exit(exitCode)
}

// mdDirectoryTreeToHTML takes startDir as the root directory,
// converts all files (except those in exclude.List) to HTML,
// and deposits them in www. Attempts to create www if it
// doesn't exist. www is expected to be a subdirectory of
// startDir.
// Return name of the root directory files are published to
//func mdDirectoryTreeToHTML(startDir string, www string, skip string, markdownExtensions searchInfo, language string, styles []string, cleanup bool) string {
func mdDirectoryTreeToHTML(startDir string, www string, skip string, markdownExtensions searchInfo, language string, stylesheets string, cleanup bool) string {

	var err error

	// Change to requested directory
	if err = os.Chdir(startDir); err != nil {
		quit(fmt.Sprintf("Unable to change to directory %s", startDir), err, 1)
	}

	// Cache project's root directory
	var currDir string
	if currDir, err = os.Getwd(); err != nil {
		quit("Unable to get name of current directory", err, 1)
	}

	if cleanup {
		delDir := filepath.Join(currDir, www)
		Verbose("Deleting directory %v", delDir)
		if err := os.RemoveAll(delDir); err != nil {
			quit(fmt.Sprintf("Unable to delete publish directory %v", delDir), err, 1)
		}
	}

	// Convert the list of exclusions into a string slice.
	var exclude searchInfo
	exclude.list = strings.Split(skip, " ")

	// Collect all the files required for this project.
	files, err := getProjectTree(".", exclude)
	if err != nil {
		quit("Unable to get directory tree", err, 1)
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
	var frontMatter map[string]string

	// Name of directory used to publish output files
	var targetDir string

	// Main loop. Traverse the list of files to be copied.
	// If a file is Markdown as determined by its file extension,
	// convert to HTML and copy to output directory.
	// If a file isn't Markdown, copy to output directory with
	// no processing.
	for _, filename := range files {

		// true if it's  Markdown file converted to HTML
		converted = false

		// Get the fully qualified pathname for this file.
		filename = filepath.Join(currDir, filename)

		// Separate out the file's origin directory
		sourceDir := filepath.Dir(filename)

		Verbose("%s", filename)

		// Get the relatve directory. For example, if your directory
		// is ~/raj/blog and you're in ~/raj/blog/2023/may, then
		// the relative directory is 2023/may.
		if rel, err = filepath.Rel(currDir, sourceDir); err != nil {
			quit(fmt.Sprintf("Unable to get relative paths of %s and %s\n", filename, www), err, 1)
		}

		// Determine the destination directory. If the base publish
		// directory is named WWW, then in the previous example
		// it would be ~/raj/blog/WWW, or ~/raj/blog/WWW/2023/may
		targetDir = filepath.Join(currDir, www, rel)
		// Obtain file extension.
		ext := path.Ext(filename)
		// Replace converted filename extension, from markdown to HTML.
		// Only convert to HTML if it has a Markdown extension.
		if markdownExtensions.Found(ext) {
			// Convert the Markdown file to an HTML string
			if HTML, frontMatter, err = mdYAMLFileToHTMLString(filename); err != nil {
				quit("Error converting Markdown file to HTML", err, 1)
			}
			// Strip original file's Markdown extension and make
			// the destination files' extension HTML
			source = filename[0:len(filename)-len(ext)] + ".html"
			converted = true
		} else {
			// Not a Markdown file. Copy unchanged.
			source = filename
			// Insert destination (WWW) directory
			converted = false
		}
		target = filepath.Join(targetDir, filepath.Base(source))

		// Create the target directory for this file if it
		// doesn't exist.
		if !dirExists(targetDir) {
			err := os.MkdirAll(targetDir, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				quit(fmt.Sprintf("Unable to create directory %s", targetDir), err, 1)
			}
		}
		if converted {
			h := assemble(HTML, frontMatter, language, stylesheets)
			writeStringToFile(target, h)
		} else {
			copyFile(source, target)
		}
	}
	// This is where the files were published
	return targetDir
}

// FILE UTILITIES
// Clear but
func copyFile(source string, target string) {
	if source == target {
		quit(fmt.Sprintf("copyFile: %s and %s are the same", source, target), nil, 1)
	}
	if source == "" {
		quit("copyFile: no source file specified", nil, 1)
	}
	if target == "" {
		quit(fmt.Sprintf("copyFile: no destination file specified for file %s", source), nil, 1)
	}
	var src, trgt *os.File
	var err error
	if src, err = os.Open(source); err != nil {
		quit(fmt.Sprintf("copyFile: Unable to open file %s", source), err, 1)
	}
	defer src.Close()

	if trgt, err = os.Create(target); err != nil {
		quit(fmt.Sprintf("copyFile: Unable to create file %s", target), err, 1)
	}
	if _, err := trgt.ReadFrom(src); err != nil {
		quit(fmt.Sprintf("Error copying file %s to %s", source, target), err, 1)
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

// writeStringToFile creates a file called filename without checking to see if it
// exists, then writes contents to it.
// filename is afully qualified pathname.
// contents is the string to write
func writeStringToFile(filename, contents string) {
	var out *os.File
	var err error
	if out, err = os.Create(filename); err != nil {
		quit(fmt.Sprintf("writeStringToFile: Unable to create file %s", filename), err, 1)
	}
	if _, err = out.WriteString(contents); err != nil {
		// TODO: Renumber error code?
		quit(fmt.Sprintf("Error writing to file %s", filename), err, 1)
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

func visit(files *[]string, exclude searchInfo) filepath.WalkFunc {
	// var exclude searchInfo
	// Find out what directories to exclude
	//exclude.list = []string{"node_modules", "main.bak", ".git", "pub", ".DS_Store", ".gitignore"}
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
		if exclude.Found(name) && isDir {
			return filepath.SkipDir
		}
		// It may be just a filename on the exclude list.
		if exclude.Found(name) {
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
// Ignore items in excluce.List
func getProjectTree(path string, exclude searchInfo) (tree []string, err error) {
	var files []string
	err = filepath.Walk(path, visit(&files, exclude))
	if err != nil {
		return []string{}, err
	}
	return files, nil
}

// mdYAMLFiletoHTML converts a Markdown document
// with YAML front matter to HTML.
// Returns a byte slice containing the HTML source.
func mdYAMLFileToHTMLString(filename string) (string, map[string]string, error) {
	source := fileToBuf(filename)
	//frontMatter := make(map[string]string)
	if HTML, frontMatter, err := mdYAMLToHTML(source); err != nil {
		return "", nil, err
	} else {
		return string(HTML), frontMatter, nil
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
func mdYAMLToHTML(source []byte) ([]byte, map[string]string, error) {
	mdParser := newGoldmark()
	mdParserCtx := parser.NewContext()
	// YAML front matter
	var frontMatter map[string]interface{}

	var buf bytes.Buffer
	// Convert Markdown source to HTML and deposit in buf.Bytes().
	if err := mdParser.Convert(source, &buf, parser.WithContext(mdParserCtx)); err != nil {
		return []byte{}, nil, err
	}
	// Obtain YAML front matter from document.
	frontMatter = meta.Get(mdParserCtx)
	frontMatterMap := make(map[string]string)
	for key, value := range frontMatter {
		k := fmt.Sprintf("%v", key)
		v := fmt.Sprintf("%v", value)
		frontMatterMap[k] = v
	}
	return buf.Bytes(), frontMatterMap, nil
}

// Generate HTML

// Generate common metatags
func metatags(frontMatter map[string]string) string {
	return metatag("description", frontMatter["Description"]) +
		metatag("author", frontMatter["Author"])
}

// metatag() generates a metatag such as <meta name="description"content="PocoCMS: Markdown-based CMS in 1 file, written in Go">
func metatag(tag string, content string) string {
	if content == "" {
		return ""
	}
	return "\t<meta name=\"" + tag + "\"" +
		" content=" + "\"" + content + "\">\n"
}

// Printy utilities
func Verbose(format string, ss ...interface{}) {
	if gVerbose {
		fmt.Println(fmtMsg(format, ss...))
	}
}

// fmtMsg() formats string like Fprintf and writes to a string
func fmtMsg(format string, ss ...interface{}) string {
	return fmt.Sprintf(format, ss...)
}

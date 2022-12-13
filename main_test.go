package main

// TODO: Come up with something better than "$$_" for filenames.
// https://pkg.go.dev/os#TempDir
import (
	//"regexp"
	//"fmt"
	"golang.org/x/exp/slices"
	"os"
	"strings"
	"testing"
)

// ********************************************************
// FRONT MATTER KEYS THAT RETURN STRINGS
// ********************************************************

// All front matter keys that return strings
var fmAllStr = `---
title: "PocoCMS title"
author: "Tom Campbell"
pagetheme: "tufte"
theme: "pocodocs"
branding: "PocoCMS for the win!"
description: "Build informational websites friction-free"

---
`

type strTest struct {
	key      string
	expected string
}

// Plan to get all string values here
var fmStrTests = []strTest{
	{"title", "PocoCMS title"},
	{"author", "Tom Campbell"},
	{"pagetheme", "tufte"},
	{"theme", "pocodocs"},
	{"branding", "PocoCMS for the win!"},
	{"description", "Build informational websites friction-free"},
}

// TestAllFmStrs tests all string values
// one can expect in front matter.
// For string slice values in front matter, see TestAllFmSlices
func TestAllFmStrs(t *testing.T) {
	source := fmAllStr
	var err error
	var fm map[string]interface{}
	if _, fm, err = mdYAMLToHTML([]byte(source)); err != nil {
		t.Errorf("Unable to get front matter from %s", source)
	}
	for _, tt := range fmStrTests {
		actual := fmStr(tt.key, fm)
		if actual != tt.expected {
			t.Errorf("FrontMatter:%s: expected \"%s\", actual \"%s\"", tt.key, actual, tt.expected)
		}
	}
}

// ********************************************************
// FRONT MATTER KEYS THAT RETURN STRING SLICES
// ********************************************************

// All front matter keys that return string slices
var fmStrSliceTests = []struct {
	code     string
	key      string
	expected []string
}{

	// ** stylesheets front matter
	// ---
	{
		`----
stylesheets:
- foo.css
---
`,
		"stylesheets",
		[]string{"foo.css"},
	},
	// ---

	{
		`----
stylesheets:
- foo.css
- bar.css
---
`,
		"stylesheets",
		[]string{"foo.css", "bar.css"},
	},

	// ** style-tags front matter
	// ---
	{
		`----
style-tags:
---
`,
		"style-tags",
		[]string{},
	},
	// ---

	// ---
	{
		`----
style-tags:
- "article>p{color:blue}"
---
`,
		"style-tags",
		[]string{"article>p{color:blue}"},
	},
	// ---

}

// TestAllFmSlices tests all string slice values
// one can expect in front matter.
// For string values in front matter, see TestAllFmStrs
func TestAllFmSlices(t *testing.T) {
	var err error
	var fm map[string]interface{}
	for _, tt := range fmStrSliceTests {
		if _, fm, err = mdYAMLToHTML([]byte(tt.code)); err != nil {
			t.Errorf("Unable to get front matter from %s", tt.code)
		}
		slice := fmStrSlice(tt.key, fm)
		if !slices.Equal(slice, tt.expected) {
			t.Errorf("%v: expected \"%s\", actual \"%s\"", tt.key, slice, tt.expected)
		}
	}
}

// ********************************************************
// RAW HTML OUTPUT WITH DEFAULT SETTINGS
// ********************************************************

// Tests only output of article, not other page
// layout elements such as header, footer, etc.
var articleMdToHTMLTests = []struct {
	code     string
	expected string
}{

	// TEST RECORD
	{
		// Markdown portion
		`hello`,
		// Expected output portion
		`<p>hello</p>`,
	},

	// TEST RECORD
	{
		// Markdown portion
		`# hello`,
		// Expected output portion
		`<h1 id="hello">hello</h1>`,
	},
}

// testArticleCode takes markup and generates the raw HTML for
// that markup. It tests only output of article, not other page
// layout elements such as header, footer, etc.
func TestArticleCode(t *testing.T) {
	for _, tt := range articleMdToHTMLTests {
		actual := mdYAMLStringToTemplatedHTMLString(newConfig(), "fakefilename", tt.code)
		actual = strings.TrimSpace(string(actual))
		if actual != tt.expected {
			t.Errorf("Expected %s. Got %s", tt.expected, actual)
		}
	}
}

// ********************************************************
// TITLE TAG
// ********************************************************

var fmTitleMissing = `---
Author: "yo mama"
---
`

func TestMissingTitleTag(t *testing.T) {
	source := fmTitleMissing
	c := newConfig()
	//var fm map[string]interface{}
	var err error
	if _, c.fm, err = mdYAMLToHTML([]byte(source)); err != nil {
		t.Errorf("Failed converting %s to HTML and obtaining front matter", source)
	}

	actual := c.titleTag()
	actual = strings.TrimSpace(actual)
	expected := strings.TrimSpace("\t<title>" + poweredBy + "</title>")
	if actual != expected {
		t.Errorf("titleTag(): expected %s. Got %s", expected, actual)
	}
}

// ********************************************************
// SEARCHINFO UTILITIES
// ********************************************************

func TestSearchInfo(t *testing.T) {
	var s searchInfo
	s.AddStr("a")
	expected := []string{"a"}
	if !slices.Equal(s.list, expected) {
		t.Errorf("%v should be empty, but it's %v",
			s.list, expected)
	}
	// Add out of alphabetical order.
	// Result should still be sorted correctly.
	s.AddStr("c")
	s.AddStr("b")
	expected = []string{"a", "b", "c"}
	if !slices.Equal(s.list, expected) {
		t.Errorf("s.list is %v, but it should be %v",
			s.list, expected)
	}
	// "d" should not be in the list
	searchFor := "d"
	found := s.Found(searchFor)
	if found {
		t.Errorf("%v reported as found, but it should not be in s.list",
			searchFor)
	}
	// "c" should be in the list
	searchFor = "c"
	found = s.Found(searchFor)
	if !found {
		t.Errorf("%v reported as not found, but it should be present in s.list",
			searchFor)
	}
}

// ********************************************************
// GETFM RETURNS FRONT MATTER FOR FILE PASSED IN
// ********************************************************

var getFmTests = []struct {
	filename        string
	code            string
	fmKey           string
	fmExpectedValue string
}{

	// TEST RECORD
	{
		// filename
		// TODO: Use temp file or defer deletion of this file
		"$$_README.md",
		// Contents of Markdown file
		`---
title: "yo mama"
---
`,
		// Front matter key to read
		"title",
		// Value the front matter should return
		"yo mama",
	},
}

// getFm() should return front matter that has
// nothing to do with the front matter passed
// in with the c, the config object.
func TestGetFm(t *testing.T) {
	for _, tt := range getFmTests {
		c := newConfig()
		// Create a Markdown file on the fly (with optional
		// front matter). Obtain its front matter.
		fm := c.getFm(stringToFile(c, tt.filename, tt.code))
		value := fmStr(tt.fmKey, fm)
		if value != tt.fmExpectedValue {
			t.Errorf("Frontmatter error. fmStr(%v) should was %v. Expected %v",
				tt.fmKey, value, tt.fmExpectedValue)
		}
	}
}

// ********************************************************
// SLICETOSTYLESHEETSTR
// ********************************************************

var sliceToStylesheetStrTests = []struct {
	slice    []string
	expected string
}{

	// TEST RECORD
	{
		// slice of stylesheet names
		[]string{"foo.css"},
		// slice should be converted to this HTML:
		`<link rel="stylesheet" href="foo.css">`,
	},

	// TEST RECORD
	{
		// slice of stylesheet names
		[]string{"foo.css", "bar.css"},
		// slice should be converted to this HTML:
		"<link rel=\"stylesheet\" href=\"foo.css\">\n\t<link rel=\"stylesheet\" href=\"bar.css\">",
	},
}

func TestSliceToStylesheetStr(t *testing.T) {
	for _, tt := range sliceToStylesheetStrTests {
		actual := strings.TrimSpace(sliceToStylesheetStr("", tt.slice))
		if actual != tt.expected {
			t.Errorf("%v should convert to:\n[%v]\nIt was actually:\n[%v]",
				tt.slice, tt.expected, actual)
		}
	}
}

// ********************************************************
// convertMdYAMLFileToHTMLFragmentStr
// ********************************************************

var convertMdYAMLFileToHTMLFragmentStrTests = []struct {
	code     string
	expected string
}{

	// TEST RECORD
	{
		// Contents of Markdown file
		// note: empty file
		`
`,
		// Expected output when Markdown file is converted to HTML
		``,
	},

	// TEST RECORD
	{
		// Contents of Markdown file
		// note: unformatted text. no YAML front matter.
		`
hello, world.
`,
		// Expected output when Markdown file is converted to HTML
		`<p>hello, world.</p>`,
	},

	// TEST RECORD
	{
		// Contents of Markdown file
		// note: unformatted text. Empty YAML front matter.
		`---
---
hello, world.
`,
		// Expected output when Markdown file is converted to HTML
		`<p>hello, world.</p>`,
	},

	// TEST RECORD
	{
		// Contents of Markdown file
		// note: YAML front matter contains description key.
		`---
description: "PocoCMS"
---
hello, {{ .description }}!
`,
		// Expected output when Markdown file is converted to HTML
		`<p>hello, PocoCMS!</p>`,
	},

	// TEST RECORD
	{
		// Contents of Markdown file
		// note: YAML front matter contains empty description value.
		`---
description: 
---
hello, {{ .description }}!
`,
		// Expected output when Markdown file is converted to HTML
		// TODO: Create an error message explaining this output
		`<p>hello, <no value>!</p>`,
	},

	// -- TEST RECORDS END HERE --
}

func TestMdYAMLStringToTemplatedHTMLString(t *testing.T) {
	for _, tt := range convertMdYAMLFileToHTMLFragmentStrTests {
		c := newConfig()
		actual := mdYAMLStringToTemplatedHTMLString(c, "dummy.md", tt.code)
		actual = strings.TrimSpace(actual)
		if actual != tt.expected {
			t.Errorf("Markdown source is\n%v\nIt converted to:\n%v\nExpected:\n%v",
				tt.code, actual, tt.expected)
		}
	}
}

// xxx

// ********************************************************
//
// ********************************************************

var getTmpTests = []struct {
	code            string
	fmKey           string
	fmExpectedValue string
}{

	// TEST RECORD
	{
		// Contents of Markdown file
		`---
title: "yo mama"
---
`,
		// Front matter key to read
		"title",
		// Value the front matter should return
		"yo mama",
	},
}

func TestGetTmpTests(t *testing.T) {
	for _, tt := range getTmpTests {
		//c := newConfig()
		// Create a Markdown file on the fly (with optional
		// front matter). Obtain its front matter.
		var f *os.File
		var err error
		// xxx
		// Create a temporary file with an .md extension in Go's
		// default temp file directory.
		if f, err = os.CreateTemp("", "*.md"); err != nil {
			t.Errorf("Unable to create temporary file")
		}
		// Delete this file when the function exits
		defer os.Remove(f.Name())
		// Create a file using the temp name and this test record's
		// source code.
		err = os.WriteFile(f.Name(), []byte(tt.code), 0666)

		/*
			fm := c.getFm(stringToFile(c, tt.filename, tt.code))
			value := fmStr(tt.fmKey, fm)
			if value != tt.fmExpectedValue {
				t.Errorf("Frontmatter error. fmStr(%v) should was %v. Expected %v",
					tt.fmKey, value, tt.fmExpectedValue)
			}
		*/
	}
}

// ********************************************************
// UTILITIES
// ********************************************************

package main

import (
	//"regexp"
	//"fmt"
	"golang.org/x/exp/slices"
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
theme: "tufte"
global-theme: "pocodocs"
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
	{"theme", "tufte"},
	{"global-theme", "pocodocs"},
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
		//fmt.Printf("SLICE: %v\n", slice)
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
  //c := newConfig()
	for _, tt := range articleMdToHTMLTests {
    actual := mdYAMLStringToTemplatedHTMLString(newConfig(),tt.code)
    // I actually don't understand why a test case like
    // `# hello` ends up with a trailing newline
    // on the actual value.
    actual = strings.TrimSpace(string(actual))
		if actual != tt.expected {
      t.Errorf("Expected %s. Got %s", tt.expected, actual)
      /*
			t.Errorf("%v: expected \"%v\", actual \"%v\"", 
        tt.code, actual, tt.expected)
        */
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
  //slice := []string{}
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
  expected = []string{"a","b","c"}  
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
    t.Errorf("%v reported as not found, but it be present in s.list",
      searchFor)
  }
}
// ********************************************************
// getFm function
// ********************************************************

// All front matter keys that return strings
/*
var fmAllStr = `---
title: "PocoCMS title"
author: "Tom Campbell"
theme: "tufte"
global-theme: "pocodocs"
branding: "PocoCMS for the win!"
description: "Build informational websites friction-free"

---
`
*/

var getFmTests = []struct {
	filename string
	code string
}{

  // TEST RECORD
	{ 
    // filename
		"README.md",
    // Contents of Markdown file
		`---
title: "yo mama"
---
`,
	},

}

// getFm() should return front matter that has
// nothing to do with the front matter passed
// in with the c, the config object.
func TestGetFm(t *testing.T) {
 	for _, tt := range getFmTests {
    c := newConfig()
    fm := c.getFm(c.instaMd(tt.filename, tt.code))
    if fmStr("title", c.fm) == fmStr("title", fm) {
      t.Errorf("c.fm and fm should be different")
    }
  }
}

// instaMd creates Markdown file on the file
// using the given filename, and 
// the code passed in as a string.
func (c *config)instaMd(filename, code string) string {
  stringToFile(c, filename, code)
  c.fileToString(filename)
  return filename
}





package main

import (
	//"regexp"
	"fmt"
	"golang.org/x/exp/slices"
  "strings"
	"testing"
)

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

var fmTitleMissing = `---
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

// fmtTestStr
func fmTestStr(key string, fm map[string]interface{}) string {
	s := fmt.Sprintf("%s", fm[key])
	return s
}

// ********************************************************
// FRONT MATTER KEYS THAT RETURN STRINGS 
// ********************************************************

// TestAllFmStrings tests all string slice values
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
		actual := fmTestStr(tt.key, fm)
		if actual != tt.expected {
			t.Errorf("FrontMatter:%s: expected \"%s\", actual \"%s\"", tt.key, actual, tt.expected)
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
	var err error
  var a []byte
	for _, tt := range articleMdToHTMLTests {
		if a, _, err = mdYAMLToHTML([]byte(tt.code)); err != nil {
      t.Errorf("Internal error converting code: %s", tt.code)
		}
    // I actually don't understand why a test case like
    // `# hello` ends up with a trailing newline
    // on the actual value.
    actual := strings.TrimSpace(string(a))
		if actual != tt.expected {
			t.Errorf("%v: expected \"%v\", actual \"%v\"", tt.code, actual, tt.expected)
		}
	}
}



/*
func TestMissingTitleTag(t *testing.T) {
	tagLine := poweredBy
	source := fmTitleMissing
	var fm map[string]interface{}
	var err error
	if _, fm, err = mdYAMLToHTML([]byte(source)); err != nil {
		t.Errorf("Failed converting %s to HTML and obtaining front matter", source)
	}

	//value := fmTest("Title", fm)
	expected := "\t<title>" + tagLine + "</title>\n"
	actual := fmStr("title", fm)
	if actual != expected {
		t.Errorf("FrontMatter[\"title\"]: expected \"%s\", actual \"%s\"", expected, actual)
	}
}
*/

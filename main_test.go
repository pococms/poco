package main

import (
	//"regexp"
	"fmt"
	"golang.org/x/exp/slices"
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

/*
// fmtTestStrSlice
func fmTestStrSlice(key string, fm map[string]interface{}) []string {
	s := fmt.Sprintf("%s", fm[key])
	return s
}
*/

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

package main

import (
	//"regexp"
	"fmt"
	"testing"
)

// All front matter keys that return strings
var fmAllStr = `---
title: PocoCMS title
---
`

var fmTitleMissing = `---
---
`

type stringTest struct {
	key      string
	expected string
}

// Plan to get all string values here
var stringTests = []stringTest{
	{"title", "PocoCMS title"},
}

// fmtTest
func fmTestStr(key string, fm map[string]interface{}) string {
	s := fmt.Sprintf("%s", fm[key])
	return s
}

// TestAllFrontMatterStrings tests all string values
// one can expect in front matter. Since it's only
// strings that means keys like Sheets aren't
// tested here.
func TestAllFrontMatterStrings(t *testing.T) {
	source := fmAllStr
	var err error
	var fm map[string]interface{}
	if _, fm, err = mdYAMLToHTML([]byte(source)); err != nil {
		t.Errorf("Unable to get front matter from %s", source)
	}
	for _, tt := range stringTests {
		actual := fmTestStr(tt.key, fm)
		if actual != tt.expected {
			t.Errorf("FrontMatter:%s: expected \"%s\", actual \"%s\"", tt.key, actual, tt.expected)
		}
	}
}

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

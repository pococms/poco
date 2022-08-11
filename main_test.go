package main

import (
	//"regexp"
	"fmt"
	"testing"
)

// All front matter keys that return strings
var fmAllStr = `---
Title: PocoCMS title
---
`

var fmTitleMissing = `---
---
`

type stringTest struct {
	key      string
	expected string
}

var stringTests = []stringTest{
	{"Title", "PocoCMS title"},
}

func fmTest(key string, fm map[string]interface{}) string {
	fmt.Println("fmTest()")
	s := fmt.Sprintf("%s", fm[key])
	fmt.Printf("fm[%s] = %s\n", key, s)
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
		actual := fmTest(tt.key, fm)
		if actual != tt.expected {
			t.Errorf("FrontMatter:%s: expected \"%s\", actual \"%s\"", tt.key, actual, tt.expected)
		}
	}
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHelloName(t *testing.T) {
	source := fmAllStr
	if HTML, fm, err := mdYAMLToHTML([]byte(source)); err != nil {
		t.Errorf("Failed converting %s to HTML and obtaining front matter", source)
	} else {
		fmt.Println(fm["Title"], HTML, err)
		s := fm["Title"]
		fmt.Println(s)
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
  actual := titleTag(fm)
	if actual != expected {
		t.Errorf("FrontMatter[\"Title\"]: expected \"%s\", actual \"%s\"", expected, actual)
	}
}

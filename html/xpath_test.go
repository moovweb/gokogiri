package html

import "testing"


func TestUnfoundFuncInXpath(t *testing.T) {
	defer CheckXmlMemoryLeaks(t)

	doc, err := Parse([]byte("<html><body><div><h1></div>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	html := doc.Root().FirstChild()
	results, _ := html.Search("./div[matches(text(), 'foo')]")
	if len(results) != 0 {
		t.Error("should match nothing because the function is not found")
	}
	doc.Free()
}
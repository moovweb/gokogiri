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
	html.Search("./div[matches(text(), 'foo')]")
	doc.Free()
}
package test

import(
	"libxml"
	"libxml/xpath"
	"testing"
)

func TestSearch(t *testing.T) {
	doc := libxml.HtmlParseString("<html><body><div>Hi<div>Mom</div></div></body></html>")
	divs := xpath.Search(doc, "//div")
	// Doctype gets returned as the first child!
	if divs.Size() != 2 {
		t.Error("Returned the two divs!")
	}
	div := divs.NodeAt(0);
	if div.Size() != 1 {
		t.Error("Only has one element in it!")
	}
	textChild := div.First();
	if textChild.Name() != "text" {
		t.Error("Should return a text child")
	}
}
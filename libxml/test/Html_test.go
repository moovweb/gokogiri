package test

import (
	"libxml"
	"testing"
)

func TestHtmlSimpleParse(t *testing.T) {
	doc := libxml.HtmlParseString("<html><head /><body /></html>")
	if doc.Size() != 1 {
		t.Error("Incorrect size")
	}
	// Doctype gets returned as the first child!
	htmlTag := doc.First().Next()
	if htmlTag.Size() != 2 {
		print(htmlTag.Name())
		t.Error("Two tags are inside of <html>")
	}

}

func TestHtmlCDataTag(t *testing.T) {
	doc := libxml.HtmlParseString(LoadFile("docs/script.html"))
	if doc.Size() != 1 {
		t.Error("Incorrect size")
	}
	scriptTag := doc.RootElement().FirstElement().FirstElement()
	if scriptTag.Name() != "script" {
		t.Error("Should have selected the script tag")
	}
	content := scriptTag.Content()
	scriptTag.SetContent(content)
}
package test

import (
	"libxml"
	"libxml/tree"
	"testing"
)

func TestTextNodeContent(t *testing.T) {
	doc := libxml.XmlParseString("<html>hi</html>")
	textNode, ok := doc.First().First().(*tree.Text)
	if !ok {
		t.Error("Should be a Text object")
	}
	if textNode.Content() != "hi" {
		t.Error("Should be equal to the string 'hi'")
	}
	textNode.SetContent("mom")
	if doc.First().String() != "<html>mom</html>" {
		t.Error("Should be able to set text content")
	}
}

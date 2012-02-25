package test

import (
	"gokogiri"
	"gokogiri/tree"
	"testing"
	"strings"
	"gokogiri/help"
)

func TestTextNodeContent(t *testing.T) {
	doc := gokogiri.XmlParseString("<html>hi</html>")
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
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestTextNodeWrap(t *testing.T) {
	doc := gokogiri.XmlParseString("<html>hi</html>")
	textNode, ok := doc.First().First().(*tree.Text)
	if !ok {
		t.Error("Should be a Text object")
	}
	wrapNode := textNode.Wrap("wrapper")
	if wrapNode.Name() != "wrapper" {
		t.Error("Should be <wrapper> node")
	}
	if !strings.Contains(doc.String(), "<wrapper>hi</wrapper>") {
		t.Error("Should have wrapped")
	}
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

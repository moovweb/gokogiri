package test

import (
	"libxml"
	"libxml/help"
	"testing"
)

func TestNewElement(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi</root>")
	root := doc.RootElement()
	child := doc.NewElement("child")
	root.AppendChildNode(child)
	Equal(t, root.String(), "<root>hi<child/></root>")
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestHtmlFragment(t *testing.T) {
	doc := libxml.XmlParseString("<meta name=\"format-detection\" content=\"telephone=no\">")
	root := doc.RootElement()
	child := doc.NewElement("child")
	root.AppendChildNode(child)
	Equal(t, root.String(), "<meta name=\"format-detection\" content=\"telephone=no\"><child/></meta>")
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

package test

import (
	"libxml"
	"libxml/help"
	"testing"
)

func TestDocXmlNoInput(t *testing.T) {
	doc := libxml.XmlParseString("")
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

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

func TestDocParseHtmlFragment(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi</root>")
	root := doc.RootElement()
	fragmentNodes := doc.ParseHtmlFragment("<div><meta style=\"cool\"></div><h1/>")
	for _, node := range(fragmentNodes) {
		root.AppendChildNode(node)
	}
	Equal(t, root.String(), "<root>hi<div><meta style=\"cool\"/></div><h1/></root>")
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}


func TestDocParseHtmlFragmentWithComment(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi</root>")
	fragmentNodes := doc.ParseHtmlFragment("<!-- comment -->")
	Equal(t, len(fragmentNodes), 0)
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}




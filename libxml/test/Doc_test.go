package test

import (
	"gokogiri0/libxml"
	"gokogiri0/libxml/help"
	"gokogiri0/libxml/tree"
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

func TestXmlDocWithComment(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<!-- comments here --></root>")
	root := doc.RootElement()
	Equal(t, root.Content(), "hi<!-- comments here -->")
	Equal(t, root.String(), "<root>hi<!-- comments here --></root>")
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
	child, _ := doc.NewElement("child")
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
	fragmentNodes := doc.ParseHtmlFragment("<div><meta style=\"cool\"></div><h1/>", "")
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
	fragmentNodes := doc.ParseHtmlFragment("<!-- comment -->", "")
	Equal(t, len(fragmentNodes), 1)
	for _, node := range(fragmentNodes) {
		node.Free()
	}
	doc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestInjectAtBottom(t *testing.T) {
	fragment := "<span class='icons-orange-link-arrow'></span>"	
	doc := libxml.XmlParseString("<root></root>")
	nodeSet := doc.ParseHtmlFragment(fragment, "")
	root := doc.RootElement()
	for _, node := range(nodeSet) {
		root.AppendChildNode(node)	
	}
	doc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

func TestWrapThenInject(t *testing.T) {
	fragment := "<span class='icons-orange-link-arrow'></span>"	
	doc := libxml.XmlParseString("<root>hi</root>")
    textNode, ok := doc.First().First().(*tree.Text)
	if !ok {
		t.Error("Should be a Text object")
	}
	wrapNode := textNode.Wrap("span")

	nodeSet := doc.ParseHtmlFragment(fragment, "")
	for _, node := range(nodeSet) {
		wrapNode.AppendChildNode(node)	
	}
	//println("doc:", doc.String())
	doc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

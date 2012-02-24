package test

import (
	"libxml"
	"testing"
	"libxml/help"
)

func TestHtmlFragment(t *testing.T) {
	doc := libxml.XmlParseString("<meta name=\"format-detection\" content=\"telephone=no\">")
	root := doc.RootElement()
	child, err := doc.NewElement("child")
	if err != nil {
		t.Fail()
	}
	root.AppendChildNode(child)
	Equal(t, root.String(), "<meta name=\"format-detection\" content=\"telephone=no\"><child/></meta>")
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestHtmlFragment2(t *testing.T) {
	doc := libxml.HtmlParseFragment("<body><div/></body>")
	f := doc.RootElement().First()
	Equal(t, f.Name(), "body")
	Equal(t, f.First().Name(), "div")
	Equal(t, f.String(), "<body><div/></body>")
	doc.Free()
	
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestHtmlFragment3(t *testing.T) {
	doc := libxml.HtmlParseFragment("<h1><div/></h1>")
	f := doc.RootElement().First()
	Equal(t, f.Name(), "h1")
	Equal(t, f.First().Name(), "div")
	Equal(t, f.String(), "<h1><div/></h1>")
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestHtmlFragmentNewlinesNokogiri(t *testing.T) {
	html := "<script src=\"blah\"></script><div id=\"blah\" class=\" mw_testing\"></div>"
	doc := libxml.HtmlParseFragment(html)
	Equal(t, doc.RootElement().Content(), html)
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

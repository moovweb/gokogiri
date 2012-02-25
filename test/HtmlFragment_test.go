package test

import (
	"gokogiri"
	"testing"
	"gokogiri/help"
)

func TestHtmlFragment(t *testing.T) {
	doc := gokogiri.XmlParseString("<meta name=\"format-detection\" content=\"telephone=no\">")
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

func TestHtmlFragment2(t *testing.T) {
	doc := gokogiri.HtmlParseFragment("<body><div/></body>")
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
	doc := gokogiri.HtmlParseFragment("<h1><div/></h1>")
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
	doc := gokogiri.HtmlParseFragment(html)
	Equal(t, doc.RootElement().Content(), html)
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

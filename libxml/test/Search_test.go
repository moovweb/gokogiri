package test

import (
	"libxml"
  "libxml/help"
	"libxml/xpath"
	"testing"
)

func TestSearch(t *testing.T) {
	doc := libxml.HtmlParseString("<html><body><div>Hi<div>Mom</div></div></body></html>")
  
	divs := xpath.Search(doc, "//div")

	// Doctype gets returned as the first child!
	if divs.Size() != 2 {
		t.Errorf("Returned the two divs!: %d", divs.Size())
	}
  /*
	div := divs.NodeAt(0)
	if div.Size() != 1 {
		t.Error("Only has one element in it!")
	}
	textChild := div.First()
	if textChild.Name() != "text" {
		t.Error("Should return a text child")
	}
  */
    doc.Free()
    
    help.XmlCleanUpParser()
    if help.XmlMemoryAllocation() != 0 {
      t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
      help.XmlMemoryLeakReport()
    }  
}
/*
// What if we remove a node we will soon match?
func TestSearchRemoval(t *testing.T) {
	doc := libxml.XmlParseString("<root><parent><child /></parent></root>")
	root := doc.RootElement()
	nodes := xpath.Search(root, "//*").Slice()
	parent := root.FirstElement()
	parent.SetContent("empty")
	if parent.IsLinked() != true {
		t.Error("Parent starts off linked")
	}
	parent.Remove()
	if parent.IsLinked() != false {
		t.Error("Parent should report being unlinked")
	}
	for i := range nodes {
		node := nodes[i]
		Equal(t, node.Type(), 1)
	}
	doc.Free()

    help.XmlCleanUpParser()
    if help.XmlMemoryAllocation() != 0 {
      t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
      help.XmlMemoryLeakReport()
    }
}

//what if a search returns a nil pointer?
func TestNilSearch(t *testing.T) {
    doc := libxml.XmlParseString("<root id=\"foo\"><h1></h1></root>")
    results := xpath.Search(doc, "//*[@id = 'foo1']//*").Slice()
    if len(results) != 0 {
        t.Error("Should return zero size node set")
    }
    doc.Free()

    help.XmlCleanUpParser()
    if help.XmlMemoryAllocation() != 0 {
      t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
      help.XmlMemoryLeakReport()
    }
}
*/


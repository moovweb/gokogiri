package xpath

import (
	"testing"
)

func NewSampleDoc() (doc *Doc, context *XPathContext) {
	doc = HtmlParse("<html><body><span /><div><span />content</div></body></html>")
	context = doc.XPathContext()
	return
}

func TestXPathRegisterNamespace(t *testing.T) {
	_, context := NewSampleDoc()
	if context.RegisterNamespace("me", "/hi") == false {
		t.Error("Should have been able to register namespace")
	}
}

func TestXPathEvaluation(t *testing.T) {
	_, context := NewSampleDoc()
	nodeSet := context.EvalToNodes("/html/body")
	if nodeSet.Size() != 1 {
		t.Error("Too many elements returned. Maybe some are nil!")
	}
	for i := 0; i < nodeSet.Size(); i++ {
		node := nodeSet.NodeAt(i)
		if node.Name() != "body" {
			t.Error("Nil node returned")
		}
		subnodes := node.Search("//div")
		if subnodes.Size() != 1 {
			t.Error("selected wrong!")
		}
	}
}

func TestXPathNodeSearches(t *testing.T) {
	doc, _ := NewSampleDoc()
	root := doc.RootNode()
	spanSet := root.Search(".//span")
	spans := spanSet.Slice()
	if len(spans) != 2 {
		t.Error("too many spans.. returned ", len(spans), " nodes")
	}
	divSet := root.Search("//div")
	divs := divSet.Slice()
	div := divs[0]

	spanSet = div.Search(".//span")
	spans = spanSet.Slice()
	if len(spans) >= 2 {
		t.Error("Search is NOT scoped: returned ", len(spans), " nodes")
	}
	if len(spans) == 0 {
		t.Error("Should have found SOMETHING. Found nothing.")
	}
}

func TestOnlyElements(t *testing.T) {
	doc, _ := NewSampleDoc()
	nodeSet := doc.Search("//*")
	for i := 0; i < nodeSet.Size(); i++ {
		node := nodeSet.NodeAt(i)
		if (Element)node {
			
		}
	}
}

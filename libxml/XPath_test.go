package libxml

import (
	"testing"
)

func BuildSampleDoc() (doc *XmlDoc, context *XPathContext) {
	doc = HtmlParse("<html><body><span /><div><span />content</div></body></html>")
	context = doc.XPathContext()
	return
}

func TestXPathRegisterNamespace(t *testing.T) {
	_, context := BuildSampleDoc()
	if context.RegisterNamespace("me", "/hi") == false {
		t.Error("Should have been able to register namespace")
	}
}

func TestXPathEvaluation(t *testing.T) {
	_, context := BuildSampleDoc()
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
	doc, _ := BuildSampleDoc()
	root := doc.RootNode()
	span_set := root.Search(".//span")
	spans := span_set.Slice()
	if len(spans) != 2 {
		t.Error("too many spans.. returned ", len(spans), " nodes")
	}
	div_set := root.Search("//div")
	divs := div_set.Slice()
	div := divs[0]

	span_set = div.Search(".//span")
	spans = span_set.Slice()
	if len(spans) >= 2 {
		t.Error("Search is NOT scoped: returned ", len(spans), " nodes")
	}
	if len(spans) == 0 {
		t.Error("Should have found SOMETHING. Found nothing.")
	}
}

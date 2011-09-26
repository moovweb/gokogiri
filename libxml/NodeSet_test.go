package libxml

import (
	"testing"
)

func BuildSampleNodeSetDoc() (doc *XmlDoc, context *XPathContext) {
	doc = HtmlParse("<html><body><span /><div><span /><span>content</span></div></body></html>")
	context = doc.XPathContext()
	return
}

func TestNodeSetRemoval(t *testing.T) {
	doc, _ := BuildSampleNodeSetDoc()
	root := doc.RootNode()
	allNodes := root.Search("//*")
	if allNodes.Size() != 6 {
		t.Error("search sucked")
	}
	spanSet := root.Search("//span")
	spanSet.RemoveAll()
	allNodes = root.Search("//*")
	if allNodes.Size() != 3 {
		t.Error("Should have removed the spans")
	}
}

package libxml

import (
  "testing"
)

func BuildSampleDoc() (doc *XmlDoc, context *XPathContext) {
  doc = HtmlReadDocSimple("<html><body><span /><div><span />content</div></body></html>")
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
  nodes := context.EvalToNodes("/html/body")
  if len(nodes) != 1 {
    t.Error("Too many elements returned. Maybe some are nil!")
  }
  for i := 0; i < len(nodes); i++ {
    node := nodes[i]
    if node.Name() != "body" {
      t.Error("Nil node returned")
    }
    subnodes := node.Search("//div")
    if len(subnodes) != 1 {
      t.Error("selected wrong!")
    }
  }
}

func TestXPathNodeSearches(t *testing.T) {
  doc, _ := BuildSampleDoc()
  root := doc.RootElement()
  spans := root.Search(".//span")
  if len(spans) != 2 {
    t.Error("too many spans.. returned ", len(spans), " nodes")
  }
  divs := root.Search("//div")
  div := divs[0]

  spans = div.Search(".//span")
  if len(spans) >= 2 {
    t.Error("Search is NOT scoped: returned ", len(spans), " nodes")
  }
  if len(spans) == 0 {
    t.Error("Should have found SOMETHING. Found nothing.")
  }
}
package libxml

import (
  "testing"
)

func BuildSampleDoc() *XmlDoc {
  return HtmlReadDocSimple("<html><body><div /><div>content</div></body></html>")
}

func TestXPathContextNamespace(t *testing.T) {
  doc := BuildSampleDoc()
  context := doc.XPathContext()
  if context.RegisterNamespace("me", "/hi") == false {
    t.Error("Should have been able to register namespace")
  }
}
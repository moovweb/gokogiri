package libxml

import (
  "testing"
)

func TestDocBuilding(t *testing.T) {
  //t.Error("hi")
  doc := HtmlReadDocSimple("<html />")
  root := doc.RootElement()
  if root == nil {
    t.Error("Should return a root node")
  }
  child := root.Next()
  if child != nil {
    t.Error("Doesn't have any children")
  }
}

func TestXPathContext(t *testing.T) {
  //doc := HtmlReadDocSimple("<html />")
  //xpath_context := doc.XPathContext()
}

func TestMetaEncoding(t *testing.T) {
  doc := HtmlReadDocSimple("<html />")
  if doc.MetaEncoding() != "" {
    t.Error("No meta encoding should return ''")
  }
  doc = HtmlReadDocSimple("<html><meta http-equiv='Content-Type' content='text/html; charset=utf-8'/></html>")
  if doc.MetaEncoding() != "utf-8" {
    t.Error("Meta data not properly returning")
  }
}



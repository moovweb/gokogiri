package libxml

import (
  "testing"
)

func TestDocBuilding(t *testing.T) {
  //t.Error("hi")
  doc := HtmlReadDocSimple("<html />")
  root := doc.RootNode()
  if root == nil {
    t.Error("Should return a root node")
  }
  child := root.Next()
  if child != nil {
    t.Error("Doesn't have any children")
  }
}

func TestXPathContext(t *testing.T) {
  doc := HtmlReadDocSimple("<html />")
  xpath_context := doc.XPathContext()
  if xpath_context == nil {
    t.Error("Didnt return a valid XPath context")
  }
}

func TestDump(t *testing.T) {
  /*doc := HtmlReadDocSimple("<html><body /></html>")
  if doc.Dump() != "<?xml version=\"1.0\" standalone=\"yes\"?>\n<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.0 Transitional//EN\" \"http://www.w3.org/TR/REC-html40/loose.dtd\">\n<html>\n  <body />\n</html>" {
    println(doc.Dump())
    t.Error("ERROR!")
  }*/
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



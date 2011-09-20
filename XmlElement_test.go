package libxml

import (
  "testing"
)

func TestXmlElementAttributes(t *testing.T) {
  doc := HtmlReadDocSimple("<div id='hi' />")
  root := doc.RootNode()
  div := root.Search("//div").First().(*XmlElement)
  if div.Attribute("id") != "hi" {
		t.Error("looking for id should return 'hi'")
	}
	if div.Attribute("class") != "" {
		t.Error("Non-existant attributes should return nil")
	}
	div.SetAttribute("class", "classy")
	if div.Attribute("class") != "classy" {
		t.Error("Attributes aren't set")
	}
}

func TestXmlElementName(t *testing.T) {
	doc := HtmlReadDocSimple("<div id='hi' />")
  root := doc.RootNode()
  div := root.Search("//div").First()
  if div.Name() != "div" {
		t.Error("Something is wrong with XMLNode.Name()")
	}
	div.SetName("span")
	if div.Name() != "span" {
		t.Error("Something is wrong with XMLNode.SetName()")
	}
}

func TestXmlElementDump(t *testing.T) {
	doc := HtmlReadDocSimple("<div id='hi' />")
  root := doc.RootNode()
  div := root.Search("//div").First()
	result := div.Dump()
	if result != "<div id=\"hi\"/>" {
		t.Error("Node dumping is being... dumpy. Got back this pile of poo: ", result)
	}
	div.SetName("span")
	result = div.Dump()
	if result != "<span id=\"hi\"/>" {
		t.Error("Node dumping is being... dumpy. Got back this pile of poo: ", result)
	}
}

func TestXmlElementRemove(t *testing.T) {
	doc := HtmlReadDocSimple("<html><body><div><span>hi</span></div></body></html>")
  root := doc.RootNode()
  span := root.Search("//span").First()
	span.Remove()
	result := doc.DumpHTML()
	if result != "<!DOCTYPE html PUBLIC \"-//W3C//DTD HTML 4.0 Transitional//EN\" \"http://www.w3.org/TR/REC-html40/loose.dtd\">\n<html><body><div></div></body></html>\n" {
		t.Error("Node dumping is being... dumpy. Got back this pile of poo: ", result)
	}

}
package libxml

import (
  "testing"
)

func TestXmlNodeAttributes(t *testing.T) {
  doc := HtmlReadDocSimple("<div id='hi' />")
  root := doc.RootNode()
  div := root.Search("//div").First()
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

func TestXmlNodeName(t *testing.T) {
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

func TestXmlNodeDump(t *testing.T) {
	
}
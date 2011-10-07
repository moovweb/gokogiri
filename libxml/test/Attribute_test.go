package test

import(
	"libxml"
	"testing"
	"strings"
)

func TestAttributeFetch(t *testing.T) {
	doc := libxml.XmlParseString("<node existing='true' />")
	node := doc.First()
	attribute, didCreate := node.Attribute("existing")
	if attribute == nil {
		t.Fail()
	}
	if didCreate == true {
		t.Error("Should be an existing attribute")
	}
	attribute, didCreate = node.Attribute("created")
	if attribute == nil {
		t.Fail()
	}
	if didCreate == false {
		t.Error("Should be a new attribute")
	}
	if !(strings.Contains(doc.String(), "created=\"\"")) {
		t.Error("Should have the 'created' attr in it")
	}
	attribute, _ = node.Attribute("existing")
	if attribute.Name() != "existing" {
		t.Error("Name isn't working with attributes")
	}
	attribute.SetName("worked")
	if !(strings.Contains(doc.String(), "worked=\"true\"")) {
		t.Error("Should have the 'worked' attr in it")
	}
	if strings.Contains(doc.String(), "existing") {
		t.Error("Existing attribute should be gone now")
	}
	
}
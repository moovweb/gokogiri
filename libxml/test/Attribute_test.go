package test

import(
	"libxml"
	"testing"
	"strings"
)

func TestAttributeFetch(t *testing.T) {
	doc := libxml.XmlParseString("<node existing='true' />")
	node := doc.First()
	exitingAttr, shouldntCreate := node.Attribute("existing")
	if exitingAttr == nil {
		t.Fail()
	}
	if shouldntCreate == true {
		t.Error("Should be an existing attribute")
	}
	createdAttr, didCreate := node.Attribute("created")
	if createdAttr == nil {
		t.Fail()
	}
	if didCreate == false {
		t.Error("Should be a new attribute")
	}
	if !(strings.Contains(doc.String(), "created=\"\"")) {
		t.Error("Should have the 'created' attr in it")
	}

	if exitingAttr.Name() != "existing" {
		t.Error("Name isn't working with attributes")
	}
	
	exitingAttr.SetName("worked") // <node worked="true" created=""/>
	
	if !(strings.Contains(doc.String(), "worked=\"true\"")) {
		t.Error("Should have the 'worked' attr in it")
	}
	if strings.Contains(doc.String(), "existing") {
		t.Error("Existing attribute should be gone now")
	}

	// Remove the created attribute
	createdAttr.Remove() 	//<node worked="true"/>
	if strings.Contains(doc.String(), "created") {
		t.Error("Created attribute should be deleted now")
	}
	
}
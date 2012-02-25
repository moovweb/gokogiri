package test

import (
	"gokogiri"
	"gokogiri/help"
	"testing"
	"strings"
)

func TestAttributeFetch(t *testing.T) {
	doc := gokogiri.XmlParseString("<node existing='true'/>")
	node := doc.First()
	existingAttr, shouldntCreate := node.Attribute("existing")
	if existingAttr == nil {
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

	Equal(t, createdAttr.Content(), "")
	if !(strings.Contains(doc.String(), "created=\"\"")) {
		t.Error("Should have the 'created' attr in it")
	}
	createdAttr.SetContent("yeah")
	Equal(t, createdAttr.Content(), "yeah")
	if !(strings.Contains(doc.String(), "created=\"yeah\"")) {
		t.Logf("createdAttr.IsValid: %v\n", createdAttr.IsValid())
		t.Logf("createdAttr: %q\n", createdAttr.String())
		t.Logf("doc: %q\n", doc.String())
		t.Error("Should have the 'created=\"yeah\"' attr in it")
	}

	Equal(t, existingAttr.Content(), "true")
	if existingAttr.Name() != "existing" {
		t.Error("Name isn't working with attributes")
	}

	existingAttr.SetName("worked") // <node worked="true" created=""/>

	if !(strings.Contains(doc.String(), "worked=\"true\"")) {
		t.Error("Should have the 'worked' attr in it")
	}
	if strings.Contains(doc.String(), "existing") {
		t.Error("Existing attribute should be gone now")
	}

	// Remove the created attribute
	createdAttr.Remove() //<node worked="true"/>
	if strings.Contains(doc.String(), "created") {
		t.Error("Created attribute have been deleted")
	}

	Equal(t, existingAttr.Content(), "true")
	//the string output of a node should not be the same as its content
	Equal(t, existingAttr.String(), " worked=\"true\"")

	existingAttr.SetContent("yes") //<node worked="yes"/>
	Equal(t, existingAttr.Content(), "yes")

	if !strings.Contains(doc.String(), "worked=\"yes\"") {
		t.Error("Should contain yes now")
	}

	doc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

package xml

import (
	"testing"
	"gokogiri/help"
)

func TestParseDocument(t *testing.T) {
	doc, err := Parse([]byte("<foo></foo>"), nil, DefaultEncodingBytes, DefaultParseOption)
	expected := 
`<?xml version="1.0" encoding="utf-8"?>
<foo/>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestEmptyDocument(t *testing.T) {
	doc, err := Parse(nil, nil, DefaultEncodingBytes, DefaultParseOption)
	expected := 
`<?xml version="1.0" encoding="utf-8"?>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match the empty xml")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}
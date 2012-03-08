package xml

import (
	"testing"
	"gokogiri/help"
)

func TestParseDocument(t *testing.T) {
	defer help.CheckXmlMemoryLeaks(t)

	doc, err := Parse([]byte("<foo></foo>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)
	expected := 
`<?xml version="1.0" encoding="utf-8"?>
<foo/>
`
	if err != nil {
		t.Error("parsing error:", err)
	}
	
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	
}

func TestParseDocumentWithBuffer(t *testing.T) {
	defer help.CheckXmlMemoryLeaks(t)
	
	buffer := make([]byte, 100)

	doc, err := ParseWithBuffer([]byte("<foo></foo>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes, buffer)
	expected := 
`<?xml version="1.0" encoding="utf-8"?>
<foo/>
`
	if err != nil {
		t.Error("parsing error:", err)
	}
	
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
}

func TestEmptyDocument(t *testing.T) {
	defer help.CheckXmlMemoryLeaks(t)

	doc, err := Parse(nil, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)
	expected := 
`<?xml version="1.0" encoding="utf-8"?>
`
	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}
	
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match the empty xml")
	}
	doc.Free()
}
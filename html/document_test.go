package html

import (
	"testing"
)

func TestParseDocument(t *testing.T) {
	doc, err := Parse([]byte("<html><body><div><h1></div>"), nil, DefaultEncodingBytes, DefaultParseOption)
	expected := 
`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">
<html><body><div><h1></h1></div></body></html>
`
	expected_xml := 
`<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">
<html><body><div><h1/></div></body></html>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	
	if doc.String() != expected {
		t.Error("the output of the html doc does not match")
	}
	
	if doc.ToXml() != expected_xml {
		t.Error("the xml output of the html doc does not match")
	}

	doc.Free()
	CheckXmlMemoryLeaks(t)
}

func TestEmptyDocument(t *testing.T) {
	doc, err := Parse(nil, nil, DefaultEncodingBytes, DefaultParseOption)
	expected := 
`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">

`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	
	if doc.String() != expected {
		t.Error("the output of the html doc does not match the empty xml")
	}
	doc.Free()
	CheckXmlMemoryLeaks(t)
}


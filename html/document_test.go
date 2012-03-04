package html

import (
	"testing"
	"gokogiri/help"
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
	help.CheckXmlMemoryLeaks(t)
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
		println(doc.String())
		t.Error("the output of the html doc does not match the empty xml")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestParseDocumentFragmentText(t *testing.T) {
	doc, err := Parse(nil, nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := ParseFragment(doc, []byte("ok\n"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	if (len(docFragment.Children) != 1 || docFragment.Children[0].String() != "ok\n") {
		t.Error("the children from the fragment text do not match")
	}
	
	docFragment.Free()
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}

func TestParseDocumentFragment(t *testing.T) {
	doc, err := Parse(nil, nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := ParseFragment(doc, []byte("<div><h1>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	if (len(docFragment.Children) != 1 || docFragment.Children[0].String() != "<div><h1/></div>") {
		t.Error("the of children from the fragment do not match")
	}
	
	docFragment.Free()
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}
package test

import (
	"libxml"
	"testing"
	"strings"
	//"fmt"
)

func TestElementRemove(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	root := doc.RootElement()
	root.Last().First().Remove()
	Equal(t, root.String(), "<root>hi<parent/></root>")
	doc.Free()
}

func TestElementClear(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	root := doc.RootElement()
	root.Clear()
	Equal(t, root.String(), "<root/>")
	doc.Free()
}

func TestElementContent(t *testing.T) {
	contents := "hi<parent><br/></parent>"
	doc := libxml.XmlParseString("<root>" + contents + "</root>")
	root := doc.RootElement()
	//Equal(t, root.Content(), contents)
	root.SetContent("<lonely/>")
	Equal(t, root.First().Name(), "lonely")
	doc.Free()
}

func TestAppendContentUnicode(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	root := doc.RootElement()
	root.AppendContent("<hello>&#x4F60;&#x597D;</hello>")
	//fmt.Printf("%q\n", doc.String());
	if !strings.Contains(doc.String(), "<hello>&#x4F60;&#x597D;</hello></root>") {
		t.Error("Append unicode content failed")
	}
	doc.Free()
}

func TestPrependContentUnicode(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	root := doc.RootElement()
	root.PrependContent("<hello>&#x4F60;&#x597D;</hello>")
	//fmt.Printf("%q\n", doc.String());
	if !strings.Contains(doc.String(), "<root><hello>&#x4F60;&#x597D;</hello>") {
		t.Error("Prepend unicode content failed")
	}
	doc.Free()
}

func TestNoAutocloseContentCall(t *testing.T) {
	doc := libxml.XmlParseString("<root></root>")
	if strings.Contains(doc.Content(), "<root/>") {
		t.Error("Should NOT autoclose tags when using Content!")
	}
	doc.Free()
}

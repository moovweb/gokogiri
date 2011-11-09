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

func TestElementAppendContentUnicode(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	root := doc.RootElement()
	root.AppendContent("<hello>&#x4F60;&#x597D;</hello>")
	//fmt.Printf("%q\n", doc.String());
	if !strings.Contains(doc.String(), "<hello>&#x4F60;&#x597D;</hello></root>") {
		t.Error("Append unicode content failed")
	}
	doc.Free()
}

func TestElementPrependContentUnicode(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	root := doc.RootElement()
	root.PrependContent("<hello>&#x4F60;&#x597D;</hello>")
	//fmt.Printf("%q\n", doc.String());
	if !strings.Contains(doc.String(), "<root><hello>&#x4F60;&#x597D;</hello>") {
		t.Error("Prepend unicode content failed")
	}
	doc.Free()
}

func TestElementNoAutocloseContentCall(t *testing.T) {
	doc := libxml.XmlParseString("<root></root>")
	if strings.Contains(doc.Content(), "<root/>") {
		t.Error("Should NOT autoclose tags when using Content!")
	}
	doc.Free()
}

func TestElementNewChild(t *testing.T) {
	doc := libxml.XmlParseString("<root></root>")
	root := doc.First()
	child := root.NewChild("child", "text")
	Equal(t, child.Name(), "child")
	Equal(t, child.Content(), "text")
	Equal(t, root.String(), "<root><child>text</child></root>")
	root.NewChild("cousin", "")
	Equal(t, root.String(), "<root><child>text</child><cousin></cousin></root>")
	doc.Free()
}


func TestElementWrap(t *testing.T) {
	doc := libxml.XmlParseString("<one/>")
	wrapperNode := doc.First().Wrap("two")
	if wrapperNode.Name() != "two" {
		t.Error("Should have returned a wrapper element")
	}
	Equal(t, doc.First().String(), "<two><one/></two>")
}

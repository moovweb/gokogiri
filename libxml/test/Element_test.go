package test

import(
	"libxml"
	"testing"
	"strings"
  //"fmt"
)

func TestAppendContentUnicode(t *testing.T) {
  doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
  root := doc.RootElement()
  root.AppendContent("<hello>&#x4F60;&#x597D;</hello>")
  //fmt.Printf("%q\n", doc.String());
  if !strings.Contains(doc.String(), "<hello>&#x4F60;&#x597D;</hello></root>") {
    t.Error("Append unicode content failed")
  }
}

func TestPrependContentUnicode(t *testing.T) {
  doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
  root := doc.RootElement()
  root.PrependContent("<hello>&#x4F60;&#x597D;</hello>")
  //fmt.Printf("%q\n", doc.String());
  if !strings.Contains(doc.String(), "<root><hello>&#x4F60;&#x597D;</hello>") {
    t.Error("Prepend unicode content failed")
  }
}

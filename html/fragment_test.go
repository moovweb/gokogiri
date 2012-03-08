package html

import (
	"testing"
	"gokogiri/help"
)


func TestParseDocumentFragmentText(t *testing.T) {
	doc, err := Parse(nil, nil, []byte("iso-8859-1"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := doc.ParseFragment([]byte("ok\r\n"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	if docFragment.Children.Length() != 1 || docFragment.Children.Nodes[0].String() != "ok\r\n" {
		t.Error("the children from the fragment text do not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}

func TestParseDocumentFragment(t *testing.T) {
	doc, err := Parse(nil, nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := doc.ParseFragment([]byte("<div><h1>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	if (len(docFragment.Children.Nodes) != 1 || docFragment.Children.Nodes[0].String() != "<div><h1></h1></div>") {
		t.Error("the of children from the fragment do not match")
	}
	
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}

func TestSearchDocumentFragment(t *testing.T) {
	doc, err := Parse([]byte("<div class='cool'></div>"), nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := doc.ParseFragment([]byte("<div class='cool'><h1>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	if (len(docFragment.Children.Nodes) != 1 || docFragment.Children.Nodes[0].String() != "<div class=\"cool\"><h1></h1></div>") {
		t.Error("the of children from the fragment do not match")
	}

	nodes, err := docFragment.Search(".//*")
	if err != nil {
		t.Error("fragment search has error")
	}
	if len(nodes) != 2 {
		t.Error("the number of children from the fragment does not match")
	}
	nodes, err = docFragment.Search("//div[@class='cool']")

	if err != nil {
		t.Error("fragment search has error")
	}

	if len(nodes) != 1 {
		println(len(nodes))
		for _, node := range(nodes) {
			println(node.String())
		}
		t.Error("the number of children from the fragment's document does not match")
	}
	
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}
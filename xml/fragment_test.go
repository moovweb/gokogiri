package xml

import (
	"testing"
	"gokogiri/help"
)

func TestParseDocumentFragment(t *testing.T) {
	doc, err := Parse(nil, nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := doc.ParseFragment([]byte("<foo></foo><!-- comment here --><bar>fun</bar>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	if (docFragment.Children.Length() != 3) {
		t.Error("the number of children from the fragment does not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}

func TestSearchDocumentFragment(t *testing.T) {
	doc, err := Parse("<moovweb><z/><s/></moovweb>", nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := doc.ParseFragment([]byte("<foo></foo><!-- comment here --><bar>fun</bar>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	nodes, err := docFragment.Search(".//*")
	if err != nil {
		t.Error("fragment search has error")
	}
	if len(nodes) != 2 {
		t.Error("the number of children from the fragment does not match")
	}
	nodes, err = docFragment.Search("//*")

	if err != nil {
		t.Error("fragment search has error")
	}

	if len(nodes) != 0 {
		t.Error("the number of children from the fragment's document does not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}

func TestSearchDocumentFragmentWithEmptyDoc(t *testing.T) {
	doc, err := Parse(nil, nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := doc.ParseFragment([]byte("<foo></foo><!-- comment here --><bar>fun</bar>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	nodes, err := docFragment.Search(".//*")
	if err != nil {
		t.Error("fragment search has error")
	}
	if len(nodes) != 2 {
		t.Error("the number of children from the fragment does not match")
	}
	nodes, err = docFragment.Search("//*")

	if err != nil {
		t.Error("fragment search has error")
	}

	if len(nodes) != 0 {
		t.Error("the number of children from the fragment's document does not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}
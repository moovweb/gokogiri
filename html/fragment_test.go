package html

import (
	"testing"
	"github.com/moovweb/gokogiri/help"
)

func TestParseDocumentFragmentText(t *testing.T) {
	doc, err := Parse(nil, []byte("iso-8859-1"), nil, DefaultParseOption, []byte("iso-8859-1"))
	if err != nil {
		println(err.Error())
	}
	docFragment, err := doc.ParseFragment([]byte("ok\r\n"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(docFragment.Children()) != 1 || docFragment.Children()[0].String() != "ok\r\n" {
		println(docFragment.String())
		t.Error("the children from the fragment text do not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestParseDocumentFragment(t *testing.T) {
	doc, err := Parse(nil, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)
	if err != nil {
		println(err.Error())
	}
	docFragment, err := doc.ParseFragment([]byte("<div><h1>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(docFragment.Children()) != 1 || docFragment.Children()[0].String() != "<div><h1></h1></div>" {
		t.Error("the of children from the fragment do not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)

}

func TestParseDocumentFragment2(t *testing.T) {
	docStr := `<html>
<head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head>
<body>
  </body>
</html>`
	doc, err := Parse([]byte(docStr), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)
	if err != nil {
		println(err.Error())
	}
	docFragment, err := doc.ParseFragment([]byte("<script>cool & fun</script>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(docFragment.Children()) != 1 || docFragment.Children()[0].String() != "<script>cool & fun</script>" {
		t.Error("the of children from the fragment do not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestSearchDocumentFragment(t *testing.T) {
	doc, err := Parse([]byte("<div class='cool'></div>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)
	if err != nil {
		println(err.Error())
	}
	docFragment, err := doc.ParseFragment([]byte("<div class='cool'><h1>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if len(docFragment.Children()) != 1 || docFragment.Children()[0].String() != "<div class=\"cool\"><h1></h1></div>" {
		t.Error("the of children from the fragment do not match")
	}

	nodes, err := docFragment.Search(".//*")
	if err != nil {
		t.Error("fragment search has error")
		return
	}
	if len(nodes) != 2 {
		t.Error("the number of children from the fragment does not match")
	}
	nodes, err = docFragment.Search("//div[@class='cool']")

	if err != nil {
		t.Error("fragment search has error")
		return
	}

	if len(nodes) != 1 {
		println(len(nodes))
		for _, node := range nodes {
			println(node.String())
		}
		t.Error("the number of children from the fragment's document does not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

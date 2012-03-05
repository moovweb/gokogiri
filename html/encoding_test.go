package html

import (
	"testing"
	"io/ioutil"
	"bytes"
	"http"
	"gokogiri/help"
)

func TestParseDocument_CP1252(t *testing.T) {
	httpClient := &http.Client{}
	response, _ := httpClient.Get("http://florist.1800flowers.com/store.php?id=123")
	responseBytes, _ := ioutil.ReadAll(response.Body)
	
	doc, err := Parse(responseBytes, nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println("err:", err.String())
	} else {
		//println("output")
		out := doc.ToHtml([]byte("iso-8859-2"))
		println(len(out))
		index := bytes.IndexByte(out, byte(146))
		println(index)
	}
	
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

/*
func TestParseDocumentFragment(t *testing.T) {
	doc, err := Parse(nil, nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		println(err.String())
	}
	docFragment, err := ParseFragment(doc, []byte("<div><h1>"), nil, DefaultParseOption)
	if err != nil {
		t.Error(err.String())
	}
	if (len(docFragment.Children) != 1 || docFragment.Children[0].String() != "<div><h1></h1></div>") {
		t.Error("the of children from the fragment do not match")
	}
	
	docFragment.Free()
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
	
}
*/
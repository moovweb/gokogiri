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
	
	doc, err := Parse(responseBytes, []byte("windows-1252"), nil, DefaultParseOption, DefaultEncodingBytes)
	if err != nil {
		println("err:", err.String())
		return
	}
	out := doc.String()
	if index := bytes.IndexByte([]byte(out), byte(146)); index >= 0 {
		t.Error("the output is not properly encoded")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestParseDocumentWithInOutEncodings(t *testing.T) {
	httpClient := &http.Client{}
	response, _ := httpClient.Get("http://florist.1800flowers.com/store.php?id=123")
	responseBytes, _ := ioutil.ReadAll(response.Body)
	
	doc, err := Parse(responseBytes, []byte("windows-1252"), nil, DefaultParseOption, []byte("windows-1252"))
	if err != nil {
		println("err:", err.String())
		return
	}
	out := doc.String()
	if index := bytes.IndexByte([]byte(out), byte(146)); index < 0 {
		t.Error("the output is not properly encoded")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}
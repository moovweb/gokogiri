package html

import (
	"testing"
	"io/ioutil"
	"bytes"
	"github.com/moovweb/gokogiri/help"
)

func TestParseDocument_CP1252(t *testing.T) {
	input, err := ioutil.ReadFile("./tests/document/encoding/input.html")
	if err != nil {
		t.Error("err:", err.Error())
		return
	}
	doc, err := Parse(input, []byte("windows-1252"), nil, DefaultParseOption, DefaultEncodingBytes)
	if err != nil {
		t.Error("err:", err.Error())
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
	input, err := ioutil.ReadFile("./tests/document/encoding/input.html")
	if err != nil {
		t.Error("err:", err.Error())
		return
	}
	doc, err := Parse(input, []byte("windows-1252"), nil, DefaultParseOption, []byte("windows-1252"))
	if err != nil {
		t.Error("err:", err.Error())
		return
	}
	out := doc.String()
	if index := bytes.IndexByte([]byte(out), byte(146)); index < 0 {
		t.Error("the output is not properly encoded")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

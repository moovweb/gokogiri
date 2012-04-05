package html

import (
	"testing"
	"gokogiri/help"
)

func TestSetAttribute(t *testing.T) {
	defer help.CheckXmlMemoryLeaks(t)

	doc, err := Parse([]byte("<html><head></head><body><div><h1></div>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	head := doc.Root().FirstChild()
	println(head.String())

	meta := head.FirstChild()
	//println(meta.String())
	println(meta)
	doc.Free()
}

package test

import (
	"libxml"
	//"libxml/tree"
	"testing"
	"strings"
)


func TestCDataNodeSetContent(t *testing.T) {
	doc := libxml.HtmlParseString("<html></html>")
	htmlNode := doc.RootElement()
	scriptNode := htmlNode.NewChild("script", "")
	scriptNode.SetCDataContent("//<![CDATA[\nalert('boo')\n//]]>")
	if !strings.Contains(doc.DumpHTML(), "<script>//<![CDATA[") {
		t.Error("Should have actually made a CDATA tag")
	}
}
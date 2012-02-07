package test


import (
	//"libxml"
	"libxml/help"
	"libxml/tree"
	"testing"
	"io/ioutil"
)


func TestEncodingRead(t *testing.T) {
	docContent, err := ioutil.ReadFile("htmldata/google-cn.html")
	if err != nil {
		t.Errorf("Err: %v", err.String())
	}
	doc := tree.HtmlParseStringWithOptions(string(docContent), "", "utf8", tree.DefaultHtmlParseOptions())
	root := doc.RootElement()
	head := root.FirstElement()
	body := head.Next()
	title := head.First().Next().Content()
	licenseCode := body.First().Next().Next().Next().First().Next().Content()
	
	if title != "Google" {
		t.Errorf("the English string does not match")
	}
	
	if licenseCode != "ICP证合字B2-20070004号" {
		t.Errorf("the English & Chinese string does not match")
	}
	doc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}


func TestEncodingSetContent(t *testing.T) {
	docContent, err := ioutil.ReadFile("htmldata/google-cn.html")
	if err != nil {
		t.Errorf("Err: %v", err.String())
	}
	doc := tree.HtmlParseStringWithOptions(string(docContent), "", "utf8", tree.DefaultHtmlParseOptions())
	root := doc.RootElement()
	head := root.FirstElement()
	body := head.Next()
	title := head.First().Next().Content()
	licenseCode := body.First().Next().Next().Next().First().Next().Content()
	
	if title != "Google" {
		t.Errorf("the English string does not match")
	}
	
	if licenseCode != "ICP证合字B2-20070004号" {
		t.Errorf("the English & Chinese string does not match")
	}
	
	newStr := "你好，Moovweb"
	body.SetContent("<p>"+newStr+"</p>")
	newContent := body.First().Content()
	
	if newContent != newStr {
		t.Errorf("the new content of English and Chinese does not match after setting content")
	}
	
	doc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}
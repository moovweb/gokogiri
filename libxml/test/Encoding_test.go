package test


import (
	//"libxml"
	"gokogiri/libxml/help"
	"gokogiri/libxml/tree"
	"testing"
	"io/ioutil"
	//"strings"
)


func TestEncodingRead(t *testing.T) {
	docContent, err := ioutil.ReadFile("htmldata/google-cn.html")
	if err != nil {
		t.Errorf("Err: %v", err.String())
	}
	doc := tree.HtmlParseString(string(docContent), "utf8")
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
	doc := tree.HtmlParseString(string(docContent), "utf8")
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

func TestEncodingHTMLFragment(t *testing.T) {
	content, err := ioutil.ReadFile("htmldata/fragment.html")
	if err != nil {
		t.Errorf("Err: %v\n", err.String())
	}

	doc := tree.HtmlParseFragment(string(content), "utf-8") // vs "utf8" ?
	newContent := doc.RootElement()
	rawContent := newContent.Content()
	
	if rawContent != string(content) {
		t.Errorf("Result of HtmlParseFragment() : [%v]\ndoesn't match original content: [%v]", newContent, string(content))
	}
	doc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestEncodingAppendChild(t *testing.T) {
	// Within the tree.HtmlParseFragment() method, using just an xml doc or just an html doc works,
	// but appending the html nodes to the xml doc breaks encoding somehow.

	content, err := ioutil.ReadFile("htmldata/fragment.html")
	if err != nil {
		t.Errorf("Err: %v\n", err.String())
	}
	encoding := "utf-8"
	controlXmlDoc := tree.XmlParseString("<root>" + string(content) + "</root>", encoding)
	testXmlDoc := tree.XmlParseString("<root></root>", encoding)

	htmlDoc := tree.HtmlParseStringWithOptions("<html><body>"+string(content), "", encoding, tree.DefaultHtmlParseOptions())
	/*
	tmpNode := htmlDoc.RootElement().First()
	if strings.Index(strings.ToLower(string(content)), "<body") < 0 {
		tmpNode = tmpNode.First()
	}

	//append all children of tmpRoot to root.
	root := testXmlDoc.RootElement()
	child := tmpNode

	for child != nil {
		nextChild := child.Next()
		root.AppendChildNode(child)
		child = nextChild
	}


	// Now compare the differences:

	testXmlContent := testXmlDoc.RootElement().Content()
	controlXmlContent := controlXmlDoc.RootElement().Content()

	if testXmlContent != controlXmlContent {
		t.Errorf("Result of appending children messed up the encoding.\nPure xml: [%v]\nXml w appended nodes:[%v]\n", controlXmlContent, testXmlContent)
	}


	if string(content) != testXmlContent {
		t.Errorf("Result of Appending children : [%v]\ndoesn't match original content: [%v]", testXmlContent, string(content))
	}
	*/
	htmlDoc.Free()
	testXmlDoc.Free()
	controlXmlDoc.Free()
	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

package test

import (
	"gokogiri/libxml"
	"gokogiri/libxml/tree"
	"testing"
	"strings"
	"gokogiri/libxml/help"
	"fmt"
)

func TestTextNodeContent(t *testing.T) {
	doc := libxml.XmlParseString("<html>hi</html>")
	textNode, ok := doc.First().First().(*tree.Text)
	if !ok {
		t.Error("Should be a Text object")
	}
	if textNode.Content() != "hi" {
		t.Error("Should be equal to the string 'hi'")
	}
	textNode.SetContent("mom")
	if doc.First().String() != "<html>mom</html>" {
		t.Error("Should be able to set text content")
	}
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestTextNodeWrap(t *testing.T) {
	doc := libxml.XmlParseString("<html>hi</html>")
	textNode, ok := doc.First().First().(*tree.Text)
	if !ok {
		t.Error("Should be a Text object")
	}
	wrapNode := textNode.Wrap("wrapper")
	if wrapNode.Name() != "wrapper" {
		t.Error("Should be <wrapper> node")
	}
	if !strings.Contains(doc.String(), "<wrapper>hi</wrapper>") {
		t.Error("Should have wrapped")
	}
	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestTextNodeContentWCarriageReturns(t *testing.T) {

	innerContent := "800.867.5309"
	rawContent := fmt.Sprintf("<html>\r\n%v\r\n</html>", innerContent)

	doc := tree.HtmlParseString(rawContent, "UTF-8")
	//doc.SetMetaEncoding(outputEncoding)

	textNode := doc.RootElement()
	actualContent := textNode.String()

	if actualContent != innerContent {
		t.Errorf("Should be equal to the string:[%v]\n Got:[%v]\n", innerContent, actualContent)
	}

	doc.Free()

	help.XmlCleanUpParser()
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}
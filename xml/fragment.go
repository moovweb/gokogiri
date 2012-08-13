package xml

//#include "helper.h"
import "C"
import (
	. "gokogiri/util"
)

type DocumentFragment struct {
	Node
}

var (
	fragmentWrapperStart = []byte("<root>")
	fragmentWrapperEnd   = []byte("</root>")
)

var fragmentParser = &XmlFragmentParser{}

const initChildrenNumber = 4

func (fragment *DocumentFragment) Remove() {
	fragment.Node.Remove()
}

func (fragment *DocumentFragment) Children() []Node {
	nodes := make([]Node, 0, initChildrenNumber)
	child := fragment.FirstChild()
	for ; child != nil; child = child.NextSibling() {
		nodes = append(nodes, child)
	}
	return nodes
}

func (fragment *DocumentFragment) ToBuffer(outputBuffer []byte) []byte {
	var b []byte
	var size int
	for _, node := range fragment.Children() {
		if docType := node.DocType(); docType == XML_HTML_DOCUMENT_NODE {
			b, size = node.ToHtml(node.OutputEncoding(), nil)
		} else {
			b, size = node.ToXml(node.OutputEncoding(), nil)
		}
		outputBuffer = append(outputBuffer, b[:size]...)
	}
	return outputBuffer
}

func (fragment *DocumentFragment) String() string {
	b := fragment.ToBuffer(nil)
	if b == nil {
		return ""
	}
	return string(b)
}

func ParseFragment(content, inEncoding, url []byte, options int, outEncoding []byte) (fragment *DocumentFragment, err error) {
	inEncoding = AppendCStringTerminator(inEncoding)
	outEncoding = AppendCStringTerminator(outEncoding)
	document := CreateEmptyDocument(inEncoding, outEncoding)
	fragment, err = document.ParseFragment(content, url, options)
	return
}
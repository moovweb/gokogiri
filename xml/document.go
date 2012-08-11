package xml

/*
#cgo CFLAGS: -I../../../clibs/include/libxml2
#cgo LDFLAGS: -lxml2 -L../../../clibs/lib

#include "helper.h"
*/
import "C"

import (
	. "gokogiri/util"
)

import (
	"runtime"
	"unsafe"
)

type Document interface {
	CreateElementNode(string) *ElementNode
	CreateCDataNode(string) *CDataNode
	CreateTextNode(string) *TextNode
	//CreateCommentNode(string) *CommentNode

	DocType() int
	Free()
	String() string
	Root() *ElementNode
}

type XmlDocument struct {
	Ptr *C.xmlDoc
	*XmlNode
}

//create a document
func NewDocument(p unsafe.Pointer, contentLen int, inEncoding, outEncoding []byte) (doc *XmlDocument) {
	docCtx := NewDocCtx(p, inEncoding, outEncoding)
	xmlNode := &XmlNode{Ptr: (*C.xmlNode)(p), DocCtx: docCtx}
	docPtr := (*C.xmlDoc)(p)
	doc = &XmlDocument{Ptr: docPtr, XmlNode: xmlNode}
	
	runtime.SetFinalizer(doc, (*XmlDocument).Free)
	return
}

func Parse(content, inEncoding, url []byte, options int, outEncoding []byte) (doc *XmlDocument, err error) {
	inEncoding = AppendCStringTerminator(inEncoding)
	outEncoding = AppendCStringTerminator(outEncoding)

	var docPtr *C.xmlDoc
	contentLen := len(content)

	if contentLen > 0 {
		var contentPtr, urlPtr, encodingPtr unsafe.Pointer
		contentPtr = unsafe.Pointer(&content[0])

		if len(url) > 0 {
			url = AppendCStringTerminator(url)
			urlPtr = unsafe.Pointer(&url[0])
		}
		if len(inEncoding) > 0 {
			encodingPtr = unsafe.Pointer(&inEncoding[0])
		}

		docPtr = C.xmlParse(contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)

		if docPtr == nil {
			err = ERR_FAILED_TO_PARSE_XML
		} else {
			doc = NewDocument(unsafe.Pointer(docPtr), contentLen, inEncoding, outEncoding)
		}

	} else {
		doc = CreateEmptyDocument(inEncoding, outEncoding)
	}
	return
}

func CreateEmptyDocument(inEncoding, outEncoding []byte) (doc *XmlDocument) {
	docPtr := C.newEmptyXmlDoc()
	doc = NewDocument(unsafe.Pointer(docPtr), 0, inEncoding, outEncoding)
	return
}

func (doc *XmlDocument) Root() (element *ElementNode) {
	nodePtr := C.xmlDocGetRootElement(doc.Ptr)
	if nodePtr != nil {
		element = NewNode(unsafe.Pointer(nodePtr), doc.DocCtx).(*ElementNode)
	}
	return
}

func (doc *XmlDocument) CreateElementNode(tag string) (element *ElementNode) {
	tagBytes := GetCString([]byte(tag))
	tagPtr := unsafe.Pointer(&tagBytes[0])
	newNodePtr := C.xmlNewNode(nil, (*C.xmlChar)(tagPtr))
	newNode := NewNode(unsafe.Pointer(newNodePtr), doc.DocCtx)
	element = newNode.(*ElementNode)
	return
}

func (doc *XmlDocument) CreateTextNode(data string) (text *TextNode) {
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewText((*C.xmlChar)(dataPtr))
	if nodePtr != nil {
		nodePtr.doc = (*_Ctype_struct__xmlDoc)(doc.DocPtr)
		text = NewNode(unsafe.Pointer(nodePtr), doc.DocCtx).(*TextNode)
	}
	return
}

func (doc *XmlDocument) CreateCDataNode(data string) (cdata *CDataNode) {
	dataLen := len(data)
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewCDataBlock(doc.Ptr, (*C.xmlChar)(dataPtr), C.int(dataLen))
	if nodePtr != nil {
		cdata = NewNode(unsafe.Pointer(nodePtr), doc.DocCtx).(*CDataNode)
	}
	return
}
/*
func (doc *XmlDocument) CreateCommentNode(data string) (cdata *CommentNode) {
	dataLen := len(data)
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewCDataBlock(doc.Ptr, (*C.xmlChar)(dataPtr), C.int(dataLen))
	if nodePtr != nil {
		cdata = NewNode(unsafe.Pointer(nodePtr), doc.DocCtx).(*CDataNode)
	}
	return
}
*/
package xml

/*
#cgo pkg-config: libxml-2.0

#include "helper.h"
*/
import "C"

import (
	"unsafe"
	"os"
	"gokogiri/xpath"
	//	"runtime/debug"
)

type Document interface {
	DocPtr() unsafe.Pointer
	DocType() int
	InputEncoding() []byte
	OutputEncoding() []byte
	DocXPathCtx() *xpath.XPath
	AddUnlinkedNode(unsafe.Pointer)
	ParseFragment([]byte, []byte, int) (*DocumentFragment, os.Error)
	CreateElementNode(string) *ElementNode
	Free()
	String() string
	BookkeepFragment(*DocumentFragment)
}

//xml parse option
const (
	XML_PARSE_RECOVER   = 1 << 0  //relaxed parsing
	XML_PARSE_NOERROR   = 1 << 5  //suppress error reports 
	XML_PARSE_NOWARNING = 1 << 6  //suppress warning reports 
	XML_PARSE_NONET     = 1 << 11 //forbid network access
)

//default parsing option: relax parsing
var DefaultParseOption = XML_PARSE_RECOVER |
	XML_PARSE_NONET |
	XML_PARSE_NOERROR |
	XML_PARSE_NOWARNING

//xml save option
const (
	XML_SAVE_FORMAT   = 1 << 0 /* format save output */
	XML_SAVE_NO_DECL  = 1 << 1 /* drop the xml declaration */
	XML_SAVE_NO_EMPTY = 1 << 2 /* no empty tags */
	XML_SAVE_NO_XHTML = 1 << 3 /* disable XHTML1 specific rules */
	XML_SAVE_XHTML    = 1 << 4 /* force XHTML1 specific rules */
	XML_SAVE_AS_XML   = 1 << 5 /* force XML serialization on HTML doc */
	XML_SAVE_AS_HTML  = 1 << 6 /* force HTML serialization on XML doc */
	XML_SAVE_WSNONSIG = 1 << 7 /* format with non-significant whitespace */
)

//libxml2 use "utf-8" by default, and so do we
const DefaultEncoding = "utf-8"

var ERR_FAILED_TO_PARSE_XML = os.NewError("failed to parse xml input")

type XmlDocument struct {
	Ptr *C.xmlDoc
	*XmlNode
	InEncoding    []byte
	OutEncoding   []byte
	UnlinkedNodes []unsafe.Pointer
	XPathCtx      *xpath.XPath
	Type          int

	fragments []*DocumentFragment //save the pointers to free them when the doc is freed
}

//default encoding in byte slice
var DefaultEncodingBytes = []byte(DefaultEncoding)

const initialUnlinkedNodes = 8
const initialFragments = 2

//create a document
func NewDocument(p unsafe.Pointer, contentLen int, inEncoding, outEncoding, buffer []byte) (doc *XmlDocument) {
	xmlNode := &XmlNode{Ptr: (*C.xmlNode)(p)}
	adjustedLen := contentLen + contentLen>>1  //1.5 of the input len
	if adjustedLen < initialOutputBufferSize { //min len
		adjustedLen = initialOutputBufferSize
	}
	if len(buffer) < adjustedLen {
		xmlNode.outputBuffer = make([]byte, adjustedLen)
	} else {
		xmlNode.outputBuffer = buffer
	}
	docPtr := (*C.xmlDoc)(p)
	doc = &XmlDocument{Ptr: docPtr, XmlNode: xmlNode, InEncoding: inEncoding, OutEncoding: outEncoding}
	doc.UnlinkedNodes = make([]unsafe.Pointer, 0, initialUnlinkedNodes)
	doc.XPathCtx = xpath.NewXPath(p)
	doc.Type = xmlNode.NodeType()
	doc.fragments = make([]*DocumentFragment, 0, initialFragments)
	xmlNode.Document = doc
	return
}

func ParseWithBuffer(content, inEncoding, url []byte, options int, outEncoding, outBuffer []byte) (doc *XmlDocument, err os.Error) {
	var docPtr *C.xmlDoc
	contentLen := len(content)

	if contentLen > 0 {
		var contentPtr, urlPtr, encodingPtr unsafe.Pointer

		contentPtr = unsafe.Pointer(&content[0])
		if len(url) > 0 {
			urlPtr = unsafe.Pointer(&url[0])
		}
		if len(inEncoding) > 0 {
			encodingPtr = unsafe.Pointer(&inEncoding[0])
		}

		docPtr = C.xmlParse(contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)

		if docPtr == nil {
			err = ERR_FAILED_TO_PARSE_XML
		} else {
			doc = NewDocument(unsafe.Pointer(docPtr), contentLen, inEncoding, outEncoding, outBuffer)
		}

	} else {
		doc = CreateEmptyDocument(inEncoding, outEncoding, outBuffer)
	}
	return
}

//parse a string to document
func Parse(content, inEncoding, url []byte, options int, outEncoding []byte) (doc *XmlDocument, err os.Error) {
	doc, err = ParseWithBuffer(content, inEncoding, url, options, outEncoding, nil)
	return
}

func CreateEmptyDocument(inEncoding, outEncoding, outBuffer []byte) (doc *XmlDocument) {
	docPtr := C.newEmptyXmlDoc()
	doc = NewDocument(unsafe.Pointer(docPtr), 0, inEncoding, outEncoding, outBuffer)
	return
}

func (document *XmlDocument) ParseFragment(input, url []byte, options int) (fragment *DocumentFragment, err os.Error) {
	fragment, err = parsefragment(document, input, document.InputEncoding(), url, options)
	return
}

func (document *XmlDocument) DocPtr() (ptr unsafe.Pointer) {
	ptr = unsafe.Pointer(document.Ptr)
	return
}

func (document *XmlDocument) DocType() (t int) {
	t = document.Type
	return
}

func (document *XmlDocument) InputEncoding() (encoding []byte) {
	encoding = document.InEncoding
	return
}

func (document *XmlDocument) OutputEncoding() (encoding []byte) {
	encoding = document.OutEncoding
	return
}

func (document *XmlDocument) DocXPathCtx() (ctx *xpath.XPath) {
	ctx = document.XPathCtx
	return
}

func (document *XmlDocument) AddUnlinkedNode(nodePtr unsafe.Pointer) {
	document.UnlinkedNodes = append(document.UnlinkedNodes, nodePtr)
}

func (document *XmlDocument) BookkeepFragment(fragment *DocumentFragment) {
	document.fragments = append(document.fragments, fragment)
}

func (document *XmlDocument) Root() (element *ElementNode) {
	nodePtr := C.xmlDocGetRootElement(document.Ptr)
	element = NewNode(unsafe.Pointer(nodePtr), document).(*ElementNode)
	return
}

func (document *XmlDocument) CreateElementNode(tag string) (element *ElementNode) {
	var tagPtr unsafe.Pointer
	if len(tag) > 0 {
		tagBytes := []byte(tag)
		tagPtr = unsafe.Pointer(&tagBytes[0])
	}
	newNodePtr := C.xmlNewNode(nil, (*C.xmlChar)(tagPtr))
	newNode := NewNode(unsafe.Pointer(newNodePtr), document)
	element = newNode.(*ElementNode)
	return
}

/*
func (document *XmlDocument) ToXml() string {
	document.outputOffset = 0
	objPtr := unsafe.Pointer(document.XmlNode)
	nodePtr      := unsafe.Pointer(document.Ptr)
	encodingPtr := unsafe.Pointer(&(document.Encoding[0]))
	C.xmlSaveNode(objPtr, nodePtr, encodingPtr, XML_SAVE_AS_XML)
	return string(document.outputBuffer[:document.outputOffset])
}

func (document *XmlDocument) ToHtml() string {
	document.outputOffset = 0
	documentPtr := unsafe.Pointer(document.XmlNode)
	docPtr      := unsafe.Pointer(document.Ptr)
	encodingPtr := unsafe.Pointer(&(document.Encoding[0]))
	C.xmlSaveNode(documentPtr, docPtr, encodingPtr, XML_SAVE_AS_HTML)
	return string(document.outputBuffer[:document.outputOffset])
}

func (document *XmlDocument) ToXml2() string {
	encodingPtr := unsafe.Pointer(&(document.Encoding[0]))
	charPtr := C.xmlDocDumpToString(document.Ptr, encodingPtr, 0)
	defer C.xmlFreeChars(charPtr)
	return C.GoString(charPtr)
}

func (document *XmlDocument) ToHtml2() string {
	charPtr := C.htmlDocDumpToString(document.Ptr, 0)
	defer C.xmlFreeChars(charPtr)
	return C.GoString(charPtr)
}

func (document *XmlDocument) String() string {
	return document.ToXml()
}
*/
func (document *XmlDocument) Free() {
	//must clear the fragments first
	//because the nodes are put in the unlinked list
	for _, fragment := range document.fragments {
		fragment.Remove()
	}

	for _, nodePtr := range document.UnlinkedNodes {
		C.xmlFreeNode((*C.xmlNode)(nodePtr))
	}

	document.XPathCtx.Free()
	C.xmlFreeDoc(document.Ptr)
}

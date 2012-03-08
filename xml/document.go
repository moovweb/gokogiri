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
)

type Document interface {
	DocPtr() unsafe.Pointer
	DocType() int
	DocEncoding() []byte
	DocXPathCtx() *xpath.XPath
	AddUnlinkedNode(unsafe.Pointer)
	ParseFragment([]byte, []byte, int) (Document, os.Error)
	Free()
}

//xml parse option
const (
	XML_PARSE_RECOVER   = 1 << 0 //relaxed parsing
    XML_PARSE_NOERROR   = 1 << 5  //suppress error reports 
    XML_PARSE_NOWARNING = 1 << 6  //suppress warning reports 
    XML_PARSE_NONET     = 1 << 11 //forbid network access
)

//default parsing option: relax parsing
var DefaultParseOption = 	XML_PARSE_RECOVER | 
    						XML_PARSE_NONET|
    						XML_PARSE_NOERROR|
    						XML_PARSE_NOWARNING

//xml save option
const (
	XML_SAVE_FORMAT     = 1<<0	/* format save output */
	XML_SAVE_NO_DECL    = 1<<1	/* drop the xml declaration */
	XML_SAVE_NO_EMPTY	= 1<<2 /* no empty tags */
	XML_SAVE_NO_XHTML	= 1<<3 /* disable XHTML1 specific rules */
	XML_SAVE_XHTML	    = 1<<4 /* force XHTML1 specific rules */
	XML_SAVE_AS_XML     = 1<<5 /* force XML serialization on HTML doc */
	XML_SAVE_AS_HTML    = 1<<6 /* force HTML serialization on XML doc */
	XML_SAVE_WSNONSIG   = 1<<7  /* format with non-significant whitespace */
)

//libxml2 use "utf-8" by default, and so do we
const DefaultEncoding = "utf-8"

type XmlDocument struct {
	Ptr *C.xmlDoc
	*XmlNode	
	Encoding []byte
	UnlinkedNodes []unsafe.Pointer
	XPathCtx *xpath.XPath
	Type int
}

//default encoding in byte slice
var DefaultEncodingBytes = []byte(DefaultEncoding)

const initialUnlinkedNodes = 8
//create a document
func NewDocument(p unsafe.Pointer, encoding []byte, buffer []byte) (doc *XmlDocument) {
	xmlNode := &XmlNode{Ptr: (*C.xmlNode)(p)}
	if len(buffer) == 0 {
		xmlNode.outputBuffer = make([]byte, initialOutputBufferSize)
	}
	docPtr := (*C.xmlDoc)(p)
	doc = &XmlDocument{Ptr: docPtr, XmlNode: xmlNode, Encoding: encoding}
	doc.UnlinkedNodes = make([]unsafe.Pointer, 0, initialUnlinkedNodes)
	doc.XPathCtx = xpath.NewXPath(p) 
	doc.Type = xmlNode.NodeType()
	xmlNode.Document = doc
	return
}

//parse a string to document
func Parse(content, url, encoding []byte, options int) (doc *XmlDocument, err os.Error) {
	var docPtr *C.xmlDoc
	contentLen := len(content)
	
	if contentLen > 0 {
		var contentPtr, urlPtr, encodingPtr unsafe.Pointer
		
		contentPtr   = unsafe.Pointer(&content[0])
		if len(url) > 0      { urlPtr       = unsafe.Pointer(&url[0]) }
		if len(encoding) > 0 { encodingPtr  = unsafe.Pointer(&encoding[0]) }
		
		docPtr = C.xmlParse(contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)
	} else {
		doc = CreateEmptyDocument(encoding)
	}
	return
}

func CreateEmptyDocument(encoding []byte) (doc *XmlDocument) {
	docPtr := C.newEmptyXmlDoc()
	doc = NewDocument(unsafe.Pointer(docPtr), encoding, nil)
	return
}

func (document *XmlDocument) ParseFragment(input, url []byte, options int) (doc Document, err os.Error) {
	content = append(fragmentWrapperStart, content...)
	content = append(content, fragmentWrapperEnd...)

	var contentPtr, urlPtr unsafe.Pointer
	contentPtr   = unsafe.Pointer(&content[0])
	contentLen   := len(content)
	if len(url) > 0  { urlPtr = unsafe.Pointer(&url[0]) }
	
	rootElementPtr := C.xmlParseFragment(document.DocPtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)

	if rootElementPtr == nil { err = ErrFailParseFragment; return }
	
	c := (*C.xmlNode)(unsafe.Pointer(rootElementPtr.children))
	var nextSibling *C.xmlNode
	
	for ; c != nil; c = nextSibling {
		nextSibling = (*C.xmlNode)(unsafe.Pointer(c.next))
		C.xmlUnlinkNode(c)
		fragment.Children = append(fragment.Children, NewNode(unsafe.Pointer(c), document))
	}
	//now we have rip all its children nodes, we should release the root node
	C.xmlFreeNode(rootElementPtr)
}

func (document *XmlDocument) DocPtr() (ptr unsafe.Pointer) {
	ptr = unsafe.Pointer(document.Ptr)
	return
}

func (document *XmlDocument) DocType() (t int) {
	t = document.Type
	return
}

func (document *XmlDocument) DocEncoding() (encoding []byte) {
	encoding = document.Encoding
	return
}

func (document *XmlDocument) DocXPathCtx() (ctx *xpath.XPath) {
	ctx = document.XPathCtx
	return
}

func (document *XmlDocument) AddUnlinkedNode(nodePtr unsafe.Pointer) {
	document.UnlinkedNodes = append(document.UnlinkedNodes, nodePtr)
}

func (document *XmlDocument) Root() (element *ElementNode) {
	nodePtr := C.xmlDocGetRootElement(document.Ptr)
	element = NewNode(unsafe.Pointer(nodePtr), document).(*ElementNode)
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
	for _, nodePtr := range(document.UnlinkedNodes) {
		C.xmlFreeNode((*C.xmlNode)(nodePtr))
	}
	document.XPathCtx.Free()
	C.xmlFreeDoc(document.Ptr)
}

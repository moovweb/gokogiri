package xml

/*
#cgo pkg-config: libxml-2.0

#include "chelper.h"
*/
import "C"

import (
	"unsafe"
)

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

type Document struct {
	DocPtr *C.xmlDoc
	*XmlNode
	
	Encoding []byte
}

//default encoding in byte slice
var DefaultEncodingBytes = []byte(DefaultEncoding)

//create a document
func NewDocument(p unsafe.Pointer, encoding []byte, buffer []byte) (doc *Document) {
	xmlNode := &XmlNode{NodePtr: (*C.xmlNode)(p)}
	if len(buffer) == 0 {
		xmlNode.outputBuffer = make([]byte, initialOutputBufferSize)
	}
	docPtr := (*C.xmlDoc)(p)
	doc = &Document{DocPtr: docPtr, XmlNode: xmlNode, Encoding: encoding}
	xmlNode.Document = doc
	return
}

//parse a string to document
func Parse(content, url, encoding []byte, options int) (doc *Document, err error) {
	var docPtr *C.xmlDoc
	contentLen := len(content)
	
	if contentLen > 0 {
		var contentPtr, urlPtr, encodingPtr unsafe.Pointer
		
		contentPtr   = unsafe.Pointer(&content[0])
		if len(url) > 0      { urlPtr       = unsafe.Pointer(&url[0]) }
		if len(encoding) > 0 { encodingPtr  = unsafe.Pointer(&encoding[0]) }
		
		docPtr = C.xmlParse(contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)
	}
	if docPtr == nil {
		//why does newEmptyXmlDoc NOT call xmlInitParser like other parse functions?
		C.xmlInitParser();
		docPtr = C.newEmptyXmlDoc()
	}
	doc = NewDocument(unsafe.Pointer(docPtr), encoding, nil)
	return
}

func (document *Document) RootElement() (element *ElementNode) {
	nodePtr := C.xmlDocGetRootElement(document.DocPtr)
	element = NewNode(nodePtr, document).(*ElementNode)
	return
}

func (document *Document) ToXml() string {
	document.outputOffset = 0
	objPtr := unsafe.Pointer(document.XmlNode)
	nodePtr      := unsafe.Pointer(document.DocPtr)
	encodingPtr := unsafe.Pointer(&(document.Encoding[0]))
	C.xmlSaveNode(objPtr, nodePtr, encodingPtr, XML_SAVE_AS_XML)
	return string(document.outputBuffer[:document.outputOffset])
}

func (document *Document) ToXml2() string {
	encodingPtr := unsafe.Pointer(&(document.Encoding[0]))
	charPtr := C.xmlDocDumpToString(document.DocPtr, encodingPtr, 0)
	defer C.xmlFreeChars(charPtr)
	return C.GoString(charPtr)
}

func (document *Document) ToHtml() string {
	document.outputOffset = 0
	documentPtr := unsafe.Pointer(document.XmlNode)
	docPtr      := unsafe.Pointer(document.DocPtr)
	encodingPtr := unsafe.Pointer(&(document.Encoding[0]))
	C.xmlSaveNode(documentPtr, docPtr, encodingPtr, XML_SAVE_AS_HTML)
	return string(document.outputBuffer[:document.outputOffset])
}

func (document *Document) ToHtml2() string {
	charPtr := C.htmlDocDumpToString(document.DocPtr, 0)
	defer C.xmlFreeChars(charPtr)
	return C.GoString(charPtr)
}

func (document *Document) String() string {
	return document.ToXml()
}

func (document *Document) Free() {
	C.xmlFreeDoc(document.DocPtr)
}
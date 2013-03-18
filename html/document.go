package html

/*
#cgo pkg-config: libxml-2.0

#include <libxml/HTMLtree.h>
#include <libxml/HTMLparser.h>
#include "helper.h"
*/
import "C"

import (
	. "gokogiri/util"
	"gokogiri/xml"
	"unsafe"
)

type HtmlDocument struct {
	*xml.XmlDocument
}

var emptyHtmlDocBytes = []byte(EmptyHtmlDoc)
var emptyStringBytes = []byte{0}

//create a document
func NewDocument(p unsafe.Pointer, contentLen int, inEncoding, outEncoding []byte) (doc *HtmlDocument) {
	doc = &HtmlDocument{}
	doc.XmlDocument = xml.NewDocument(p, contentLen, inEncoding, outEncoding)
	doc.SetFragmentParser(&HtmlFragmentParser{})
	return
}

//parse a string to document
func Parse(content, inEncoding, url []byte, options int, outEncoding []byte) (doc *HtmlDocument, err error) {
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

		docPtr = C.htmlParse(contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)

		if docPtr == nil {
			err = ERR_FAILED_TO_PARSE_HTML
		} else {
			doc = NewDocument(unsafe.Pointer(docPtr), contentLen, inEncoding, outEncoding)
		}
	}
	if docPtr == nil {
		doc = CreateEmptyDocument(inEncoding, outEncoding)
	}
	return
}

func CreateEmptyDocument(inEncoding, outEncoding []byte) (doc *HtmlDocument) {
	C.xmlInitParser()
	docPtr := C.htmlNewDoc(nil, nil)
	doc = NewDocument(unsafe.Pointer(docPtr), 0, inEncoding, outEncoding)
	return
}

func (doc *HtmlDocument) MetaEncoding() string {
	metaEncodingXmlCharPtr := C.htmlGetMetaEncoding((*C.xmlDoc)(doc.DocPtr))
	return C.GoString((*C.char)(unsafe.Pointer(metaEncodingXmlCharPtr)))
}

func (doc *HtmlDocument) SetMetaEncoding(encoding string) (err error) {
	var encodingPtr unsafe.Pointer = nil
	if len(encoding) > 0 {
		encodingBytes := AppendCStringTerminator([]byte(encoding))
		encodingPtr = unsafe.Pointer(&encodingBytes[0])
	}
	ret := int(C.htmlSetMetaEncoding((*C.xmlDoc)(doc.DocPtr), (*C.xmlChar)(encodingPtr)))
	if ret == -1 {
		err = ErrSetMetaEncoding
	}
	return
}

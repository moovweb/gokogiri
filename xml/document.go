package xml

/*
#cgo pkg-config: libxml-2.0

#include "helper.h"
*/
import "C"

import (
	"errors"
	. "github.com/moovweb/gokogiri/util"
	"github.com/moovweb/gokogiri/xpath"
	//"runtime"
	"unsafe"
)

type Document interface {
	/* Nokogiri APIs */
	CreateElementNode(string) *ElementNode
	CreateCDataNode(string) *CDataNode
	CreateTextNode(string) *TextNode
	//CreateCommentNode(string) *CommentNode
	ParseFragment([]byte, []byte, int) (*DocumentFragment, error)

	DocPtr() unsafe.Pointer
	DocType() int
	DocRef() Document
	InputEncoding() []byte
	OutputEncoding() []byte
	DocXPathCtx() *xpath.XPath
	AddUnlinkedNode(unsafe.Pointer)
	RemoveUnlinkedNode(unsafe.Pointer) bool
	Free()
	String() string
	Root() *ElementNode
	BookkeepFragment(*DocumentFragment)
}

//xml parse option
// See: http://xmlsoft.org/html/libxml-parser.html#xmlParserOption
const (
	XML_PARSE_RECOVER    = 1 << 0  // recover on errors
	XML_PARSE_NOENT      = 1 << 1  // substitute entities
	XML_PARSE_DTDLOAD    = 1 << 2  // load the external subset
	XML_PARSE_DTDATTR    = 1 << 3  // default DTD attributes
	XML_PARSE_DTDVALID   = 1 << 4  // validate with the DTD
	XML_PARSE_NOERROR    = 1 << 5  // suppress error reports
	XML_PARSE_NOWARNING  = 1 << 6  // suppress warning reports
	XML_PARSE_PEDANTIC   = 1 << 7  // pedantic error reporting
	XML_PARSE_NOBLANKS   = 1 << 8  // remove blank nodes
	XML_PARSE_SAX1       = 1 << 9  // use the SAX1 interface internally
	XML_PARSE_XINCLUDE   = 1 << 10 // Implement XInclude substitition
	XML_PARSE_NONET      = 1 << 11 // Forbid network access
	XML_PARSE_NODICT     = 1 << 12 // Do not reuse the context dictionnary
	XML_PARSE_NSCLEAN    = 1 << 13 // remove redundant namespaces declarations
	XML_PARSE_NOCDATA    = 1 << 14 // merge CDATA as text nodes
	XML_PARSE_NOXINCNODE = 1 << 15 // do not generate XINCLUDE START/END nodes
	XML_PARSE_COMPACT    = 1 << 16 // compact small text nodes; no modification of the tree allowed afterwards (will possibly crash if you try to modify the tree)
	XML_PARSE_OLD10      = 1 << 17 // parse using XML-1.0 before update 5
	XML_PARSE_NOBASEFIX  = 1 << 18 // do not fixup XINCLUDE xml:base uris
	XML_PARSE_HUGE       = 1 << 19 // relax any hardcoded limit from the parser
	XML_PARSE_OLDSAX     = 1 << 20 // parse using SAX2 interface before 2.7.0
	XML_PARSE_IGNORE_ENC = 1 << 21 // ignore internal document encoding hint
	XML_PARSE_BIG_LINES  = 1 << 22 // Store big lines numbers in text PSVI field
)

//default parsing option: relax parsing
var DefaultParseOption = XML_PARSE_RECOVER |
	XML_PARSE_NONET |
	XML_PARSE_NOERROR |
	XML_PARSE_NOWARNING

//libxml2 use "utf-8" by default, and so do we
const DefaultEncoding = "utf-8"

var ERR_FAILED_TO_PARSE_XML = errors.New("failed to parse xml input")

type XmlDocument struct {
	Ptr *C.xmlDoc
	Me  Document
	Node
	InEncoding    []byte
	OutEncoding   []byte
	UnlinkedNodes map[*C.xmlNode]bool
	XPathCtx      *xpath.XPath
	Type          int
	InputLen      int

	fragments []*DocumentFragment //save the pointers to free them when the doc is freed
}

//default encoding in byte slice
var DefaultEncodingBytes = []byte(DefaultEncoding)

const initialFragments = 2

//create a document
func NewDocument(p unsafe.Pointer, contentLen int, inEncoding, outEncoding []byte) (doc *XmlDocument) {
	inEncoding = AppendCStringTerminator(inEncoding)
	outEncoding = AppendCStringTerminator(outEncoding)

	xmlNode := &XmlNode{Ptr: (*C.xmlNode)(p)}
	docPtr := (*C.xmlDoc)(p)
	doc = &XmlDocument{Ptr: docPtr, Node: xmlNode, InEncoding: inEncoding, OutEncoding: outEncoding, InputLen: contentLen}
	doc.UnlinkedNodes = make(map[*C.xmlNode]bool)
	doc.XPathCtx = xpath.NewXPath(p)
	doc.Type = xmlNode.NodeType()
	doc.fragments = make([]*DocumentFragment, 0, initialFragments)
	doc.Me = doc
	xmlNode.Document = doc
	//runtime.SetFinalizer(doc, (*XmlDocument).Free)
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

func (document *XmlDocument) DocPtr() (ptr unsafe.Pointer) {
	ptr = unsafe.Pointer(document.Ptr)
	return
}

func (document *XmlDocument) DocType() (t int) {
	t = document.Type
	return
}

func (document *XmlDocument) DocRef() (d Document) {
	d = document.Me
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
	p := (*C.xmlNode)(nodePtr)
	document.UnlinkedNodes[p] = true
}

func (document *XmlDocument) RemoveUnlinkedNode(nodePtr unsafe.Pointer) bool {
	p := (*C.xmlNode)(nodePtr)
	if document.UnlinkedNodes[p] {
		delete(document.UnlinkedNodes, p)
		return true
	}
	return false
}

func (document *XmlDocument) BookkeepFragment(fragment *DocumentFragment) {
	document.fragments = append(document.fragments, fragment)
}

func (document *XmlDocument) Root() (element *ElementNode) {
	nodePtr := C.xmlDocGetRootElement(document.Ptr)
	if nodePtr != nil {
		element = NewNode(unsafe.Pointer(nodePtr), document).(*ElementNode)
	}
	return
}

func (document *XmlDocument) CreateElementNode(tag string) (element *ElementNode) {
	tagBytes := GetCString([]byte(tag))
	tagPtr := unsafe.Pointer(&tagBytes[0])
	newNodePtr := C.xmlNewNode(nil, (*C.xmlChar)(tagPtr))
	newNode := NewNode(unsafe.Pointer(newNodePtr), document)
	element = newNode.(*ElementNode)
	return
}

func (document *XmlDocument) CreateTextNode(data string) (text *TextNode) {
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewText((*C.xmlChar)(dataPtr))
	if nodePtr != nil {
		nodePtr.doc = (*_Ctype_struct__xmlDoc)(document.DocPtr())
		text = NewNode(unsafe.Pointer(nodePtr), document).(*TextNode)
	}
	return
}

func (document *XmlDocument) CreateCDataNode(data string) (cdata *CDataNode) {
	dataLen := len(data)
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewCDataBlock(document.Ptr, (*C.xmlChar)(dataPtr), C.int(dataLen))
	if nodePtr != nil {
		cdata = NewNode(unsafe.Pointer(nodePtr), document).(*CDataNode)
	}
	return
}

/*
func (document *XmlDocument) CreateCommentNode(data string) (cdata *CommentNode) {
	dataLen := len(data)
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewCDataBlock(document.Ptr, (*C.xmlChar)(dataPtr), C.int(dataLen))
	if nodePtr != nil {
		cdata = NewNode(unsafe.Pointer(nodePtr), document).(*CDataNode)
	}
	return
}
*/

func (document *XmlDocument) ParseFragment(input, url []byte, options int) (fragment *DocumentFragment, err error) {
	root := document.Root()
	if root == nil {
		fragment, err = parsefragment(document, nil, input, url, options)
	} else {
		fragment, err = parsefragment(document, root.XmlNode, input, url, options)
	}
	return
}

func (document *XmlDocument) Free() {
	//must clear the fragments first
	//because the nodes are put in the unlinked list
	if document.fragments != nil {
		for _, fragment := range document.fragments {
			fragment.Remove()
		}
	}
	document.fragments = nil
	var p *C.xmlNode
	if document.UnlinkedNodes != nil {
		for p, _ = range document.UnlinkedNodes {
			C.xmlFreeNode(p)
		}
	}
	document.UnlinkedNodes = nil
	if document.XPathCtx != nil {
		document.XPathCtx.Free()
		document.XPathCtx = nil
	}
	if document.Ptr != nil {
		C.xmlFreeDoc(document.Ptr)
		document.Ptr = nil
	}
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

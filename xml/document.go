package xml

/*
#cgo pkg-config: libxml-2.0

#include "helper.h"
*/
import "C"

import (
	"errors"
	"github.com/moovweb/gokogiri/help"
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
	CreateCommentNode(string) *CommentNode
	CreatePINode(string, string) *ProcessingInstructionNode
	ParseFragment([]byte, []byte, ParseOption) (*DocumentFragment, error)

	DocPtr() unsafe.Pointer
	DocType() NodeType
	DocRef() Document
	InputEncoding() []byte
	OutputEncoding() []byte
	DocXPathCtx() *xpath.XPath
	AddUnlinkedNode(unsafe.Pointer)
	RemoveUnlinkedNode(unsafe.Pointer) bool
	Free()
	String() string
	Root() *ElementNode
	NodeById(string) *ElementNode
	BookkeepFragment(*DocumentFragment)

	RecursivelyRemoveNamespaces() error
}

type ParseOption int

const (
	XML_PARSE_RECOVER   ParseOption = 1 << 0  //relaxed parsing
	XML_PARSE_NOERROR   ParseOption = 1 << 5  //suppress error reports
	XML_PARSE_NOWARNING ParseOption = 1 << 6  //suppress warning reports
	XML_PARSE_NONET     ParseOption = 1 << 11 //forbid network access
)

//default parsing option: relax parsing
var DefaultParseOption ParseOption = XML_PARSE_RECOVER |
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
	Type          NodeType
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

func Parse(content, inEncoding, url []byte, options ParseOption, outEncoding []byte) (doc *XmlDocument, err error) {
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
	help.LibxmlInitParser()
	docPtr := C.newEmptyXmlDoc()
	doc = NewDocument(unsafe.Pointer(docPtr), 0, inEncoding, outEncoding)
	return
}

func (document *XmlDocument) DocPtr() (ptr unsafe.Pointer) {
	ptr = unsafe.Pointer(document.Ptr)
	return
}

func (document *XmlDocument) DocType() (t NodeType) {
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

// Get an element node by the value of its ID attribute. By convention this attribute
// is named id, but the actual name of the attribute is set by the document's DTD or schema.

// The value for an ID attribute is guaranteed to be unique within a valid document.
func (document *XmlDocument) NodeById(id string) (element *ElementNode) {
	dataBytes := GetCString([]byte(id))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlGetID(document.Ptr, (*C.xmlChar)(dataPtr))
	if nodePtr != nil {
		idattr := NewNode(unsafe.Pointer(nodePtr), document).(*AttributeNode)
		element = idattr.Parent().(*ElementNode)
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

func (document *XmlDocument) CreateCommentNode(data string) (comment *CommentNode) {
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewComment((*C.xmlChar)(dataPtr))
	if nodePtr != nil {
		comment = NewNode(unsafe.Pointer(nodePtr), document).(*CommentNode)
	}
	return
}

func (document *XmlDocument) CreatePINode(name, data string) (pi *ProcessingInstructionNode) {
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewDocPI(document.Ptr, (*C.xmlChar)(namePtr), (*C.xmlChar)(dataPtr))
	if nodePtr != nil {
		pi = NewNode(unsafe.Pointer(nodePtr), document).(*ProcessingInstructionNode)
	}
	return
}

func (document *XmlDocument) ParseFragment(input, url []byte, options ParseOption) (fragment *DocumentFragment, err error) {
	root := document.Root()
	if root == nil {
		fragment, err = parsefragment(document, nil, input, url, options)
	} else {
		fragment, err = parsefragment(document, root.XmlNode, input, url, options)
	}
	return
}

func (document *XmlDocument) Free() {
	//must free the xpath context before freeing the fragments or unlinked nodes
	//otherwise, it causes memory leaks and crashes when dealing with very large documents (a few MB)
	if document.XPathCtx != nil {
		document.XPathCtx.Free()
		document.XPathCtx = nil
	}
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

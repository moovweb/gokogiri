package xml

//#include "chelper.h"
//#include <string.h>
import "C"

import (
	"errors"
	"unsafe"
)

var (
	ERR_UNDEFINED_ADD_CHILD_PARAM 				= errors.New("unexpected parameter type in AddChild")
	ERR_UNDEFINED_SET_CONTENT_PARAM 			= errors.New("unexpected parameter type in SetContent")
	ERR_CANNOT_MAKE_DUCMENT_AS_CHILD 			= errors.New("cannot add a document node as a child")
	ERR_CANNOT_COPY_TEXT_NODE_WHEN_ADD_CHILD 	= errors.New("cannot copy a text node when adding it")
)

//xmlNode types
const (
	XML_ELEMENT_NODE       = 1
	XML_ATTRIBUTE_NODE     = 2
	XML_TEXT_NODE          = 3
	XML_CDATA_SECTION_NODE = 4
	XML_ENTITY_REF_NODE    = 5
	XML_ENTITY_NODE        = 6
	XML_PI_NODE            = 7
	XML_COMMENT_NODE       = 8
	XML_DOCUMENT_NODE      = 9
	XML_DOCUMENT_TYPE_NODE = 10
	XML_DOCUMENT_FRAG_NODE = 11
	XML_NOTATION_NODE      = 12
	XML_HTML_DOCUMENT_NODE = 13
	XML_DTD_NODE           = 14
	XML_ELEMENT_DECL       = 15
	XML_ATTRIBUTE_DECL     = 16
	XML_ENTITY_DECL        = 17
	XML_NAMESPACE_DECL     = 18
	XML_XINCLUDE_START     = 19
	XML_XINCLUDE_END       = 20
	XML_DOCB_DOCUMENT_NODE = 21
)

type Node interface {
	GetNodePointer() *C.xmlNode
	ResetNodePointer()
	GetDocument() *Document
	
	//
	GetNodeType() int
	GetNextSibling() Node
	GetPreviousSibling() Node
	
	GetFirstChild() Node
	GetLastChild() Node
	//Attributes() map[string]*AttributeNode
	
	//
	AddChild(interface{}) error
	AddPreviousSibling(interface{}) error
	AddNextSibling(interface{}) error
	//InsertBefore(interface{}) error
	//InsertAfter(interface{}) error
	//SetInnerHtml(interface{}) error
	//Replace(interface{}) error
	//Swap(interface{}) error
	//
	////
	SetContent(interface{}) error
	//SetParent(Node)
	//IsComment() bool
	//IsCData() bool
	//IsXml() bool
	//IsHtml() bool
	//IsText() bool
	//IsElement() bool
	//IsFragment() bool
	//
	
	//
	Free()
	////
	ToXml() []byte
	ToHtml() []byte
	String() string
}

//run out of memory
var ErrTooLarge = errors.New("Output buffer too large")

//pre-allocate a buffer for serializing the document
const initialOutputBufferSize = 100*1024 //100K

type XmlNode struct {
	NodePtr *C.xmlNode
	*Document
	
	outputBuffer []byte
	outputOffset int
}

func NewNode(nodePtr *C.xmlNode, document *Document) (node Node) {
	if nodePtr == nil {
		return nil
	}
	
	xmlNode := &XmlNode{NodePtr: nodePtr, Document: document}
	nodeType := C.getNodeType(nodePtr)
	
	switch nodeType {
	default:
		node = xmlNode
	case XML_ATTRIBUTE_NODE:
		node = &AttributeNode{XmlNode: xmlNode}
	case XML_ELEMENT_NODE:
		node = &ElementNode{XmlNode: xmlNode}
	case XML_CDATA_SECTION_NODE:
		node = &CDataNode{XmlNode: xmlNode}
	case XML_TEXT_NODE:
		node = &TextNode{XmlNode: xmlNode}
	}
	return
}

//
func (xmlNode *XmlNode) AddChild(tag interface{}) (err error) {
	switch t := tag.(type) {
	default:
		err = ERR_UNDEFINED_ADD_CHILD_PARAM
	case *XmlNode:
		err = xmlNode.addChild(t)
	case *DocumentFragment:
		for _, child := range(t.Children) {
			if err = xmlNode.addChild(child); err != nil {
				break
			}
		}
	case string:
		f, err := ParseFragment(xmlNode.Document, []byte(t), xmlNode.Document.Encoding, DefaultParseOption)
		if err == nil {
			xmlNode.AddChild(f)
		}
	case []byte:
		f, err := ParseFragment(xmlNode.Document, t, xmlNode.Document.Encoding, DefaultParseOption)
		if err == nil {
			xmlNode.AddChild(f)
		}
	}
	return
}

func (xmlNode *XmlNode) AddPreviousSibling(tag interface{}) (err error) {
	switch t := tag.(type) {
	default:
		err = ERR_UNDEFINED_ADD_CHILD_PARAM
	case *XmlNode:
		err = xmlNode.addPreviousSibling(t)
	case *DocumentFragment:
		for _, child := range(t.Children) {
			if err = xmlNode.addPreviousSibling(child); err != nil {
				break
			}
		}
	case string:
		f, err := ParseFragment(xmlNode.Document, []byte(t), xmlNode.Document.Encoding, DefaultParseOption)
		if err == nil {
			xmlNode.AddPreviousSibling(f)
		}
	case []byte:
		f, err := ParseFragment(xmlNode.Document, t, xmlNode.Document.Encoding, DefaultParseOption)
		if err == nil {
			xmlNode.AddPreviousSibling(f)
		}
	}
	return
}

func (xmlNode *XmlNode) AddNextSibling(tag interface{}) (err error) {
	switch t := tag.(type) {
	default:
		err = ERR_UNDEFINED_ADD_CHILD_PARAM
	case *XmlNode:
		err = xmlNode.addNextSibling(t)
	case *DocumentFragment:
		for _, child := range(t.Children) {
			if err = xmlNode.addNextSibling(child); err != nil {
				break
			}
		}
	case string:
		f, err := ParseFragment(xmlNode.Document, []byte(t), xmlNode.Document.Encoding, DefaultParseOption)
		if err == nil {
			xmlNode.AddNextSibling(f)
		}
	case []byte:
		f, err := ParseFragment(xmlNode.Document, t, xmlNode.Document.Encoding, DefaultParseOption)
		if err == nil {
			xmlNode.AddNextSibling(f)
		}
	}
	return
}

func (xmlNode *XmlNode) ResetNodePointer() {
	xmlNode.NodePtr = nil
	return
}

func (xmlNode *XmlNode) GetDocument() (document *Document) {
	document = xmlNode.Document
	return
}

func (xmlNode *XmlNode) GetNodePointer() (p *C.xmlNode) {
	p = xmlNode.NodePtr
	return
}

func (xmlNode *XmlNode) GetNodeType() (nodeType int) {
	nodeType = int(C.getNodeType(xmlNode.NodePtr))
	return
}

func (xmlNode *XmlNode) GetNextSibling() Node {
	siblingPtr := (*C.xmlNode)(xmlNode.NodePtr.next);
	return NewNode(siblingPtr, xmlNode.Document)
}

func (xmlNode *XmlNode) GetPreviousSibling() Node {
	siblingPtr := (*C.xmlNode)(xmlNode.NodePtr.prev);
	return NewNode(siblingPtr, xmlNode.Document)
}

func (node *XmlNode) GetFirstChild() Node {
	return NewNode((*C.xmlNode)(node.NodePtr.children), node.Document)
}

func (node *XmlNode) GetLastChild() Node {
	return NewNode((*C.xmlNode)(node.NodePtr.last), node.Document)
}

func (xmlNode *XmlNode) SetContent(content interface{}) (err error) {
	switch data := content.(type) {
	default:
		err = ERR_UNDEFINED_SET_CONTENT_PARAM
	case string:
		err = xmlNode.SetContent([]byte(data))
	case []byte:
		if len(data) > 0 {
			contentPtr := unsafe.Pointer(&data[0])
			C.xmlSetContent(unsafe.Pointer(xmlNode), contentPtr)
		}
	}
	return
}

/*
//func (xmlNode *XmlNode) Attributes() map[string]*AttributeNode

func (xmlNode *XmlNode) InsertBefore(interface{}) error {
	
}
func (xmlNode *XmlNode) InsertAfter(interface{}) error {
	
}
func (xmlNode *XmlNode) SetInnerHtml(interface{}) error {
	
}
func (xmlNode *XmlNode) Replace(interface{}) error {
	
}
func (xmlNode *XmlNode) Swap(interface{}) error {
	
}
func (xmlNode *XmlNode) SetParent(Node) {
	
}
func (xmlNode *XmlNode) IsComment() bool {
	
}
func (xmlNode *XmlNode) IsCData() bool {
	
}
func (xmlNode *XmlNode) IsXml() bool {
	
}
func (xmlNode *XmlNode) IsHtml() bool {
	
}
func (xmlNode *XmlNode) IsText() bool {
	
}
func (xmlNode *XmlNode) IsElement() bool {
	
}
func (xmlNode *XmlNode) IsFragment() bool {
	
}
*/

func (xmlNode *XmlNode) ToXml() []byte {
	xmlNode.outputOffset = 0
	if len(xmlNode.outputBuffer) == 0 {
		xmlNode.outputBuffer = make([]byte, initialOutputBufferSize)
	}
	objPtr := unsafe.Pointer(xmlNode)
	nodePtr      := unsafe.Pointer(xmlNode.NodePtr)
	encodingPtr := unsafe.Pointer(&(xmlNode.Document.Encoding[0]))
	C.xmlSaveNode(objPtr, nodePtr, encodingPtr, XML_SAVE_AS_XML)
	return xmlNode.outputBuffer[:xmlNode.outputOffset]
}

func (xmlNode *XmlNode) ToHtml() []byte {
	return nil
}

func (xmlNode *XmlNode) String() string {
	b := xmlNode.ToXml()
	if b == nil {
		return ""
	}
	return string(b)
}

func (xmlNode *XmlNode) Free() {
	if xmlNode.NodePtr != nil {
		C.xmlFreeNode(xmlNode.NodePtr)
		xmlNode.NodePtr = nil
	}
}

func (xmlNode *XmlNode) addChild(node Node) (err error) {
	nodeType := node.GetNodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.GetNodePointer()
	C.xmlUnlinkNode(nodePtr)
	
	childPtr := C.xmlAddChild(xmlNode.NodePtr, nodePtr)
	if nodeType == XML_TEXT_NODE && childPtr != nodePtr {
		//check the retured pointer
		//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
		//then we should invalid this node, because we do not want to have a dangling pointer
		node.ResetNodePointer()
	}
	return
}

func (xmlNode *XmlNode) addPreviousSibling(node Node) (err error) {
	nodeType := node.GetNodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.GetNodePointer()
	C.xmlUnlinkNode(nodePtr)
	
	childPtr := C.xmlAddPrevSibling(xmlNode.NodePtr, nodePtr)
	if nodeType == XML_TEXT_NODE && childPtr != nodePtr {
		//check the retured pointer
		//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
		//then we should invalid this node, because we do not want to have a dangling pointer
		node.ResetNodePointer()
	}
	return
}

func (xmlNode *XmlNode) addNextSibling(node Node) (err error) {
	nodeType := node.GetNodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.GetNodePointer()
	C.xmlUnlinkNode(nodePtr)
	
	childPtr := C.xmlAddNextSibling(xmlNode.NodePtr, nodePtr)
	if nodeType == XML_TEXT_NODE && childPtr != nodePtr {
		//check the retured pointer
		//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
		//then we should invalid this node, because we do not want to have a dangling pointer
		node.ResetNodePointer()
	}
	return
}


//export xmlNodeWriteCallback
func xmlNodeWriteCallback(obj unsafe.Pointer, data unsafe.Pointer, data_len C.int) {
	node := (*XmlNode)(obj)
	dataLen := int(data_len)

	if node.outputOffset + dataLen > cap(node.outputBuffer) {
		node.outputBuffer = grow(node.outputBuffer, dataLen)
	}
	if dataLen > 0 {
		destBufPtr := unsafe.Pointer(&(node.outputBuffer[node.outputOffset]))
		C.memcpy(destBufPtr, data, C.size_t(data_len))
		node.outputOffset += dataLen
	}
}

func grow(buffer []byte, n int) (newBuffer []byte) {
	newBuffer = makeSlice(2*cap(buffer) + n)
    copy(newBuffer, buffer)
	return
}

func makeSlice(n int) []byte {
    // If the make fails, give a known error.
    defer func() {
        if recover() != nil {
            panic(ErrTooLarge)
        }
    }()
    return make([]byte, n)
}

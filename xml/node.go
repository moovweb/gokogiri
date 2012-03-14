package xml

//#include "helper.h"
//#include <string.h>
import "C"

import (
	"os"
	"unsafe"
	"gokogiri/xpath"
)

var (
	ERR_UNDEFINED_COERCE_PARAM               = os.NewError("unexpected parameter type in coerce")
	ERR_UNDEFINED_SET_CONTENT_PARAM          = os.NewError("unexpected parameter type in SetContent")
	ERR_UNDEFINED_SEARCH_PARAM               = os.NewError("unexpected parameter type in Search")
	ERR_CANNOT_MAKE_DUCMENT_AS_CHILD         = os.NewError("cannot add a document node as a child")
	ERR_CANNOT_COPY_TEXT_NODE_WHEN_ADD_CHILD = os.NewError("cannot copy a text node when adding it")
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

const (
	XML_SAVE_FORMAT   = 1   // format save output
	XML_SAVE_NO_DECL  = 2   //drop the xml declaration
	XML_SAVE_NO_EMPTY = 4   //no empty tags
	XML_SAVE_NO_XHTML = 8   //disable XHTML1 specific rules
	XML_SAVE_XHTML    = 16  //force XHTML1 specific rules
	XML_SAVE_AS_XML   = 32  //force XML serialization on HTML doc
	XML_SAVE_AS_HTML  = 64  //force HTML serialization on XML doc
	XML_SAVE_WSNONSIG = 128 //format with non-significant whitespace
)

type Node interface {
	NodePtr() unsafe.Pointer
	ResetNodePtr()
	MyDocument() Document

	IsValid() bool

	//
	NodeType() int
	NextSibling() Node
	PreviousSibling() Node

	FirstChild() Node
	LastChild() Node
	Attributes() map[string]*AttributeNode

	//
	AddChild(interface{}) os.Error
	AddPreviousSibling(interface{}) os.Error
	AddNextSibling(interface{}) os.Error
	InsertBefore(interface{}) os.Error
	InsertAfter(interface{}) os.Error
	SetInnerHtml(interface{}) os.Error
	SetChildren(interface{}) os.Error
	Replace(interface{}) os.Error
	//Swap(interface{}) os.Error
	//
	////
	SetContent(interface{}) os.Error

	//
	Name() string
	SetName(string)

	//
	Attribute(string) string
	SetAttribute(string, string) *AttributeNode

	//
	Path() string

	//
	Duplicate(int) Node

	Search(interface{}) ([]Node, os.Error)

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
	Unlink()
	Remove()
	ResetChildren()
	//Free()
	////
	ToXml([]byte) []byte
	ToHtml([]byte) []byte
	String() string
	Content() string
}

//run out of memory
var ErrTooLarge = os.NewError("Output buffer too large")

//pre-allocate a buffer for serializing the document
const initialOutputBufferSize = 100 * 1024 //100K

type XmlNode struct {
	Ptr *C.xmlNode
	Document

	outputBuffer []byte
	outputOffset int

	valid bool
}

func NewNode(nodePtr unsafe.Pointer, document Document) (node Node) {
	if nodePtr == nil {
		return nil
	}

	xmlNode := &XmlNode{
		Ptr:      (*C.xmlNode)(nodePtr),
		Document: document,
		valid:    true,
	}
	nodeType := C.getNodeType((*C.xmlNode)(nodePtr))

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

func (xmlNode *XmlNode) coerce(data interface{}) (nodes []Node, err os.Error) {
	switch t := data.(type) {
	default:
		err = ERR_UNDEFINED_COERCE_PARAM
	case []Node:
		nodes = t
	case *DocumentFragment:
		nodes = t.Children.Nodes
	case string:
		f, err := xmlNode.Document.ParseFragment([]byte(t), nil, DefaultParseOption)
		if err == nil {
			nodes = f.Children.Nodes
		}
	case []byte:
		f, err := xmlNode.Document.ParseFragment(t, nil, DefaultParseOption)
		if err == nil {
			nodes = f.Children.Nodes
		}
	}
	return
}

//
func (xmlNode *XmlNode) AddChild(data interface{}) (err os.Error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addChild(node); err != nil {
					break
				}
			}
		}
	case *XmlNode:
		err = xmlNode.addChild(t)
	}
	return
}

func (xmlNode *XmlNode) AddPreviousSibling(data interface{}) (err os.Error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addPreviousSibling(node); err != nil {
					break
				}
			}
		}
	case *XmlNode:
		err = xmlNode.addChild(t)
	}
	return
}

func (xmlNode *XmlNode) AddNextSibling(data interface{}) (err os.Error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addNextSibling(node); err != nil {
					break
				}
			}
		}
	case *XmlNode:
		err = xmlNode.addChild(t)
	}
	return
}

func (xmlNode *XmlNode) ResetNodePtr() {
	xmlNode.Ptr = nil
	return
}

func (xmlNode *XmlNode) IsValid() bool {
	return xmlNode.valid
}

func (xmlNode *XmlNode) MyDocument() (document Document) {
	document = xmlNode.Document
	return
}

func (xmlNode *XmlNode) NodePtr() (p unsafe.Pointer) {
	p = unsafe.Pointer(xmlNode.Ptr)
	return
}

func (xmlNode *XmlNode) NodeType() (nodeType int) {
	nodeType = int(C.getNodeType(xmlNode.Ptr))
	return
}

func (xmlNode *XmlNode) Path() (path string) {
	pathPtr := C.xmlGetNodePath(xmlNode.Ptr)
	if pathPtr != nil {
		p := (*C.char)(unsafe.Pointer(pathPtr))
		defer C.xmlFreeChars(p)
		path = C.GoString(p)
	}
	return
}

func (xmlNode *XmlNode) NextSibling() Node {
	siblingPtr := (*C.xmlNode)(xmlNode.Ptr.next)
	return NewNode(unsafe.Pointer(siblingPtr), xmlNode.Document)
}

func (xmlNode *XmlNode) PreviousSibling() Node {
	siblingPtr := (*C.xmlNode)(xmlNode.Ptr.prev)
	return NewNode(unsafe.Pointer(siblingPtr), xmlNode.Document)
}

func (node *XmlNode) FirstChild() Node {
	return NewNode(unsafe.Pointer(node.Ptr.children), node.Document)
}

func (node *XmlNode) LastChild() Node {
	return NewNode(unsafe.Pointer(node.Ptr.last), node.Document)
}

func (xmlNode *XmlNode) ResetChildren() {
	var p unsafe.Pointer
	for childPtr := xmlNode.Ptr.children; childPtr != nil; {
		nextPtr := childPtr.next
		p = unsafe.Pointer(childPtr)
		C.xmlUnlinkNode((*C.xmlNode)(p))
		xmlNode.Document.AddUnlinkedNode(p)
		childPtr = nextPtr
	}
}

func (xmlNode *XmlNode) SetContent(content interface{}) (err os.Error) {
	switch data := content.(type) {
	default:
		err = ERR_UNDEFINED_SET_CONTENT_PARAM
	case string:
		err = xmlNode.SetContent([]byte(data))
	case []byte:
		if len(data) > 0 {
			contentPtr := unsafe.Pointer(&data[0])
			C.xmlSetContent(unsafe.Pointer(xmlNode.Ptr), contentPtr)
		}
	}
	return
}

func (xmlNode *XmlNode) InsertBefore(data interface{}) os.Error {
	return xmlNode.AddPreviousSibling(data)
}

func (xmlNode *XmlNode) InsertAfter(data interface{}) os.Error {
	return xmlNode.AddNextSibling(data)
}

func (xmlNode *XmlNode) SetChildren(data interface{}) (err os.Error) {
	nodes, err := xmlNode.coerce(data)
	if err != nil {
		return
	}
	xmlNode.ResetChildren()
	err = xmlNode.AddChild(nodes)
	return
}

func (xmlNode *XmlNode) SetInnerHtml(data interface{}) (err os.Error) {
	err = xmlNode.SetChildren(data)
	return
}

func (xmlNode *XmlNode) Replace(data interface{}) (err os.Error) {
	err = xmlNode.AddPreviousSibling(data)
	if err != nil {
		return
	}
	xmlNode.Unlink()
	return
}

func (xmlNode *XmlNode) Attributes() (attributes map[string]*AttributeNode) {
	attributes = make(map[string]*AttributeNode)
	for prop := xmlNode.Ptr.properties; prop != nil; prop = prop.next {
		if prop.name != nil {
			namePtr := unsafe.Pointer(prop.name)
			name := C.GoString((*C.char)(namePtr))
			attrPtr := unsafe.Pointer(prop)
			attributeNode := NewNode(attrPtr, xmlNode.Document)
			if attr, ok := attributeNode.(*AttributeNode); ok {
				attributes[name] = attr
			}
		}
	}
	return
}

func (xmlNode *XmlNode) Attribute(name string) (val string) {
	if len(name) == 0 {
		return
	}
	nameBytes := []byte(name)
	namePtr := unsafe.Pointer(&nameBytes[0])
	valPtr := C.xmlGetProp(xmlNode.Ptr, (*C.xmlChar)(namePtr))
	if valPtr == nil {
		return
	}
	p := unsafe.Pointer(valPtr)
	defer C.xmlFreeChars((*C.char)(p))
	val = C.GoString((*C.char)(p))
	return
}

func (xmlNode *XmlNode) SetAttribute(name, value string) (attribute *AttributeNode) {
	if len(name) == 0 {
		return
	}
	nameBytes := []byte(name)
	namePtr := unsafe.Pointer(&nameBytes[0])
	prop := C.xmlHasProp(xmlNode.Ptr, (*C.xmlChar)(namePtr))
	if prop != nil && prop.children != nil {
		c := prop.children
		for c != nil {
			nextC := c
			child := NewNode(unsafe.Pointer(c), xmlNode.Document)
			if child != nil {
				child.Remove()
			}
			c = nextC
		}
	}

	if len(value) == 0 {
		prop = C.xmlSetProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), nil)
	} else {
		valueBytes := []byte(name)
		valuePtr := unsafe.Pointer(&valueBytes[0])
		prop = C.xmlSetProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), (*C.xmlChar)(valuePtr))
	}
	if prop != nil {
		node := NewNode(unsafe.Pointer(prop), xmlNode.Document)
		if node, ok := node.(*AttributeNode); ok {
			attribute = node
		}
	}
	return
}

func (xmlNode *XmlNode) Search(data interface{}) (result []Node, err os.Error) {
	switch data := data.(type) {
	default:
		err = ERR_UNDEFINED_SEARCH_PARAM
	case string:
		if xpathExpr := xpath.Compile(data); xpathExpr != nil {
			result, err = xmlNode.Search(xpathExpr)
			defer xpathExpr.Free()
		} else {
			err = os.NewError("cannot compile xpath: " + data)
		}
	case []byte:
		result, err = xmlNode.Search(string(data))
	case *xpath.Expression:
		xpathCtx := xmlNode.Document.DocXPathCtx()
		nodePtrs := xpathCtx.Evaluate(unsafe.Pointer(xmlNode.Ptr), data)
		for _, nodePtr := range nodePtrs {
			result = append(result, NewNode(nodePtr, xmlNode.Document))
		}
	}
	return
}

/*
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

func (xmlNode *XmlNode) Name() (name string) {
	if xmlNode.Ptr.name != nil {
		p := unsafe.Pointer(xmlNode.Ptr.name)
		name = C.GoString((*C.char)(p))
	}
	return ""
}

func (xmlNode *XmlNode) SetName(name string) {
	if len(name) > 0 {
		nameBytes := []byte(name)
		namePtr := unsafe.Pointer(&nameBytes[0])
		C.xmlNodeSetName(xmlNode.Ptr, (*C.xmlChar)(namePtr))
	}
}

func (xmlNode *XmlNode) Duplicate(level int) (dup Node) {
	if xmlNode.valid {
		dupPtr := C.xmlDocCopyNode(xmlNode.Ptr, (*C.xmlDoc)(xmlNode.Document.DocPtr()), C.int(level))
		if dupPtr != nil {
			dup = NewNode(unsafe.Pointer(dupPtr), xmlNode.Document)
		}
	}
	return
}

func (xmlNode *XmlNode) to_s(format int, encoding []byte) []byte {
	xmlNode.outputOffset = 0
	if len(xmlNode.outputBuffer) == 0 {
		xmlNode.outputBuffer = make([]byte, initialOutputBufferSize)
	}
	nodePtr := unsafe.Pointer(xmlNode.Ptr)

	var encodingPtr unsafe.Pointer
	if len(encoding) == 0 {
		encoding = xmlNode.Document.OutputEncoding()
	}
	if len(encoding) > 0 {
		encodingPtr = unsafe.Pointer(&(encoding[0]))
	} else {
		encodingPtr = nil
	}
	bufferPtr := unsafe.Pointer(&xmlNode.outputBuffer[0])
	bufferLen := len(xmlNode.outputBuffer)
	format |= XML_SAVE_FORMAT
	ret := int(C.xmlSaveNode(bufferPtr, C.int(bufferLen), nodePtr, encodingPtr, C.int(format)))
	if ret < 0 {
		println("output error!!!")
		return nil
	}
	return xmlNode.outputBuffer[:ret]
}

func (xmlNode *XmlNode) ToXml(encoding []byte) []byte {
	return xmlNode.to_s(XML_SAVE_AS_XML, encoding)
}

func (xmlNode *XmlNode) ToHtml(encoding []byte) []byte {
	return xmlNode.to_s(XML_SAVE_AS_HTML, encoding)
}

func (xmlNode *XmlNode) String() string {
	var b []byte
	if docType := xmlNode.Document.DocType(); docType == XML_HTML_DOCUMENT_NODE {
		b = xmlNode.ToHtml(nil)
	} else {
		b = xmlNode.ToXml(nil)
	}
	if b == nil {
		return ""
	}
	return string(b)
}

func (xmlNode *XmlNode) Content() string {
	contentPtr := C.xmlNodeGetContent(xmlNode.Ptr)
	charPtr := (*C.char)(unsafe.Pointer(contentPtr))
	defer C.xmlFreeChars(charPtr)
	return C.GoString(charPtr)
}

func (xmlNode *XmlNode) Unlink() {
	C.xmlUnlinkNode(xmlNode.Ptr)
	xmlNode.Document.AddUnlinkedNode(unsafe.Pointer(xmlNode.Ptr))
}

func (xmlNode *XmlNode) Remove() {
	if xmlNode.valid {
		xmlNode.Unlink()
		xmlNode.valid = false
	}
}

func (xmlNode *XmlNode) addChild(node Node) (err os.Error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	C.xmlUnlinkNode((*C.xmlNode)(nodePtr))

	childPtr := C.xmlAddChild(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
		//check the retured pointer
		//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
		//then we should invalid this node, because we do not want to have a dangling pointer
		node.Remove()
	}
	return
}

func (xmlNode *XmlNode) addPreviousSibling(node Node) (err os.Error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	C.xmlUnlinkNode((*C.xmlNode)(nodePtr))

	childPtr := C.xmlAddPrevSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
		//check the retured pointer
		//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
		//then we should invalid this node, because we do not want to have a dangling pointer
		node.Remove()
	}
	return
}

func (xmlNode *XmlNode) addNextSibling(node Node) (err os.Error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	C.xmlUnlinkNode((*C.xmlNode)(nodePtr))

	childPtr := C.xmlAddNextSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
		//check the retured pointer
		//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
		//then we should invalid this node, because we do not want to have a dangling pointer
		node.Remove()
	}
	return
}

/*
//export xmlNodeWriteCallback
func xmlNodeWriteCallback(obj unsafe.Pointer, data unsafe.Pointer, data_len C.int) {
	node := (*XmlNode)(obj)
	dataLen := int(data_len)

	if node.outputOffset+dataLen > cap(node.outputBuffer) {
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
*/

package xml

//#include "helper.h"
//#include <string.h>
import "C"

import "time"

import (
	"errors"
	. "gokogiri/util"
	"gokogiri/xpath"
	"unsafe"
)

type Node interface {
	NodePtr() unsafe.Pointer
	ResetNodePtr()
	DocType() int
	InputEncoding() []byte
	OutputEncoding() []byte
	//IsValid() bool

	ParseFragment([]byte, []byte, int) (*DocumentFragment, error)

	//
	NodeType() int
	NextSibling() Node
	PreviousSibling() Node

	Parent() Node
	FirstChild() Node
	LastChild() Node
	CountChildren() int
	Attributes() map[string]*AttributeNode

	//
	Coerce(interface{}) ([]Node, error)

	//
	AddChild(interface{}) error
	AddPreviousSibling(interface{}) error
	AddNextSibling(interface{}) error
	InsertBefore(interface{}) error
	InsertAfter(interface{}) error
	InsertBegin(interface{}) error
	InsertEnd(interface{}) error
	SetInnerHtml(interface{}) error
	SetChildren(interface{}) error
	Replace(interface{}) error
	Wrap(string) error
	//Swap(interface{}) os.Error
	//
	////
	SetContent(interface{}) error

	//
	Name() string
	SetName(string)

	//
	Attr(string) string
	SetAttr(string, string) string
	Attribute(string) *AttributeNode

	//
	Path() string

	//
	Duplicate(int) Node

	Search(interface{}) ([]Node, error)

	//SetParent(Node)
	//IsComment() bool
	//IsCData() bool
	//IsXml() bool
	//IsHtml() bool
	//IsText() bool
	//IsElement() bool
	//IsFragment() bool
	//

	Root() *ElementNode
	CreateCDataNode(string) *CDataNode
	CreateTextNode(string) *TextNode
	CreateElementNode(string) *ElementNode

	//
	Unlink()
	Remove()
	ResetChildren()
	//Free()
	////
	ToXml([]byte, []byte) ([]byte, int)
	ToHtml([]byte, []byte) ([]byte, int)
	ToBuffer([]byte) []byte
	String() string
	Content() string
	InnerHtml() string
}

//pre-allocate a buffer for serializing the document
const initialOutputBufferSize = 10 //100K

type XmlNode struct {
	Ptr *C.xmlNode
	*DocCtx
}


func NewNode(nodePtr unsafe.Pointer, docCtx *DocCtx) (node Node) {
	if nodePtr == nil {
		return nil
	}
	xmlNode := &XmlNode{
		Ptr:      (*C.xmlNode)(nodePtr),
		DocCtx:   docCtx,
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

func (xmlNode *XmlNode) coerce(data interface{}) (nodes []Node, err error) {
	switch t := data.(type) {
	default:
		err = ERR_UNDEFINED_COERCE_PARAM
	case []Node:
		nodes = t
	case *DocumentFragment:
		nodes = t.Children()
	case string:
		f, err := xmlNode.ParseFragment([]byte(t), nil, DefaultParseOption)
		if err == nil {
			nodes = f.Children()
		}
	case []byte:
		f, err := xmlNode.ParseFragment(t, nil, DefaultParseOption)
		if err == nil {
			nodes = f.Children()
		}
	}

	return
}

func (xmlNode *XmlNode) Coerce(data interface{}) (nodes []Node, err error) {
	return xmlNode.coerce(data)
}

//
func (xmlNode *XmlNode) AddChild(data interface{}) (err error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addChild(node); err != nil {
					break
				}
			}
		}
	case *DocumentFragment:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				println("trying to add ", node.NodePtr())
				if err = xmlNode.addChild(node); err != nil {
					break
				}
			}
		}
	case Node:
		err = xmlNode.addChild(t)
	}
	return
}

func (xmlNode *XmlNode) AddPreviousSibling(data interface{}) (err error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addPreviousSibling(node); err != nil {
					break
				}
			}
		}
	case *DocumentFragment:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addPreviousSibling(node); err != nil {
					break
				}
			}
		}
	case Node:
		err = xmlNode.addPreviousSibling(t)
	}
	return
}

func (xmlNode *XmlNode) AddNextSibling(data interface{}) (err error) {
	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for i := len(nodes) - 1; i >= 0; i-- {
				node := nodes[i]
				if err = xmlNode.addNextSibling(node); err != nil {
					break
				}
			}
		}
	case *DocumentFragment:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for i := len(nodes) - 1; i >= 0; i-- {
				node := nodes[i]
				if err = xmlNode.addNextSibling(node); err != nil {
					break
				}
			}
		}
	case Node:
		err = xmlNode.addNextSibling(t)
	}
	return
}

func (xmlNode *XmlNode) ResetNodePtr() {
	xmlNode.Ptr = nil
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
	return NewNode(unsafe.Pointer(siblingPtr), xmlNode.DocCtx)
}

func (xmlNode *XmlNode) PreviousSibling() Node {
	siblingPtr := (*C.xmlNode)(xmlNode.Ptr.prev)
	return NewNode(unsafe.Pointer(siblingPtr), xmlNode.DocCtx)
}

func (xmlNode *XmlNode) CountChildren() int {
	return int(C.xmlLsCountNode(xmlNode.Ptr))
}

func (xmlNode *XmlNode) FirstChild() Node {
	return NewNode(unsafe.Pointer(xmlNode.Ptr.children), xmlNode.DocCtx)
}

func (xmlNode *XmlNode) LastChild() Node {
	return NewNode(unsafe.Pointer(xmlNode.Ptr.last), xmlNode.DocCtx)
}

func (xmlNode *XmlNode) Parent() Node {
	if C.xmlNodePtrCheck(unsafe.Pointer(xmlNode.Ptr.parent)) == C.int(0) {
		return nil
	}
	return NewNode(unsafe.Pointer(xmlNode.Ptr.parent), xmlNode.DocCtx)
}

func (xmlNode *XmlNode) ResetChildren() {
	var p unsafe.Pointer
	for childPtr := xmlNode.Ptr.children; childPtr != nil; {
		nextPtr := childPtr.next
		p = unsafe.Pointer(childPtr)
		C.xmlUnlinkNodeWithCheck((*C.xmlNode)(p))
		xmlNode.AddUnlinkedNode(p)
		childPtr = nextPtr
	}
}

func (xmlNode *XmlNode) SetContent(content interface{}) (err error) {
	switch data := content.(type) {
	default:
		err = ERR_UNDEFINED_SET_CONTENT_PARAM
	case string:
		err = xmlNode.SetContent([]byte(data))
	case []byte:
		contentBytes := GetCString(data)
		contentPtr := unsafe.Pointer(&contentBytes[0])
		C.xmlSetContent(unsafe.Pointer(xmlNode), unsafe.Pointer(xmlNode.Ptr), contentPtr)
	}
	return
}

func (xmlNode *XmlNode) InsertBefore(data interface{}) (err error) {
	err = xmlNode.AddPreviousSibling(data)
	return
}

func (xmlNode *XmlNode) InsertAfter(data interface{}) (err error) {
	err = xmlNode.AddNextSibling(data)
	return
}

func (xmlNode *XmlNode) InsertBegin(data interface{}) (err error) {
	if parent := xmlNode.Parent(); parent != nil {
		if last := parent.LastChild(); last != nil {
			err = last.AddPreviousSibling(data)
		}
	}
	return
}

func (xmlNode *XmlNode) InsertEnd(data interface{}) (err error) {
	if parent := xmlNode.Parent(); parent != nil {
		if first := parent.FirstChild(); first != nil {
			err = first.AddPreviousSibling(data)
		}
	}
	return
}

func (xmlNode *XmlNode) SetChildren(data interface{}) (err error) {
	nodes, err := xmlNode.coerce(data)
	if err != nil {
		return
	}
	xmlNode.ResetChildren()
	err = xmlNode.AddChild(nodes)
	return nil
}

func (xmlNode *XmlNode) SetInnerHtml(data interface{}) (err error) {
	err = xmlNode.SetChildren(data)
	return
}

func (xmlNode *XmlNode) Replace(data interface{}) (err error) {
	err = xmlNode.AddPreviousSibling(data)
	if err != nil {
		return
	}
	xmlNode.Remove()
	return
}

func (xmlNode *XmlNode) Attributes() (attributes map[string]*AttributeNode) {
	attributes = make(map[string]*AttributeNode)
	for prop := xmlNode.Ptr.properties; prop != nil; prop = prop.next {
		if prop.name != nil {
			namePtr := unsafe.Pointer(prop.name)
			name := C.GoString((*C.char)(namePtr))
			attrPtr := unsafe.Pointer(prop)
			attributeNode := NewNode(attrPtr, xmlNode.DocCtx)
			if attr, ok := attributeNode.(*AttributeNode); ok {
				attributes[name] = attr
			}
		}
	}
	return
}

func (xmlNode *XmlNode) Attribute(name string) (attribute *AttributeNode) {
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])
	attrPtr := C.xmlHasNsProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), nil)
	if attrPtr == nil {
		return
	} else {
		node := NewNode(unsafe.Pointer(attrPtr), xmlNode.DocCtx)
		if node, ok := node.(*AttributeNode); ok {
			attribute = node
		}
	}
	return
}

func (xmlNode *XmlNode) Attr(name string) (val string) {
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
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

func (xmlNode *XmlNode) SetAttr(name, value string) (val string) {
	val = value
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])

	valueBytes := GetCString([]byte(value))
	valuePtr := unsafe.Pointer(&valueBytes[0])

	C.xmlSetProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), (*C.xmlChar)(valuePtr))
	return
}

func (xmlNode *XmlNode) Search(data interface{}) (result []Node, err error) {
	switch data := data.(type) {
	default:
		err = ERR_UNDEFINED_SEARCH_PARAM
	case string:
		if xpathExpr := xpath.Compile(data); xpathExpr != nil {
			result, err = xmlNode.Search(xpathExpr)
			defer xpathExpr.Free()
		} else {
			err = errors.New("cannot compile xpath: " + data)
		}
	case []byte:
		result, err = xmlNode.Search(string(data))
	case *xpath.Expression:
		xpathCtx := xmlNode.XPathCtx
		nodePtrs, err := xpathCtx.Evaluate(unsafe.Pointer(xmlNode.Ptr), data)
		if nodePtrs == nil || err != nil {
			return nil, err
		}
		for _, nodePtr := range nodePtrs {
			result = append(result, NewNode(nodePtr, xmlNode.DocCtx))
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
	return
}

func (xmlNode *XmlNode) SetName(name string) {
	if len(name) > 0 {
		nameBytes := GetCString([]byte(name))
		namePtr := unsafe.Pointer(&nameBytes[0])
		C.xmlNodeSetName(xmlNode.Ptr, (*C.xmlChar)(namePtr))
	}
}

func (xmlNode *XmlNode) Duplicate(level int) (dup Node) {
	dupPtr := C.xmlDocCopyNode(xmlNode.Ptr, (*C.xmlDoc)(xmlNode.DocPtr), C.int(level))
	if dupPtr != nil {
		dup = NewNode(unsafe.Pointer(dupPtr), xmlNode.DocCtx)
	}
	return
}

func (xmlNode *XmlNode) ToXml(encoding, outputBuffer []byte) ([]byte, int) {
	return serialize(xmlNode, XML_SAVE_AS_XML, encoding, outputBuffer)
}

func (xmlNode *XmlNode) ToHtml(encoding, outputBuffer []byte) ([]byte, int) {
	return serialize(xmlNode, XML_SAVE_AS_HTML, encoding, outputBuffer)
}

func (xmlNode *XmlNode) ToBuffer(outputBuffer []byte) []byte {
	var b []byte
	var size int
	if docType := xmlNode.DocType(); docType == XML_HTML_DOCUMENT_NODE {
		b, size = xmlNode.ToHtml(nil, outputBuffer)
	} else {
		b, size = xmlNode.ToXml(nil, outputBuffer)
	}
	return b[:size]
}

func (xmlNode *XmlNode) String() string {
	b := xmlNode.ToBuffer(nil)
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

func (xmlNode *XmlNode) InnerHtml() string {
	out := ""

	for child := xmlNode.FirstChild(); child != nil; child = child.NextSibling() {
		out += child.String()
	}
	return out
}

func (xmlNode *XmlNode) Unlink() {
	if int(C.xmlUnlinkNodeWithCheck(xmlNode.Ptr)) != 0 {
		xmlNode.AddUnlinkedNode(unsafe.Pointer(xmlNode.Ptr))
	}
}

func (xmlNode *XmlNode) Remove() {
	if unsafe.Pointer(xmlNode.Ptr) != xmlNode.DocPtr {
		xmlNode.Unlink()
	}
}

func (xmlNode *XmlNode) addChild(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	if xmlNode.NodePtr() == nodePtr {
		return
	}
	ret := xmlNode.isAccestor(nodePtr)
	if ret < 0 {
		return
	} else if ret == 0 {
		if ! xmlNode.RemoveUnlinkedNode(nodePtr) {
			C.xmlUnlinkNodeWithCheck((*C.xmlNode)(nodePtr))
		}
		C.xmlAddChild(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	} else if ret > 0 {
		node.Remove()
	}
	return
}

func (xmlNode *XmlNode) addPreviousSibling(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	if xmlNode.NodePtr() == nodePtr {
		return
	}
	ret := xmlNode.isAccestor(nodePtr)
	if ret < 0 {
		return
	} else if ret == 0 {
		if ! xmlNode.RemoveUnlinkedNode(nodePtr) {
			C.xmlUnlinkNodeWithCheck((*C.xmlNode)(nodePtr))
		}
		C.xmlAddPrevSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	} else if ret > 0 {
		node.Remove()
	}
	return
}

func (xmlNode *XmlNode) addNextSibling(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	if xmlNode.NodePtr() == nodePtr {
		return
	}
	ret := xmlNode.isAccestor(nodePtr)
	if ret < 0 {
		return
	} else if ret == 0 {
		if ! xmlNode.RemoveUnlinkedNode(nodePtr) {
			C.xmlUnlinkNodeWithCheck((*C.xmlNode)(nodePtr))
		}
		C.xmlAddNextSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	} else if ret > 0 {
		node.Remove()
	}
	return
}

func (xmlNode *XmlNode) Wrap(data string) (err error) {
	newNodes, err := xmlNode.coerce(data)
	if err == nil && len(newNodes) > 0 {
		newParent := newNodes[0]
		xmlNode.addNextSibling(newParent)
		newParent.AddChild(xmlNode)
	}
	return
}

//export xmlUnlinkNodeCallback
func xmlUnlinkNodeCallback(nodePtr unsafe.Pointer, gonodePtr unsafe.Pointer) {
	xmlNode := (*XmlNode)(gonodePtr)
	xmlNode.AddUnlinkedNode(nodePtr)
}

func (xmlNode *XmlNode) isAccestor(nodePtr unsafe.Pointer) int {
	parentPtr := xmlNode.Ptr.parent

	if C.xmlNodePtrCheck(unsafe.Pointer(parentPtr)) == C.int(0) {
		return -1
	}
	for ; parentPtr != nil; parentPtr = parentPtr.parent {
		if C.xmlNodePtrCheck(unsafe.Pointer(parentPtr)) == C.int(0) {
			return -1
		}
		p := unsafe.Pointer(parentPtr)
		if p == nodePtr {
			return 1
		}
	}
	return 0
}

package xml

//#include "helper.h"
//#include <string.h>
import "C"

import (
	"errors"
	. "gokogiri/util"
	"gokogiri/xpath"
	"unsafe"

	// the following two packages are imported for profiling stuff
	"fmt"
)

var (
	ERR_UNDEFINED_COERCE_PARAM               = errors.New("unexpected parameter type in coerce")
	ERR_UNDEFINED_SET_CONTENT_PARAM          = errors.New("unexpected parameter type in SetContent")
	ERR_UNDEFINED_SEARCH_PARAM               = errors.New("unexpected parameter type in Search")
	ERR_CANNOT_MAKE_DUCMENT_AS_CHILD         = errors.New("cannot add a document node as a child")
	ERR_CANNOT_COPY_TEXT_NODE_WHEN_ADD_CHILD = errors.New("cannot copy a text node when adding it")
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

//run out of memory
var ErrTooLarge = errors.New("Output buffer too large")

//pre-allocate a buffer for serializing the document
const initialOutputBufferSize = 10 //100K

type XmlNode struct {
	Ptr *C.xmlNode
	Document
	valid bool
}

type WriteBuffer struct {
	Node   *XmlNode
	Buffer []byte
	Offset int
}

func NewNode(nodePtr unsafe.Pointer, document Document) (node Node) {

	// document.StartProfiling("NewNode")

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

	// document.StopProfiling()

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
		f, err := xmlNode.MyDocument().ParseFragment([]byte(t), nil, DefaultParseOption)
		if err == nil {
			nodes = f.Children()
		}
	case []byte:
		f, err := xmlNode.MyDocument().ParseFragment(t, nil, DefaultParseOption)
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

	xmlNode.Document.StartProfiling("AddChild")

	switch t := data.(type) {
	default:
		if nodes, err := xmlNode.coerce(data); err == nil {
			for _, node := range nodes {
				if err = xmlNode.addChild(node); err != nil {
					break
				}
			}
		}
	case Node:
		err = xmlNode.addChild(t)
	}

	xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) AddPreviousSibling(data interface{}) (err error) {

	xmlNode.Document.StartProfiling("AddPreviousSibling")

	switch t := data.(type) {
	default:
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

	xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) AddNextSibling(data interface{}) (err error) {

	xmlNode.Document.StartProfiling("AddNextSibling")

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
	case Node:
		err = xmlNode.addNextSibling(t)
	}

	xmlNode.Document.StopProfiling()

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
	document = xmlNode.Document.DocRef()
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

	xmlNode.Document.StartProfiling("Path")

	pathPtr := C.xmlGetNodePath(xmlNode.Ptr)
	if pathPtr != nil {
		p := (*C.char)(unsafe.Pointer(pathPtr))
		defer C.xmlFreeChars(p)
		path = C.GoString(p)
	}

	xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) NextSibling() (n Node) {
	// xmlNode.Document.StartProfiling("NextSibling")
	siblingPtr := (*C.xmlNode)(xmlNode.Ptr.next)
	n = NewNode(unsafe.Pointer(siblingPtr), xmlNode.Document)
	// xmlNode.Document.StopProfiling()
	return
}

func (xmlNode *XmlNode) PreviousSibling() (n Node) {
	// xmlNode.Document.StartProfiling("PreviousSibling")
	siblingPtr := (*C.xmlNode)(xmlNode.Ptr.prev)
	n = NewNode(unsafe.Pointer(siblingPtr), xmlNode.Document)
	// xmlNode.Document.StopProfiling()
	return
}

func (xmlNode *XmlNode) CountChildren() (i int) {
	// xmlNode.Document.StartProfiling("CountChildren")
	i = int(C.xmlLsCountNode(xmlNode.Ptr))
	// xmlNode.Document.StopProfiling()
	return
}

func (xmlNode *XmlNode) FirstChild() (n Node) {
	// xmlNode.Document.StartProfiling("FirstChild")
	n = NewNode(unsafe.Pointer(xmlNode.Ptr.children), xmlNode.Document)
	// xmlNode.Document.StopProfiling()
	return
}

func (xmlNode *XmlNode) LastChild() (n Node) {
	// xmlNode.Document.StartProfiling("LastChild")
	n = NewNode(unsafe.Pointer(xmlNode.Ptr.last), xmlNode.Document)
	// xmlNode.Document.StopProfiling()
	return
}

func (xmlNode *XmlNode) Parent() (n Node) {
	// xmlNode.Document.StartProfiling("Parent")
	if C.xmlNodePtrCheck(unsafe.Pointer(xmlNode.Ptr.parent)) == C.int(0) {
		return nil
	}
	n = NewNode(unsafe.Pointer(xmlNode.Ptr.parent), xmlNode.Document)
	// xmlNode.Document.StopProfiling()
	return
}

func (xmlNode *XmlNode) ResetChildren() {

	xmlNode.Document.StartProfiling("ResetChildren")

	var p unsafe.Pointer
	for childPtr := xmlNode.Ptr.children; childPtr != nil; {
		nextPtr := childPtr.next
		p = unsafe.Pointer(childPtr)
		C.xmlUnlinkNode((*C.xmlNode)(p))
		xmlNode.Document.AddUnlinkedNode(p)
		childPtr = nextPtr
	}

	xmlNode.Document.StopProfiling()
}

func (xmlNode *XmlNode) SetContent(content interface{}) (err error) {
	switch data := content.(type) {
	default:
		err = ERR_UNDEFINED_SET_CONTENT_PARAM
	case string:
		err = xmlNode.SetContent([]byte(data))
	case []byte:
		xmlNode.Document.StartProfiling("SetContent")
		contentBytes := GetCString(data)
		contentPtr := unsafe.Pointer(&contentBytes[0])
		C.xmlSetContent(unsafe.Pointer(xmlNode.Ptr), contentPtr)
		xmlNode.Document.StopProfiling()
	}
	return
}

func (xmlNode *XmlNode) InsertBefore(data interface{}) (err error) {

	// xmlNode.Document.StartProfiling("InsertBefore")

	err = xmlNode.AddPreviousSibling(data)

	// xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) InsertAfter(data interface{}) (err error) {

	// xmlNode.Document.StartProfiling("InsertAfter")

	err = xmlNode.AddNextSibling(data)

	// xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) InsertBegin(data interface{}) (err error) {

	// xmlNode.Document.StartProfiling("InsertBegin")

	if parent := xmlNode.Parent(); parent != nil {
		if last := parent.LastChild(); last != nil {
			err = last.AddPreviousSibling(data)
		}
	}

	// xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) InsertEnd(data interface{}) (err error) {

	// xmlNode.Document.StartProfiling("InsertEnd")

	if parent := xmlNode.Parent(); parent != nil {
		if first := parent.FirstChild(); first != nil {
			err = first.AddPreviousSibling(data)
		}
	}

	// xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) SetChildren(data interface{}) (err error) {

	xmlNode.Document.StartProfiling("SetChildren")

	nodes, err := xmlNode.coerce(data)
	if err != nil {
		return
	}
	xmlNode.ResetChildren()
	err = xmlNode.AddChild(nodes)

	xmlNode.Document.StopProfiling()

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

	// xmlNode.Document.StartProfiling("Attributes")

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

	// xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) Attribute(name string) (attribute *AttributeNode) {

	// xmlNode.Document.StartProfiling("Attribute")

	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])
	attrPtr := C.xmlHasNsProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), nil)
	if attrPtr == nil {
		return
	} else {
		node := NewNode(unsafe.Pointer(attrPtr), xmlNode.Document)
		if node, ok := node.(*AttributeNode); ok {
			attribute = node
		}
	}

	// xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) Attr(name string) (val string) {

	xmlNode.Document.StartProfiling("Attr")

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

	xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) SetAttr(name, value string) (val string) {

	xmlNode.Document.StartProfiling("SetAttr")

	val = value
	if xmlNode.NodeType() != XML_ELEMENT_NODE {
		return
	}
	nameBytes := GetCString([]byte(name))
	namePtr := unsafe.Pointer(&nameBytes[0])

	valueBytes := GetCString([]byte(value))
	valuePtr := unsafe.Pointer(&valueBytes[0])

	C.xmlSetProp(xmlNode.Ptr, (*C.xmlChar)(namePtr), (*C.xmlChar)(valuePtr))

	xmlNode.Document.StopProfiling()

	return
}

func TimedXPathCompile(data string, doc Document) (expr *xpath.Expression) {
	doc.StartProfiling("xpath.Compile")
	expr = xpath.Compile(data)
	doc.StopProfiling()
	return
}

func (xmlNode *XmlNode) Search(data interface{}) (result []Node, err error) {
	switch data := data.(type) {
	default:
		err = ERR_UNDEFINED_SEARCH_PARAM
	case string:
		// xpathExpr := TimedXPathCompile(data, xmlNode.Document)
		xpathExpr := xpath.Compile(data)
		if xpathExpr != nil {
			result, err = xmlNode.Search(xpathExpr)
			defer xpathExpr.Free()
		} else {
			err = errors.New("cannot compile xpath: " + data)
		}
	case []byte:
		result, err = xmlNode.Search(string(data))
	case *xpath.Expression:
		xmlNode.Document.StartProfiling("Search")
		xpathCtx := xmlNode.Document.DocXPathCtx()
		nodePtrs := xpathCtx.Evaluate(unsafe.Pointer(xmlNode.Ptr), data)
		for _, nodePtr := range nodePtrs {
			result = append(result, NewNode(nodePtr, xmlNode.Document))
		}
		xmlNode.Document.StopProfiling()
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

	xmlNode.Document.StartProfiling("Name")

	if xmlNode.Ptr.name != nil {
		p := unsafe.Pointer(xmlNode.Ptr.name)
		name = C.GoString((*C.char)(p))
	}

	xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) SetName(name string) {

	xmlNode.Document.StartProfiling("SetName")

	if len(name) > 0 {
		nameBytes := GetCString([]byte(name))
		namePtr := unsafe.Pointer(&nameBytes[0])
		C.xmlNodeSetName(xmlNode.Ptr, (*C.xmlChar)(namePtr))
	}

	xmlNode.Document.StopProfiling()
}

func (xmlNode *XmlNode) Duplicate(level int) (dup Node) {

	// xmlNode.Document.StartProfiling("Duplicate")

	if xmlNode.valid {
		dupPtr := C.xmlDocCopyNode(xmlNode.Ptr, (*C.xmlDoc)(xmlNode.Document.DocPtr()), C.int(level))
		if dupPtr != nil {
			dup = NewNode(unsafe.Pointer(dupPtr), xmlNode.Document)
		}
	}

	// xmlNode.Document.StopProfiling()

	return
}

func (xmlNode *XmlNode) serialize(format int, encoding, outputBuffer []byte) ([]byte, int) {

	xmlNode.Document.StartProfiling("serialize")

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

	wbuffer := &WriteBuffer{Node: xmlNode, Buffer: outputBuffer}
	wbufferPtr := unsafe.Pointer(wbuffer)

	format |= XML_SAVE_FORMAT
	ret := int(C.xmlSaveNode(wbufferPtr, nodePtr, encodingPtr, C.int(format)))
	if ret < 0 {
		println("output error!!!")
		return nil, 0
	}

	xmlNode.Document.StopProfiling()

	return wbuffer.Buffer, wbuffer.Offset
}

func (xmlNode *XmlNode) ToXml(encoding, outputBuffer []byte) ([]byte, int) {
	return xmlNode.serialize(XML_SAVE_AS_XML, encoding, outputBuffer)
}

func (xmlNode *XmlNode) ToHtml(encoding, outputBuffer []byte) ([]byte, int) {
	return xmlNode.serialize(XML_SAVE_AS_HTML, encoding, outputBuffer)
}

func (xmlNode *XmlNode) ToBuffer(outputBuffer []byte) []byte {
	var b []byte
	var size int
	if docType := xmlNode.Document.DocType(); docType == XML_HTML_DOCUMENT_NODE {
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

	xmlNode.Document.StartProfiling("Content")

	contentPtr := C.xmlNodeGetContent(xmlNode.Ptr)
	charPtr := (*C.char)(unsafe.Pointer(contentPtr))
	defer C.xmlFreeChars(charPtr)

	xmlNode.Document.StopProfiling()

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

	// xmlNode.Document.StartProfiling("Unlink")

	if int(C.xmlUnlinkNodeWithCheck(xmlNode.Ptr)) != 0 {
		xmlNode.Document.AddUnlinkedNode(unsafe.Pointer(xmlNode.Ptr))
	}

	// xmlNode.Document.StopProfiling()

}

func (xmlNode *XmlNode) Remove() {

	xmlNode.Document.StartProfiling("Remove")

	if xmlNode.valid {
		xmlNode.Unlink()
		xmlNode.valid = false
	}

	xmlNode.Document.StopProfiling()
}

func (xmlNode *XmlNode) addChild(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	parentPtr := xmlNode.Ptr.parent

	if C.xmlNodePtrCheck(unsafe.Pointer(parentPtr)) == C.int(0) {
		return
	}

	isNodeAccestor := false
	for ; parentPtr != nil; parentPtr = parentPtr.parent {
		if C.xmlNodePtrCheck(unsafe.Pointer(parentPtr)) == C.int(0) {
			return
		}
		p := unsafe.Pointer(parentPtr)
		if p == nodePtr {
			isNodeAccestor = true
		}
	}
	if !isNodeAccestor {
		C.xmlUnlinkNode((*C.xmlNode)(nodePtr))
		C.xmlAddChild(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	} else {
		node.Remove()
	}

	/*
		childPtr := C.xmlAddChild(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
		if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
			//check the retured pointer
			//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
			//then we should invalid this node, because we do not want to have a dangling pointer
			node.Remove()
		}
	*/
	return
}

func (xmlNode *XmlNode) addPreviousSibling(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	C.xmlUnlinkNode((*C.xmlNode)(nodePtr))

	C.xmlAddPrevSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	/*
		childPtr := C.xmlAddPrevSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
		if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
			//check the retured pointer
			//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
			//then we should invalid this node, because we do not want to have a dangling pointer
			//xmlNode.Document.AddUnlinkedNode(unsafe.Pointer(nodePtr))
		}
	*/
	return
}

func (xmlNode *XmlNode) addNextSibling(node Node) (err error) {
	nodeType := node.NodeType()
	if nodeType == XML_DOCUMENT_NODE || nodeType == XML_HTML_DOCUMENT_NODE {
		err = ERR_CANNOT_MAKE_DUCMENT_AS_CHILD
		return
	}
	nodePtr := node.NodePtr()
	C.xmlUnlinkNode((*C.xmlNode)(nodePtr))
	C.xmlAddNextSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
	/*
		childPtr := C.xmlAddNextSibling(xmlNode.Ptr, (*C.xmlNode)(nodePtr))
		if nodeType == XML_TEXT_NODE && childPtr != (*C.xmlNode)(nodePtr) {
			//check the retured pointer
			//if it is not the text node just added, it means that the text node is freed because it has merged into other nodes
			//then we should invalid this node, because we do not want to have a dangling pointer
			//node.Remove()
		}
	*/
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

func (xmlNode *XmlNode) ParseFragment(input, url []byte, options int) (fragment *DocumentFragment, err error) {
	fragment, err = parsefragment(xmlNode.Document, xmlNode, input, url, options)
	return
}

//export xmlNodeWriteCallback
func xmlNodeWriteCallback(wbufferObj unsafe.Pointer, data unsafe.Pointer, data_len C.int) {
	wbuffer := (*WriteBuffer)(wbufferObj)
	offset := wbuffer.Offset

	if offset > len(wbuffer.Buffer) {
		panic("fatal error in xmlNodeWriteCallback")
	}

	buffer := wbuffer.Buffer[:offset]
	dataLen := int(data_len)

	if dataLen > 0 {
		if len(buffer)+dataLen > cap(buffer) {
			newBuffer := grow(buffer, dataLen)
			wbuffer.Buffer = newBuffer
		}
		destBufPtr := unsafe.Pointer(&(wbuffer.Buffer[offset]))
		C.memcpy(destBufPtr, data, C.size_t(dataLen))
		wbuffer.Offset += dataLen
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

var (
	SearchCount int64
	SearchTime	int64
)

func init() {
	fmt.Println("Just loaded node.go!")
	// SearchCount = 0
	// SearchTime  = 0
}
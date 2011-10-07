package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h>
*/
import "C"
import "unsafe"

type Node interface {
	ptr() *C.xmlNode
	Ptr() unsafe.Pointer // Used to access the C.Ptr's externally
	Doc() *Doc // reference to doc

	String() string
	Remove() bool

	// Traversal Methods
	Parent() Node // child->parent link
	First()  Node // first child link
	Last()   Node // last child link
	Next()   Node // next sibling link
	Prev()   Node // previous sibling link

	// Informational Methods
	Size() int
	Type() int

	// Node Getters and Setters
	Name() string
	SetName(name string)
	
	//Content() string
	//SetContent(content string)

	Attribute(name string) (*Attribute, bool) // First, the attribute, then if it is new or not
}

type XmlNode struct {
	NodePtr *C.xmlNode
	DocRef  *Doc
}

func NewNode(ptr unsafe.Pointer, doc *Doc) Node {
	cPtr := (*C.xmlNode)(ptr)
	if cPtr == nil {
		return nil
	}
	node_type := xmlNodeType(cPtr)
	xml_node := &XmlNode{NodePtr: cPtr, DocRef: doc}
	if doc == nil {
		doc := &Doc{XmlNode: xml_node}
		// If we are a doc, then we reference ourselves
		doc.XmlNode.DocRef = doc
		return doc
	} else if node_type == C.XML_ELEMENT_NODE {
		return &Element{XmlNode: xml_node}
	} else if node_type == C.XML_ATTRIBUTE_NODE {
		return &Attribute{XmlNode: xml_node}
	}
	return xml_node
}

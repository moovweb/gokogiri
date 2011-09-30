package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/HTMLtree.h>
*/
import "C"
import "unsafe"
//import . "libxml/help"

type Node interface {
	ptr() *C.xmlNode
	Ptr() unsafe.Pointer // Used to access the C.Ptr's externally
	Doc() *Doc // reference to doc

	String() string
	Remove()

	// Standard libxml Node interface
	//Children() []Node;
	First() Node  // first child link
	Last() Node   // last child link
	Parent() Node // child->parent link
	Next() Node   // next sibling link
	Prev() Node   // previous sibling link
	Size() int
	Type() int

	Name() string
	SetName(name string)
	AttributeValue(name string) string
	SetAttributeValue(name string, value string)
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
	} else if node_type == XML_ELEMENT_NODE {
		return &Element{XmlNode: xml_node}
	}
	return xml_node
}

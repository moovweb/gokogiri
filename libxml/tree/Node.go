package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/HTMLtree.h>
*/
import "C"

//import . "libxml/help"

type Node interface {
	Ptr() *C.xmlNode
	AnonPtr() interface{} // Used to access the C.Ptr's externally (and Type Assert them)
	Doc() *Doc // reference to doc

	Dump() string
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

func NewNode(undefined_ptr interface{}, doc *Doc) Node {
	ptr := undefined_ptr.(*C.xmlNode)
	if ptr == nil {
		return nil
	}
	node_type := xmlNodeType(ptr)
	xml_node := &XmlNode{NodePtr: ptr, DocRef: doc}
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

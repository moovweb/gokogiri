package libxml

/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/HTMLtree.h>
*/
import "C"

type Node interface {
	Ptr() *C.xmlNode
	Doc() *XmlDoc // reference to doc

	Dump() string
	Search(xpath string) *NodeSet
	Remove()

	// Standard libxml Node interface
	//Children() []Node;
	First() Node  // first child link
	Last() Node   // last child link
	Parent() Node // child->parent link
	Next() Node   // next sibling link
	Prev() Node   // previous sibling link
	Type() int

	Name() string
	SetName(name string)
	AttributeValue(name string) string
	SetAttributeValue(name string, value string)
}

func buildNode(ptr *C.xmlNode, doc *XmlDoc) Node {
	if ptr == nil {
		return nil
	}
	node_type := xmlNodeType(ptr)
	xml_node := &XmlNode{NodePtr: ptr, DocRef: doc}
	if doc == nil {
		doc := &XmlDoc{XmlNode: xml_node}
		// If we are a doc, then we reference ourselves
		doc.XmlNode.DocRef = doc
		return doc
	} else if node_type == XML_ELEMENT_NODE {
		return &XmlElement{XmlNode: xml_node}
	}
	return xml_node
}

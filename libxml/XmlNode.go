package libxml
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 

xmlNode * GoNodeNext(xmlNode *node) { return node->next; } 
xmlNode * GoNodePrev(xmlNode *node) { return node->prev; } 
xmlNode * GoNodeChildren(xmlNode *node) { return node->children; }
xmlNode * GoNodeLast(xmlNode *node) { return node->last; } 
xmlNode * GoNodeParent(xmlNode *node) { return node->parent; } 
int NodeType(xmlNode *node) { return (int)node->type; }

const xmlChar * GoNodeName(xmlNode *node) { return node->name; }

xmlChar *
DumpNodeToXmlChar(xmlNode *node, xmlDoc *doc) {
  xmlBuffer *buff = xmlBufferCreate();
  xmlNodeDump(buff, doc, node, 0, 0);
  return buff->content;
}
*/
import "C"

type XmlNode struct {
	NodePtr *C.xmlNode
	DocRef  *XmlDoc
}

type Element struct {
	*XmlNode
}

func xmlNodeType(node *C.xmlNode) int {
	return int(C.NodeType(node))
}

func (node *XmlNode) Doc() *XmlDoc {
	return node.DocRef
}

func (node *XmlNode) Ptr() *C.xmlNode {
	return node.NodePtr
}

func (node *XmlNode) Type() int {
	return int(C.NodeType(node.Ptr()))
}

func (node *XmlNode) Search(xpath_expression string) *NodeSet {
	if node.Doc() == nil {
		println("Must define document in node")
	}
	ctx := node.Doc().XPathContext()
	ctx.SetNode(node)
	return ctx.EvalToNodes(xpath_expression)
}

func (node *XmlNode) Parent() Node {
	return buildNode(C.GoNodeParent(node.Ptr()), node.Doc())
}

func (node *XmlNode) Next() Node {
	return buildNode(C.GoNodeNext(node.Ptr()), node.Doc())
}

func (node *XmlNode) Prev() Node {
	return buildNode(C.GoNodePrev(node.Ptr()), node.Doc())
}

func (node *XmlNode) First() Node {
	return buildNode(C.GoNodeChildren(node.Ptr()), node.Doc())
}

func (node *XmlNode) Last() Node {
	return buildNode(C.GoNodeLast(node.Ptr()), node.Doc())
}

func (node *XmlNode) Remove() {
	C.xmlUnlinkNode(node.Ptr())
}

func (node *XmlNode) Name() string {
	return XmlChar2String(C.GoNodeName(node.Ptr()))
}

func (node *XmlNode) SetName(name string) {
	C.xmlNodeSetName(node.Ptr(), C.xmlCharStrdup(C.CString(name)))
}

func (node *XmlNode) Dump() string {
	return XmlChar2String(C.DumpNodeToXmlChar(node.Ptr(), node.Doc().DocPtr))
}

func (node *XmlNode) AttributeValue(name string) string {
	c := C.xmlCharStrdup(C.CString(name))
	s := C.xmlGetProp(node.Ptr(), c)
	return XmlChar2String(s)
}

func (node *XmlNode) SetAttributeValue(name string, value string) {
	c_name := C.xmlCharStrdup(C.CString(name))
	c_value := C.xmlCharStrdup(C.CString(value))
	C.xmlSetProp(node.Ptr(), c_name, c_value)
}

package libxml
/* 
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

type Node interface {
	Ptr() *C.xmlNode;
	Doc() *XmlDoc;    // reference to doc
	
	Dump() string;
	Search(xpath string) *NodeSet;
	Remove();
	
	// Standard libxml Node interface
	//Children() []Node;
	First()    Node;   // first child link
	Last()     Node;   // last child link
	Parent()   Node;   // child->parent link
	Next()     Node;   // next sibling link
	Prev()     Node;   // previous sibling link
	Type()     int;
	
	Name() string;
	SetName(name string);
	Attribute(name string) string;
	SetAttribute(name string, value string);
}

type XmlNode struct {
	NodePtr *C.xmlNode
	DocRef  *XmlDoc
}

func BuildNode(ptr *C.xmlNode, doc *XmlDoc) Node {
  if ptr == nil {
    return nil
  }
	node_type := int(C.NodeType(ptr))
	xml_node := &XmlNode{NodePtr: ptr, DocRef: doc}
	if node_type == 1 {
		return &XmlElement{XmlNode: xml_node}
	}
	return xml_node
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
  return BuildNode(C.GoNodeParent(node.Ptr()), node.Doc()) 
}

func (node *XmlNode) Next() Node { 
  return BuildNode(C.GoNodeNext(node.Ptr()), node.Doc()) 
}

func (node *XmlNode) Prev() Node { 
  return BuildNode(C.GoNodePrev(node.Ptr()), node.Doc()) 
}

func (node *XmlNode) First() Node { 
  return BuildNode(C.GoNodeChildren(node.Ptr()), node.Doc()) 
}

func (node *XmlNode) Last() Node { 
  return BuildNode(C.GoNodeLast(node.Ptr()), node.Doc()) 
}

func (node *XmlNode) Remove() {
  C.xmlUnlinkNode(node.Ptr())
}

func (node *XmlNode) Name() string { 
  return XmlChar2String(C.GoNodeName(node.Ptr()))
}

func (node *XmlNode) SetName(name string) {
	C.xmlNodeSetName(node.Ptr(), C.xmlCharStrdup( C.CString(name) ))
}

func (node *XmlNode) Dump() string {
	return XmlChar2String(C.DumpNodeToXmlChar(node.Ptr(), node.Doc().Ptr))
}

func (node *XmlNode) Attribute(name string) string { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  s := C.xmlGetProp(node.Ptr(), c) 
  return XmlChar2String(s)
}

func (node *XmlNode) SetAttribute(name string, value string) {
	c_name  := C.xmlCharStrdup( C.CString(name) ) 
	c_value := C.xmlCharStrdup( C.CString(value) ) 
	C.xmlSetProp(node.Ptr(), c_name, c_value)
}
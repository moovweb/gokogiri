package libxml
/* 
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
xmlNode * NodeNext(xmlNode *node) { return node->next; } 
xmlNode * NodeChildren(xmlNode *node) { return node->children; }
int NodeType(xmlNode *node) { return (int)node->type; }

xmlChar *
DumpNodeToXmlChar(xmlNode *node, xmlDoc *doc) {
  xmlBuffer *buff = xmlBufferCreate();
  xmlNodeDump(buff, doc, node, 0, 0);
  return buff->content;
}
*/
import "C"

type XmlNode interface {
	Search(xpath string) *XmlNodeSet
	Dump() string;
	DumpHTML() string;
}

type XmlBaseNode struct {
	Ptr	*C.xmlNode
	Doc	*XmlDoc
}

func BuildXmlNode(ptr *C.xmlNode, doc *XmlDoc) *XmlNode {
  if ptr == nil {
    return nil
  }
  node_type := C.NodeType(ptr)
	switch {
		case node_type == XML_ELEMENT_NODE:
        return &XmlElement{Ptr: ptr, Doc: doc}
    case node_type == XML_DOCUMENT_NODE:
        return &XmlDoc{Ptr: ptr, Doc: doc}
	}
  return &XmlBaseNode{Ptr: ptr, Doc: doc}
}

func (node *XmlNode) Type() int {
	return C.NodeType(node.ptr)
}

func (node *XmlNode) Search(xpath_expression string) *XmlNodeSet {
  if node.Doc == nil {
    println("Must define document in node")
  }
  ctx := node.Doc.XPathContext()
  ctx.SetNode(node)
  return ctx.EvalToNodes(xpath_expression)
}

func (node *XmlNode) Next() *XmlNode { 
  return BuildXmlNode(C.NodeNext(node.Ptr), node.Doc) 
}

func (node *XmlNode) Children() *XmlNode { 
  return BuildXmlNode(C.NodeChildren(node.Ptr), node.Doc) 
}

func (node *XmlNode) Remove() {
  C.xmlUnlinkNode(node.Ptr)
}

func (node *XmlNode) Dump() string {
	return XmlChar2String(C.DumpNodeToXmlChar(node.Ptr, node.Doc.Ptr))
}
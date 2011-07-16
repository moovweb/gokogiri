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
*/
import "C"

type XmlNode struct { 
  Ptr    *C.xmlNode
  Doc    *XmlDoc
}

func BuildXmlNode(ptr *C.xmlNode, doc *XmlDoc) *XmlNode {
  if ptr == nil {
    return nil
  }
  return &XmlNode{Ptr: ptr, Doc: doc}
}

func (node *XmlNode) GetProp(name string) string { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  s := C.xmlGetProp(node.Ptr, c) 
  return XmlChar2String(s)
}

func (node *XmlNode) Next() *XmlNode { 
  return BuildXmlNode(C.NodeNext(node.Ptr), node.Doc) 
}

func (node *XmlNode) Children() *XmlNode { 
  return BuildXmlNode(C.NodeChildren(node.Ptr), node.Doc) 
}

func (node *XmlNode) Name() string { 
  return XmlChar2String(node.Ptr.name)
}

func (node *XmlNode) Type() int { 
  return int(C.NodeType(node.Ptr)) 
}

func (node *XmlNode) Remove() {
  C.xmlUnlinkNode(node.Ptr)
}

func (node *XmlNode) Search(xpath_expression string) *XmlNodeSet {
  if node.Doc == nil {
    println("Must define document in node")
  }
  ctx := node.Doc.XPathContext()
  ctx.SetNode(node)
  return ctx.EvalToNodes(xpath_expression)
} 
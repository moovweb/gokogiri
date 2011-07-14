package libxml 
/* 
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
void xmlXPathContextSetNode(xmlXPathContext *ctx, xmlNode *new_node) { ctx->node = new_node; } 
xmlNode ** ObjectNode(xmlXPathObject *obj) { 
  xmlNodeSet *nodes;
  nodes = obj->nodesetval;
  if(nodes != NULL) {
    return (xmlNode **)nodes->nodeTab; 
  }
  return NULL;
}
xmlNode * FetchNode(xmlNode **nodes, int index) { return nodes[index]; }
int SizeOf(xmlNode **nodes) { return sizeof(*nodes) / sizeof(int); }
*/ 
import "C"

type XPathContext struct { 
  Ptr *C.xmlXPathContext
  Doc *XmlDoc
}

type XPathObject struct {
  Ptr *C.xmlXPathObject
  Doc *XmlDoc
}

func (context *XPathContext) RegisterNamespace(prefix, href string) bool {
  result := C.xmlXPathRegisterNs(context.Ptr, String2XmlChar(prefix), String2XmlChar(href))
  return result == 0
}

func (context *XPathContext) SetNode(node *XmlNode) {
  C.xmlXPathContextSetNode(context.Ptr, node.Ptr)
}

func (context *XPathContext) Eval(expression string) *XPathObject {
  object_pointer := C.xmlXPathEvalExpression(String2XmlChar(expression), context.Ptr)
  return &XPathObject{Ptr: object_pointer, Doc: context.Doc}
}

func (context *XPathContext) EvalToNodes(expression string) []XmlNode {
  obj := context.Eval(expression)
  return obj.NodeSet()
}

func (context *XPathContext) Free() {
  C.xmlXPathFreeContext(context.Ptr)
}

func (obj *XPathObject) NodeSet() []XmlNode {
  list := make([]XmlNode, 0, 100)
  nodes := C.ObjectNode(obj.Ptr)
  if nodes == nil {
    return list;
  }
  size := int(C.SizeOf(nodes));
  
  for i := 0; i < size; i++ {
    node := C.FetchNode(nodes, C.int(i));
    if node != nil {
      list = append(list, *BuildXmlNode(node, obj.Doc))
    }
  }
  //size  := C.int(C.sizeof(nodes))
  //println(size)
  //list := make([]XmlNode)
  return list
}
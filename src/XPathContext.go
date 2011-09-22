package libxml 
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
void 
xmlXPathContextSetNode(xmlXPathContext *ctx, xmlNode *new_node) { 
  ctx->node = new_node; } 

xmlNodeSet * 
FetchNodeSet(xmlXPathObject *obj) {
  return obj->nodesetval; }
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

func (context *XPathContext) SetNode(node Node) {
  C.xmlXPathContextSetNode(context.Ptr, node.Ptr())
}

func (context *XPathContext) Eval(expression string) *XPathObject {
  object_pointer := C.xmlXPathEvalExpression(String2XmlChar(expression), context.Ptr)
  return &XPathObject{Ptr: object_pointer, Doc: context.Doc}
}

func (context *XPathContext) EvalToNodes(expression string) *NodeSet {
  obj := context.Eval(expression)
  return obj.NodeSet()
}

func (context *XPathContext) Free() {
  C.xmlXPathFreeContext(context.Ptr)
}

func (obj *XPathObject) NodeSet() *NodeSet {
  return buildNodeSet(C.FetchNodeSet(obj.Ptr), obj.Doc)
}
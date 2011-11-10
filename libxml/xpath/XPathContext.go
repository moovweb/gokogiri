package xpath
/* 
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
void 
xmlXPathContextSetNode(xmlXPathContext *ctx, xmlNode *new_node) { 
  ctx->node = new_node; } 

xmlNodeSet * 
FetchNodeSet(xmlXPathObject *obj) {
  return obj->nodesetval; }
void xmlFreeXmlChars3(xmlChar* buf) { xmlFree(buf); }
*/
import "C"
import "unsafe"
import . "libxml/tree"

type XPathContext struct {
	Ptr *C.xmlXPathContext
	Doc *Doc
}

type XPathObject struct {
	Ptr *C.xmlXPathObject
	Doc *Doc
}

func ContextNew(node Node) *XPathContext {
	doc := node.Doc()
	docPtr := (*C.xmlDoc)(unsafe.Pointer(doc.Ptr()))
	ctx := &XPathContext{Ptr: C.xmlXPathNewContext(docPtr), Doc: doc}
	ctx.SetNode(node)
	return ctx
}

func Search(node Node, xpath_expression string) *NodeSet {
	if node.Doc() == nil {
		println("Must define document in node")
	}
	ctx := ContextNew(node)
  defer ctx.Free()
	return ctx.EvalToNodes(xpath_expression)
}

func (context *XPathContext) RegisterNamespace(prefix, href string) bool {
	cPrefix := C.xmlCharStrdup(C.CString(prefix))
  defer C.xmlFreeXmlChars3(cPrefix)
	cHref := C.xmlCharStrdup(C.CString(href))
  defer C.xmlFreeXmlChars3(cHref)
	result := C.xmlXPathRegisterNs(context.Ptr, cPrefix, cHref)
	return result == 0
}

func (context *XPathContext) SetNode(node Node) {
	C.xmlXPathContextSetNode(context.Ptr, (*C.xmlNode)(node.Ptr()))
}

func (context *XPathContext) Eval(expression string) *XPathObject {
	cExpression := C.xmlCharStrdup(C.CString(expression))
  defer C.xmlFreeXmlChars3(cExpression)
	object_pointer := C.xmlXPathEvalExpression(cExpression, context.Ptr)
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
  defer C.xmlXPathFreeObject(obj.Ptr)
	return NewNodeSet(unsafe.Pointer(C.FetchNodeSet(obj.Ptr)), obj.Doc)
}

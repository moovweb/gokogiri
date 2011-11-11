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

func Search(node Node, xpath_expression string) (*NodeSet, *XPathObject) {
	if node.Doc() == nil {
		println("Must define document in node")
	}
	ctx := ContextNew(node)
  defer ctx.Free()
	return ctx.EvalToNodes(xpath_expression)
}

func (context *XPathContext) RegisterNamespace(prefix, href string) bool {
  prefixCharPtr := C.CString(prefix)
  defer C.free(unsafe.Pointer(prefixCharPtr))
  prefixXmlCharPtr := C.xmlCharStrdup(prefixCharPtr)
  defer XmlFreeChars(unsafe.Pointer(prefixXmlCharPtr))

  hrefCharPtr := C.CString(href)
  defer C.free(unsafe.Pointer(hrefCharPtr))
  hrefXmlCharPtr := C.xmlCharStrdup(hrefCharPtr)
  defer XmlFreeChars(unsafe.Pointer(hrefXmlCharPtr))

	result := C.xmlXPathRegisterNs(context.Ptr, prefixXmlCharPtr, hrefXmlCharPtr)
	return result == 0
}

func (context *XPathContext) SetNode(node Node) {
	C.xmlXPathContextSetNode(context.Ptr, (*C.xmlNode)(node.Ptr()))
}

func (context *XPathContext) Eval(expression string) *XPathObject {
  expressionCharPtr := C.CString(expression)
  defer C.free(unsafe.Pointer(expressionCharPtr))
  expressionXmlCharPtr := C.xmlCharStrdup(expressionCharPtr)
  defer XmlFreeChars(unsafe.Pointer(expressionXmlCharPtr))

	object_pointer := C.xmlXPathEvalExpression(expressionXmlCharPtr, context.Ptr)
	return &XPathObject{Ptr: object_pointer, Doc: context.Doc}
}

func (context *XPathContext) EvalToNodes(expression string) (*NodeSet, *XPathObject) {
	obj := context.Eval(expression)
	return obj.NodeSet(), obj
}

func (context *XPathContext) Free() {
	C.xmlXPathFreeContext(context.Ptr)
}

func (obj *XPathObject) NodeSet() *NodeSet {
	return NewNodeSet(unsafe.Pointer(C.FetchNodeSet(obj.Ptr)), obj.Doc)
}

func (obj *XPathObject) Free() {
  C.xmlXPathFreeObject(obj.Ptr)
}

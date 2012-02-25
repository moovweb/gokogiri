package xpath
/* 
#cgo pkg-config: libxml-2.0

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
import "fmt"
import . "gokogiri/tree"

type XPath struct {
	context C.xmlXPathContextPtr
	result  C.xmlXPathObjectPtr
	Doc     *Doc
}

func NewXPath(doc *Doc) *XPath {
	if doc.Ptr() == nil {
		return nil
	}
	docPtr := (*C.xmlDoc)(unsafe.Pointer(doc.Ptr()))
	xpath := &XPath{context: C.xmlXPathNewContext(docPtr), result: nil, Doc: doc}
	return xpath
}


func (xpath *XPath) Search(node Node, xpathExp string) *NodeSet {
	exp := CompileXPath(xpathExp)
	if exp == nil {
		panic(fmt.Sprintf("cannot compile xpath: %q", xpathExp))
	}
	defer exp.Free()
	return xpath.SearchByCompiledXPath(node, exp)
}

func (xpath *XPath) SearchByCompiledXPath(node Node, exp *Expression) *NodeSet {
	if node.Doc().Ptr() != xpath.Doc.Ptr() {
		panic("this node's document does NOT match the document of the XPath context")
	}
	xpath.SetNode(node)
	if xpath.result != nil {
		//free the previous result if the XPath objecy is being reused.
		C.xmlXPathFreeObject(xpath.result)
	}
	xpath.result = C.xmlXPathCompiledEval(exp.ptr, xpath.context)
	return NewNodeSet(unsafe.Pointer(C.FetchNodeSet(xpath.result)), xpath.Doc)
}

func (xpath *XPath) RegisterNamespace(prefix, href string) bool {
	prefixCharPtr := C.CString(prefix)
	defer C.free(unsafe.Pointer(prefixCharPtr))
	prefixXmlCharPtr := C.xmlCharStrdup(prefixCharPtr)
	defer XmlFreeChars(unsafe.Pointer(prefixXmlCharPtr))

	hrefCharPtr := C.CString(href)
	defer C.free(unsafe.Pointer(hrefCharPtr))
	hrefXmlCharPtr := C.xmlCharStrdup(hrefCharPtr)
	defer XmlFreeChars(unsafe.Pointer(hrefXmlCharPtr))

	result := C.xmlXPathRegisterNs(xpath.context, prefixXmlCharPtr, hrefXmlCharPtr)
	return result == 0
}

func (xpath *XPath) SetNode(node Node) {
	C.xmlXPathContextSetNode(xpath.context, (*C.xmlNode)(node.Ptr()))
}

func (xpath *XPath) Free() {
	C.xmlXPathFreeContext(xpath.context)
	if xpath.result != nil {
		C.xmlXPathFreeObject(xpath.result)
	}
}

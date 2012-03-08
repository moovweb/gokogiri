package xpath
/* 
#cgo pkg-config: libxml-2.0
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
#include <libxml/parser.h>

void xmlXPathContextSetNode(xmlXPathContext *ctx, xmlNode *new_node) { 
	ctx->node = new_node;
}

xmlNode* fetchNode(xmlNodeSet *nodeset, int index) {
  	return nodeset->nodeTab[index];
}
*/
import "C"
import "unsafe"

type XPath struct {
	ContextPtr *C.xmlXPathContext
	ResultPtr  *C.xmlXPathObject
}

func NewXPath(docPtr unsafe.Pointer) (xpath *XPath) {
	if docPtr == nil {
		return
	}
	xpath = &XPath{ContextPtr: C.xmlXPathNewContext((*C.xmlDoc)(docPtr)), ResultPtr: nil}
	return
}

func (xpath *XPath) RegisterNamespace(prefix, href string) bool {
	var prefixPtr unsafe.Pointer = nil
	if len(prefix) > 0 {
		prefixBytes := []byte(prefix)
		prefixPtr = unsafe.Pointer(&prefixBytes[0])
	}
	
	var hrefPtr unsafe.Pointer = nil
	if len(href) > 0 {
		hrefBytes := []byte(href)
		hrefPtr = unsafe.Pointer(&hrefBytes[0])
	}

	result := C.xmlXPathRegisterNs(xpath.ContextPtr, (*C.xmlChar)(prefixPtr), (*C.xmlChar)(hrefPtr))
	return result == 0
}

func (xpath *XPath) Evaluate(nodePtr unsafe.Pointer, xpathExpr *Expression) (nodes []unsafe.Pointer){
	if nodePtr == nil {
		return
	}
	C.xmlXPathContextSetNode(xpath.ContextPtr, (*C.xmlNode)(nodePtr))
	if xpath.ResultPtr != nil {
		C.xmlXPathFreeObject(xpath.ResultPtr)
	}
	xpath.ResultPtr = C.xmlXPathCompiledEval(xpathExpr.Ptr, xpath.ContextPtr)
	if nodesetPtr := xpath.ResultPtr.nodesetval; nodesetPtr != nil {
		if nodesetSize := int(nodesetPtr.nodeNr); nodesetSize > 0 {
			nodes = make([]unsafe.Pointer, nodesetSize)
			for i := 0; i < nodesetSize; i ++ {
				nodes[i] = unsafe.Pointer(C.fetchNode(nodesetPtr, C.int(i)))
			}
		}
	}
	return
}

func (xpath *XPath) Free() {
	if xpath.ContextPtr != nil {
		C.xmlXPathFreeContext(xpath.ContextPtr)
	}
	if xpath.ResultPtr != nil {
		C.xmlXPathFreeObject(xpath.ResultPtr)
	}
}
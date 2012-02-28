package xpath
/* 
#cgo pkg-config: libxml-2.0
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
*/
import("C")
import "unsafe"
import . "gokogiri/libxml/tree"

type Expression struct {
	ptr C.xmlXPathCompExprPtr
}

func CompileXPath(xpathExp string) (*Expression) {
	expressionCharPtr := C.CString(xpathExp)
	defer C.free(unsafe.Pointer(expressionCharPtr))
	expressionXmlCharPtr := C.xmlCharStrdup(expressionCharPtr)
	defer XmlFreeChars(unsafe.Pointer(expressionXmlCharPtr))
	ptr := C.xmlXPathCompile(expressionXmlCharPtr)
	if ptr == nil {
		return nil
	}
	return &Expression{ptr: ptr}
}

func (exp *Expression) Free() {
	C.xmlXPathFreeCompExpr(exp.ptr)
}

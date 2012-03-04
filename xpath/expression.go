package xpath
/*
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
*/
import "C"
import "unsafe"

type Expression struct {
	Ptr *C.xmlXPathCompExpr
}

func Compile(xpath string) (expr *Expression) {
	if len(xpath) == 0 {
		return
	}
	
	xpathBytes := []byte(xpath)
	xpathPtr := unsafe.Pointer(&xpathBytes[0])
	ptr := C.xmlXPathCompile((*C.xmlChar)(xpathPtr))
	if ptr == nil {
		return
	}
	expr = &Expression{Ptr: ptr}
	return
}

func (exp *Expression) Free() {
	if exp.Ptr != nil {
		C.xmlXPathFreeCompExpr(exp.Ptr)
	}
}
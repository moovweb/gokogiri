package xpath
/*
#cgo pkg-config: libxml-2.0
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
*/
import "C"
import "unsafe"
import . "gokogiri/util"

type Expression struct {
	Ptr *C.xmlXPathCompExpr
}

func Compile(path string) (expr *Expression) {
	if len(path) == 0 {
		return
	}

	xpathBytes := AppendCStringTerminator([]byte(path))
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

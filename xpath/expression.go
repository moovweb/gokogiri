package xpath
/*
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
*/
import "C"

import (
	"unsafe"
	. "gokogiri/util"

	"time"
)

type Expression struct {
	Ptr *C.xmlXPathCompExpr
}

var (
	CompileCount int64
	CompileTime int64
)

func Compile(path string) (expr *Expression) {

	CompileCount++
	startTime := time.Now().UnixNano()

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

	CompileTime += time.Now().UnixNano() - startTime

	return
}

func (exp *Expression) Free() {
	if exp.Ptr != nil {
		C.xmlXPathFreeCompExpr(exp.Ptr)
	}
}

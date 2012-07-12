package xpath

/*
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>

void check_xpath_syntax_noop() {
}

char *check_xpath_syntax(const char *xpath) {
	char *rval = NULL;
	xmlXPathContextPtr ctx = xmlXPathNewContext(NULL);
	ctx->error = check_xpath_syntax_noop;
	xmlXPathCtxtCompile(ctx, (const xmlChar *)xpath);
	if (ctx->lastError.domain > 0) {
		//fprintf(stderr, "%s Code: %d Domain: %d", ctx->lastError.str2, ctx->lastError.code, ctx->lastError.domain);
		rval = "ERROR";
	}
	xmlXPathFreeContext(ctx);
	return rval;
}
*/
import "C"
import "unsafe"
import . "gokogiri/util"
import "runtime"
import "errors"

type Expression struct {
	Ptr *C.xmlXPathCompExpr
}

func Check(path string) (err error) {
	str := C.CString(path)
	defer C.free(unsafe.Pointer(str))
	cstr := C.check_xpath_syntax(str)
	if cstr != nil {
		err = errors.New(C.GoString(cstr))
	}
	return
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
	runtime.SetFinalizer(expr, (*Expression).Free)
	return
}

func (exp *Expression) Free() {
	if exp.Ptr != nil {
		C.xmlXPathFreeCompExpr(exp.Ptr)
		exp.Ptr = nil
	}
}

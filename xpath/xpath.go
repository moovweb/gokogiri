package xpath

/*
#cgo CFLAGS: -I../../../clibs/include/libxml2
#cgo LDFLAGS: -lxml2 -L../../../clibs/lib
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <libxml/parser.h>

xmlNode* fetchNode(xmlNodeSet *nodeset, int index) {
  	return nodeset->nodeTab[index];
}

xmlXPathObjectPtr go_resolve_variables(void* ctxt, char* name, char* ns);

static void set_var_lookup(xmlXPathContext* c, void* data) {
    c->varLookupFunc = (void *)go_resolve_variables;
    c->varLookupData = data;
}

int getXPathObjectType(xmlXPathObject* o) {
    if(o == 0)
        return 0;
    return o->type;
}

*/
import "C"

import "time"
import "unsafe"
import . "gokogiri/util"
import "runtime"
import "errors"

type XPath struct {
	ContextPtr *C.xmlXPathContext
	ResultPtr  *C.xmlXPathObject
}

type XPathObjectType int

const (
	XPATH_UNDEFINED   XPathObjectType = 0
	XPATH_NODESET                     = 1
	XPATH_BOOLEAN                     = 2
	XPATH_NUMBER                      = 3
	XPATH_STRING                      = 4
	XPATH_POINT                       = 5
	XPATH_RANGE                       = 6
	XPATH_LOCATIONSET                 = 7
	XPATH_USERS                       = 8
	XPATH_XSLT_TREE                   = 9 // An XSLT value tree, non modifiable
)

// Types that provide the VariableScope interface know how to resolve
// XPath variable names into values.

//This interface exist primarily for the benefit of XSLT processors.
type VariableScope interface {
	Resolve(string, string) interface{}
}

func NewXPath(docPtr unsafe.Pointer) (xpath *XPath) {
	if docPtr == nil {
		return
	}
	xpath = &XPath{ContextPtr: C.xmlXPathNewContext((*C.xmlDoc)(docPtr)), ResultPtr: nil}
	runtime.SetFinalizer(xpath, (*XPath).Free)
	return
}

func (xpath *XPath) RegisterNamespace(prefix, href string) bool {
	var prefixPtr unsafe.Pointer = nil
	if len(prefix) > 0 {
		prefixBytes := AppendCStringTerminator([]byte(prefix))
		prefixPtr = unsafe.Pointer(&prefixBytes[0])
	}

	var hrefPtr unsafe.Pointer = nil
	if len(href) > 0 {
		hrefBytes := AppendCStringTerminator([]byte(href))
		hrefPtr = unsafe.Pointer(&hrefBytes[0])
	}

	result := C.xmlXPathRegisterNs(xpath.ContextPtr, (*C.xmlChar)(prefixPtr), (*C.xmlChar)(hrefPtr))
	return result == 0
}

// Evaluate an XPath and attempt to consume the result as a nodeset.
func (xpath *XPath) EvaluateAsNodeset(nodePtr unsafe.Pointer, xpathExpr *Expression) (nodes []unsafe.Pointer, err error) {
	if nodePtr == nil {
		//evaluating xpath on a  nil node returns no result.
		return
	}

	err = xpath.Evaluate(nodePtr, xpathExpr)
	if err != nil {
		return
	}

	nodes, err = xpath.ResultAsNodeset()
	return
}

// Evaluate an XPath. The returned result is stored in the struct. Call ReturnType to
// discover the type of result, and call one of the ResultAs* functions to return a
// copy of the result as a particular type.
func (xpath *XPath) Evaluate(nodePtr unsafe.Pointer, xpathExpr *Expression) (err error) {
	if nodePtr == nil {
		//evaluating xpath on a nil node returns no result.
		return
	}
	xpath.ContextPtr.node = (*C.xmlNode)(nodePtr)
	if xpath.ResultPtr != nil {
		C.xmlXPathFreeObject(xpath.ResultPtr)
	}

	xpath.ResultPtr = C.xmlXPathCompiledEval(xpathExpr.Ptr, xpath.ContextPtr)
	if xpath.ResultPtr == nil {
		err = errors.New("err in evaluating xpath: " + xpathExpr.String())
		return
	}
	return
}

// Determine the actual return type of the XPath evaluation.
func (xpath *XPath) ReturnType() XPathObjectType {
	return XPathObjectType(C.getXPathObjectType(xpath.ResultPtr))
}

// Get the XPath result as a nodeset.
func (xpath *XPath) ResultAsNodeset() (nodes []unsafe.Pointer, err error) {
	if xpath.ResultPtr == nil {
		return
	}

	if xpath.ReturnType() != XPATH_NODESET {
		err = errors.New("Cannot convert XPath result to nodeset")
	}

	if nodesetPtr := xpath.ResultPtr.nodesetval; nodesetPtr != nil {
		if nodesetSize := int(nodesetPtr.nodeNr); nodesetSize > 0 {
			nodes = make([]unsafe.Pointer, nodesetSize)
			for i := 0; i < nodesetSize; i++ {
				nodes[i] = unsafe.Pointer(C.fetchNode(nodesetPtr, C.int(i)))
			}
		}
	}
	return
}

// Coerce the result into a string
func (xpath *XPath) ResultAsString() (val string, err error) {
	if xpath.ReturnType() != XPATH_STRING {
		xpath.ResultPtr = C.xmlXPathConvertString(xpath.ResultPtr)
	}
	val = C.GoString((*C.char)(unsafe.Pointer(xpath.ResultPtr.stringval)))
	return
}

// Coerce the result into a number
func (xpath *XPath) ResultAsNumber() (val float64, err error) {
	if xpath.ReturnType() != XPATH_NUMBER {
		xpath.ResultPtr = C.xmlXPathConvertNumber(xpath.ResultPtr)
	}
	val = float64(xpath.ResultPtr.floatval)
	return
}

// Add a variable resolver.
func (xpath *XPath) SetResolver(v VariableScope) {
	C.set_var_lookup(xpath.ContextPtr, unsafe.Pointer(&v))
}

func (xpath *XPath) SetDeadline(deadline *time.Time) {
	if deadline == nil {
		C.xmlXPathContextSetDeadline(xpath.ContextPtr, C.time_t(0))
	} else {
		t := deadline.Unix()
		C.xmlXPathContextSetDeadline(xpath.ContextPtr, C.time_t(t))
	}
}

func (xpath *XPath) Free() {
	if xpath.ContextPtr != nil {
		C.xmlXPathFreeContext(xpath.ContextPtr)
		xpath.ContextPtr = nil
	}
	if xpath.ResultPtr != nil {
		C.xmlXPathFreeObject(xpath.ResultPtr)
		xpath.ResultPtr = nil
	}
}

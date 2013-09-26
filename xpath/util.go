package xpath

/*
#cgo pkg-config: libxml-2.0
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <libxml/parser.h>

*/
import "C"

import "unsafe"
import "fmt"
import "reflect"
import . "github.com/moovweb/gokogiri/util"

//export go_resolve_variables
func go_resolve_variables(ctxt unsafe.Pointer, name, ns *C.char) (ret C.xmlXPathObjectPtr) {
	variable := C.GoString(name)
	namespace := C.GoString(ns)

	context := (*VariableScope)(ctxt)
	if context != nil {
		val := (*context).Resolve(variable, namespace)
		if val == nil {
			//fmt.Println("go-resolve wrong-type nil")
			//return the empty node set
			ret = C.xmlXPathNewNodeSet(nil)
			return
		}
		switch val.(type) {
		case []unsafe.Pointer:
			ptrs := val.([]unsafe.Pointer)
			if len(ptrs) > 0 {
				//default - return a node set
				ret = C.xmlXPathNewNodeSet(nil)
				for _, p := range ptrs {
					_ = C.xmlXPathNodeSetAdd(ret.nodesetval, (*C.xmlNode)(p))
				}
			} else {
				ret = C.xmlXPathNewNodeSet(nil)
				return
			}
		case float64:
			content := val.(float64)
			ret = C.xmlXPathNewFloat(C.double(content))
		case string:
			content := val.(string)
			fmt.Println("go-resolve string", content)
			xpathBytes := GetCString([]byte(content))
			xpathPtr := unsafe.Pointer(&xpathBytes[0])
			ret = C.xmlXPathNewString((*C.xmlChar)(xpathPtr))
		default:
			typ := reflect.TypeOf(val)
			// if a pointer to a struct is passed, get the type of the dereferenced object
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			fmt.Println("go-resolve wrong-type", typ.Kind())
		}
	}
	return
}

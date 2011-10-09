package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h>
*/
import "C"
import "unsafe"

func XmlChar2String(chars *C.xmlChar) string {
	return C.GoString((*C.char)(unsafe.Pointer(chars)))
}

func String2XmlChar(input string) *C.xmlChar {
	cString := C.CString(input)

	defer C.free(unsafe.Pointer(cString))
	return C.xmlCharStrdup(cString)
}

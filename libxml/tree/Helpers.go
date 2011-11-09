package tree
/* 
#include <libxml/tree.h>
//xmlFree is not really a function but a macro in libxml2, so we have to define a function like the following and thus cgo can use it
void xmlFreeChars(char* buf) { xmlFree((xmlChar*)buf); } 
*/
import "C"
import "unsafe"

func XmlChar2String(chars *C.xmlChar) string {
  cPtr := (*C.char)(unsafe.Pointer(chars))
	return C.GoString(cPtr)
}

func String2XmlChar(input string) *C.xmlChar {
	cString := C.CString(input)

	defer C.free(unsafe.Pointer(cString))
	return C.xmlCharStrdup(cString)
}

func XmlFreeChars(chars *C.char) {
  C.xmlFreeChars(chars)
}

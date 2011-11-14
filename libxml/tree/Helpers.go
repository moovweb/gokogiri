package tree
/* 
#include <libxml/tree.h>
//xmlFree is not really a function but a macro in libxml2, so we have to define a function like the following and thus cgo can use it
void xmlFreeChars(void* buf) { xmlFree((xmlChar*)buf); } 
*/
import "C"
import "unsafe"

func XmlChar2String(xmlCharPtr *C.xmlChar) string {
	cCharPtr := (*C.char)(unsafe.Pointer(xmlCharPtr))
	return C.GoString(cCharPtr)
}

func String2XmlChar(str string) *C.xmlChar {
	cCharPtr := C.CString(str)
	defer C.free(unsafe.Pointer(cCharPtr))
	return C.xmlCharStrdup(cCharPtr)
}

func XmlFreeChars(chars unsafe.Pointer) {
	C.xmlFreeChars(chars)
}

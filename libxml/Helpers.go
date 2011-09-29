package libxml
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/xmlstring.h> 
char* xmlChar2C(xmlChar* x) { return (char *) x; }
xmlChar* C2xmlChar(char* x) { return (xmlChar *) x; }
*/
import "C"
//import "unsafe"

func XmlCheckVersion() int {
	var v C.int
	C.xmlCheckVersion(v)
	return int(v)
}

func XmlCleanUpParser() {
	C.xmlCleanupParser()
}

func XmlChar2String(s *C.xmlChar) string {
	cString := C.xmlChar2C(s)
	//defer C.free(unsafe.Pointer(cString))
	return C.GoString(cString)
}

func String2XmlChar(s string) *C.xmlChar {
	return C.C2xmlChar(C.CString(s))
}

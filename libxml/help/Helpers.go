package help
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

func XmlCheckVersion() int {
	var v C.int
	C.xmlCheckVersion(v)
	return int(v)
}

func XmlCleanUpParser() {
	C.xmlCleanupParser()
}
package html
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h>

xmlDoc* 
htmlDocToXmlDoc(htmlDocPtr doc) { return (xmlDocPtr)doc; }

*/
import "C"
import "unsafe"
import . "libxml/tree"

func ParseStringWithOptions(content string, url string, encoding string, opts int) *Doc {
	cString := C.CString(content)
	cXmlChar := C.xmlCharStrdup(cString)
	htmlDocPtr := C.htmlReadDoc(cXmlChar, C.CString(url), C.CString(encoding), C.int(opts))
	xmlDocPtr := C.htmlDocToXmlDoc(htmlDocPtr)
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func ParseString(content string) *Doc {
	return ParseStringWithOptions(content, "", "",
		HTML_PARSE_RECOVER|
			HTML_PARSE_NONET|
			HTML_PARSE_NOERROR|
			HTML_PARSE_NOWARNING)
}

func ParseFile(url string, encoding string, opts int) *Doc {
	htmlDocPtr := C.htmlReadFile(C.CString(url), C.CString(encoding), C.int(opts))
	xmlDocPtr := C.htmlDocToXmlDoc(htmlDocPtr)
	return NewDoc(xmlDocPtr)
}

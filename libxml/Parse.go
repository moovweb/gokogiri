package libxml
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

func HtmlParseStringWithOptions(content string, url string, encoding string, opts int) *Doc {
	cString := C.CString(content)
	cXmlChar := C.xmlCharStrdup(cString)
	htmlDocPtr := C.htmlReadDoc(cXmlChar, C.CString(url), C.CString(encoding), C.int(opts))
	xmlDocPtr := C.htmlDocToXmlDoc(htmlDocPtr)
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func HtmlParseString(content string) *Doc {
	return HtmlParseStringWithOptions(content, "", "",
		HTML_PARSE_RECOVER|
			HTML_PARSE_NONET|
			HTML_PARSE_NOERROR|
			HTML_PARSE_NOWARNING)
}

func HtmlParseFile(url string, encoding string, opts int) *Doc {
	htmlDocPtr := C.htmlReadFile(C.CString(url), C.CString(encoding), C.int(opts))
	xmlDocPtr := C.htmlDocToXmlDoc(htmlDocPtr)
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func XmlParseWithOption(content string, url string, encoding string, opts int) *Doc {
	c := C.xmlCharStrdup(C.CString(content))
	c_encoding := C.CString(encoding)
	if encoding == "" {
		c_encoding = nil
	}
	xmlDocPtr := C.xmlReadDoc(c, C.CString(url), c_encoding, C.int(opts))
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func XmlParseString(content string) *Doc {
	return XmlParseWithOption(content, "", "", 0)
}

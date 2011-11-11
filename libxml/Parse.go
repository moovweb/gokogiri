package libxml
/*
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h>

xmlDoc* 
htmlDocToXmlDoc(htmlDocPtr doc) { return (xmlDocPtr)doc; }
*/
import "C"
import "unsafe"
import . "libxml/tree"

func HtmlParseStringWithOptions(content string, url string, encoding string, opts int) *Doc {
	contentCharPtr := C.CString(content)
	defer C.free(unsafe.Pointer(contentCharPtr))
	contentXmlCharPtr := C.xmlCharStrdup(contentCharPtr)
	defer XmlFreeChars(unsafe.Pointer(contentXmlCharPtr))
	urlCharPtr := C.CString(url)
	defer C.free(unsafe.Pointer(urlCharPtr))

	var encodingCharPtr *C.char = nil
	if encoding != "" {
		encodingCharPtr = C.CString(encoding)
		defer C.free(unsafe.Pointer(encodingCharPtr))
	}
	htmlDocPtr := C.htmlReadDoc(contentXmlCharPtr, urlCharPtr, encodingCharPtr, C.int(opts))
	if htmlDocPtr == nil {
		return nil
	}
	xmlDocPtr := C.htmlDocToXmlDoc(htmlDocPtr)
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func HtmlParseString(content string) *Doc {
	doc := HtmlParseStringWithOptions(content, "", "",
		HTML_PARSE_RECOVER|
			HTML_PARSE_NONET|
			HTML_PARSE_NOERROR|
			HTML_PARSE_NOWARNING)
	if doc == nil {
		return HtmlParseString("<html />")
	}
	return doc
}

func HtmlParseFile(url string, encoding string, opts int) *Doc {
	htmlDocPtr := C.htmlReadFile(C.CString(url), C.CString(encoding), C.int(opts))
	xmlDocPtr := C.htmlDocToXmlDoc(htmlDocPtr)
	if xmlDocPtr == nil {
		return nil
	}
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func XmlParseWithOption(content string, url string, encoding string, opts int) *Doc {
	contentCharPtr := C.CString(content)
	defer C.free(unsafe.Pointer(contentCharPtr))
	contentXmlCharPtr := C.xmlCharStrdup(contentCharPtr)
	defer XmlFreeChars(unsafe.Pointer(contentXmlCharPtr))
	urlCharPtr := C.CString(url)
	defer C.free(unsafe.Pointer(urlCharPtr))

	var encodingCharPtr *C.char = nil
	if encoding != "" {
		encodingCharPtr = C.CString(encoding)
		defer C.free(unsafe.Pointer(encodingCharPtr))
	}
	xmlDocPtr := C.xmlReadDoc(contentXmlCharPtr, urlCharPtr, encodingCharPtr, C.int(opts))
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func XmlParseString(content string) *Doc {
	return XmlParseWithOption(content, "", "", 1)
}

package tree
/* 
#include <libxml/tree.h>
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h>

xmlNode * GoXmlCastDocToNode(xmlDoc *doc) { return (xmlNode *)doc; }
xmlDoc * htmlDocToXmlDoc(htmlDocPtr doc) { return (xmlDocPtr)doc; }
xmlDoc * newEmptyXmlDoc2() { return xmlNewDoc(BAD_CAST XML_DEFAULT_VERSION); }
*/
import "C"
import "unsafe"
import "strings"

/*
func Parse(input string) *Doc {
	cCharInput := C.CString(input)
	defer C.free(unsafe.Pointer(cCharInput))
	doc := C.xmlParseMemory(cCharInput, C.int(len(input)))
	return NewNode(unsafe.Pointer(doc), nil).(*Doc)
}*/

func XmlParseWithOptions(content string, url string, encoding string, opts int) *Doc {
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
	if xmlDocPtr == nil {
		xmlDocPtr = C.newEmptyXmlDoc2()
	}
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

// Returns the first element in the input string.
// Use Next() to access siblings
func XmlParseFragmentWithOptions(input string, url string, encoding string, opts int) *Doc {
	return XmlParseWithOptions("<root>" + input + "</root>", url, encoding, opts)
}

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


func HtmlParseFile(url string, encoding string, opts int) *Doc {
	htmlDocPtr := C.htmlReadFile(C.CString(url), C.CString(encoding), C.int(opts))
	xmlDocPtr := C.htmlDocToXmlDoc(htmlDocPtr)
	if xmlDocPtr == nil {
		return nil
	}
	return NewDoc(unsafe.Pointer(xmlDocPtr))
}

func HtmlParseString(content string, encoding string) *Doc {
	doc := HtmlParseStringWithOptions(content, "", encoding, DefaultHtmlParseOptions())
	if doc == nil {
		return HtmlParseString("<html />", "")
	}
	return doc
}

func XmlParseString(content string, encoding string) *Doc {
	return XmlParseWithOptions(content, "", encoding, DefaultXmlParseOptions())
}

func XmlParseFragment(content string, encoding string) *Doc {
	return XmlParseFragmentWithOptions(content, "", encoding, DefaultXmlParseOptions())
}

func HtmlParseFragment(content string, encoding string) *Doc {
	htmlDoc := HtmlParseStringWithOptions("<html><body>"+content, "", encoding, DefaultHtmlParseOptions())
	html := htmlDoc.RootElement()
	body := html.First()
	if strings.Index(strings.ToLower(content), "<body") < 0 {
		child := body.First()
		for child != nil {
			nextChild := child.Next()
			html.AppendChildNode(child)
			child = nextChild
		}
		body.Free()
	}
	return htmlDoc
}
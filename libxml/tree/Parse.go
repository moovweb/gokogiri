package tree
/* 
#include <libxml/tree.h>
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h>

xmlNode * GoXmlCastDocToNode(xmlDoc *doc) { return (xmlNode *)doc; }
xmlDoc * htmlDocToXmlDoc(htmlDocPtr doc) { return (xmlDocPtr)doc; }
xmlDoc * newEmptyXmlDoc() { return xmlNewDoc(BAD_CAST XML_DEFAULT_VERSION); }
htmlDocPtr	htmlReadDoc2(const void * cur, const void * URL, const void * encoding, int options) {
	return htmlReadDoc((xmlChar*)cur, (char*)URL, (char*)encoding, options);
}
*/
import "C"
import "unsafe"
import "strings"

func Parse(input string) *Doc {
	cCharInput := C.CString(input)
	defer C.free(unsafe.Pointer(cCharInput))
	doc := C.xmlParseMemory(cCharInput, C.int(len(input)))
	return NewNode(unsafe.Pointer(doc), nil).(*Doc)
}

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
		xmlDocPtr = C.newEmptyXmlDoc()
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

func HtmlParseBytesWithOptions(content []byte, url []byte, encoding []byte, opts int) *Doc {
	contentCharPtr := unsafe.Pointer(&content[0])
	/*
	urlCharPtr := unsafe.Pointer(nil)
	if url != nil {
		urlCharPtr = unsafe.Pointer(&url[0])
	}
	encodingCharPtr := unsafe.Pointer(nil)
	if encoding != nil {
		encodingCharPtr = unsafe.Pointer(&encoding[0])
	}*/
	htmlDocPtr := C.htmlReadDoc2(contentCharPtr, nil, nil, C.int(opts))
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

func HtmlParseString(content string) *Doc {
	doc := HtmlParseStringWithOptions(content, "", "", DefaultHtmlParseOptions())
	if doc == nil {
		return HtmlParseString("<html />")
	}
	return doc
}

func HtmlParseBytes(content []byte) *Doc {
	doc := HtmlParseBytesWithOptions(content, nil, nil, DefaultHtmlParseOptions())
	if doc == nil {
		return HtmlParseString("<html />")
	}
	return doc
}

func XmlParseString(content string) *Doc {
	return XmlParseWithOptions(content, "", "", DefaultXmlParseOptions())
}

func XmlParseFragment(content string) *Doc {
	return XmlParseFragmentWithOptions(content, "", "", DefaultXmlParseOptions())
}

func HtmlParseFragment(content string) *Doc {
	doc := XmlParseString("<root></root>")
	tmpDoc := HtmlParseStringWithOptions("<html><body>"+content, "", "", DefaultHtmlParseOptions())
	defer tmpDoc.Free()

	tmpNode := tmpDoc.RootElement().First()
	if strings.Index(strings.ToLower(content), "<body") < 0 {
		tmpNode = tmpNode.First()
	}

	//append all children of tmpRoot to root.
	root := doc.RootElement()
	child := tmpNode
	for child != nil {
		nextChild := child.Next()
		root.AppendChildNode(child)
		child = nextChild
	}
	return doc
}



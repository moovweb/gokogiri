package tree
/* 
#include <libxml/tree.h>
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h>

xmlNode * GoXmlCastDocToNode(xmlDoc *doc) { return (xmlNode *)doc; }
xmlDoc * htmlDocToXmlDoc(htmlDocPtr doc) { return (xmlDocPtr)doc; }
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

func HtmlParseString(content string) *Doc {
	doc := HtmlParseStringWithOptions(content, "", "", DefaultHtmlParseOptions())
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
	tmpDoc := HtmlParseStringWithOptions(content, "", "", DefaultHtmlParseOptions())
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



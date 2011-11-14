package tree
/* 
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
#include <libxml/xpath.h> 

char *
DumpXmlToString(xmlDoc *doc) {
  xmlChar *buff;
  int buffersize;
  xmlDocDumpFormatMemory(doc, 
                         &buff,
                         &buffersize, 1);
  return (char *)buff;
}

char *
DumpHtmlToString(xmlDoc *doc) {
  xmlChar *buff;
  int buffersize;
  htmlDocDumpMemory(doc, &buff, &buffersize);
  return (char *)buff;
}

xmlNode * GoXmlCastDocToNode(xmlDoc *doc) { return (xmlNode *)doc; }
xmlDoc * htmlDocToXmlDoc(htmlDocPtr doc) { return (xmlDocPtr)doc; }
*/
import "C"
import "unsafe"

type Doc struct {
	DocPtr *C.xmlDoc
	*XmlNode
}

func Parse(input string) *Doc {
	cCharInput := C.CString(input)
	defer C.free(unsafe.Pointer(cCharInput))
	doc := C.xmlParseMemory(cCharInput, C.int(len(input)))
	return NewNode(unsafe.Pointer(doc), nil).(*Doc)
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

func CreateHtmlDoc() *Doc {
	cDoc := C.htmlNewDoc(String2XmlChar(""), String2XmlChar(""))
	return NewNode(unsafe.Pointer(cDoc), nil).(*Doc)
}

// Returns the first element in the input string.
// Use Next() to access siblings
func (doc *Doc) ParseFragment(input string) *Doc {
	newDoc := Parse("<root>" + input + "</root>")
	return newDoc
}

func NewDoc(ptr unsafe.Pointer) *Doc {
	doc := NewNode(ptr, nil).(*Doc)
	doc.DocPtr = (*C.xmlDoc)(ptr)
	return doc
}

func (doc *Doc) NewElement(name string) *Element {
	nameXmlCharPtr := String2XmlChar(name)
	defer XmlFreeChars(unsafe.Pointer(nameXmlCharPtr))
	return NewNode(unsafe.Pointer(C.xmlNewNode(nil, nameXmlCharPtr)), doc).(*Element)
}

func (doc *Doc) Free() {
	C.xmlFreeDoc(doc.DocPtr)
}

func (doc *Doc) MetaEncoding() string {
	metaEncodingXmlCharPtr := C.htmlGetMetaEncoding(doc.DocPtr)
	return C.GoString((*C.char)(unsafe.Pointer(metaEncodingXmlCharPtr)))
}

func (doc *Doc) String() string {
	// TODO: Decide what type of return to do HTML or XML
	dumpCharPtr := C.DumpXmlToString(doc.DocPtr)
	defer XmlFreeChars(unsafe.Pointer(dumpCharPtr))
	return C.GoString(dumpCharPtr)
}

func (doc *Doc) DumpHTML() string {
	dumpCharPtr := C.DumpHtmlToString(doc.DocPtr)
	defer XmlFreeChars(unsafe.Pointer(dumpCharPtr))
	return C.GoString(dumpCharPtr)
}

func (doc *Doc) DumpXML() string {
	dumpCharPtr := C.DumpXmlToString(doc.DocPtr)
	defer XmlFreeChars(unsafe.Pointer(dumpCharPtr))
	return C.GoString(dumpCharPtr)
}

func (doc *Doc) RootElement() *Element {
	return NewNode(unsafe.Pointer(C.xmlDocGetRootElement(doc.DocPtr)), doc).(*Element)
}

func (doc *Doc) NewCData(content string) *CData {
	length := C.int(len([]byte(content)))
	cData := C.xmlNewCDataBlock(doc.DocPtr, String2XmlChar(content), length)
	return NewNode(unsafe.Pointer(cData), doc).(*CData)
}

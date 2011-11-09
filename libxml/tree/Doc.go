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
*/
import "C"
import "unsafe"

type Doc struct {
	DocPtr *C.xmlDoc
	*XmlNode
}

func Parse(input string) *Doc {
	cInput := C.CString(input)
	doc := C.xmlParseMemory(cInput, C.int(len(input)))
	return NewNode(unsafe.Pointer(doc), nil).(*Doc)
}

func CreateHtmlDoc() *Doc {
	cDoc := C.htmlNewDoc(String2XmlChar(""), String2XmlChar(""))
	return NewNode(unsafe.Pointer(cDoc), nil).(*Doc)
}

// Returns the first element in the input string.
// Use Next() to access siblings
func (doc *Doc) ParseFragment(input string) Node {
	newDoc := Parse("<root>" + input + "</root>")
	defer newDoc.Free()
	res := newDoc.First().First()
	res.SetDoc(doc)
	return res
}

func NewDoc(ptr unsafe.Pointer) *Doc {
	doc := NewNode(ptr, nil).(*Doc)
	doc.DocPtr = (*C.xmlDoc)(ptr)
	return doc
}

func (doc *Doc) NewElement(named string) *Element {
	return NewNode(unsafe.Pointer(C.xmlNewNode(nil, String2XmlChar(named))), doc).(*Element)
}

func (doc *Doc) Free() {
	C.xmlFreeDoc(doc.DocPtr)
}

func (doc *Doc) MetaEncoding() string {
	return C.GoString((*C.char)(unsafe.Pointer(C.htmlGetMetaEncoding(doc.DocPtr))))
}

func (doc *Doc) String() string {
	// TODO: Decide what type of return to do HTML or XML
	cString := C.DumpXmlToString(doc.DocPtr)
	defer C.free(unsafe.Pointer(cString))
	return C.GoString(cString)
}

func (doc *Doc) DumpHTML() string {
	return C.GoString(C.DumpHtmlToString(doc.DocPtr))
}

func (doc *Doc) DumpXML() string {
	return C.GoString(C.DumpXmlToString(doc.DocPtr))
}

func (doc *Doc) RootElement() *Element {
	return NewNode(unsafe.Pointer(C.xmlDocGetRootElement(doc.DocPtr)), doc).(*Element)
}

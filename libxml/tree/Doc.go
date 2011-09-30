package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
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

func NewDoc(ptr unsafe.Pointer) *Doc {
	doc := NewNode((*C.xmlNode)(ptr), nil).(*Doc)
	doc.DocPtr = (*C.xmlDoc)(ptr)
	return doc
}

func (doc *Doc) Free() {
	C.xmlFreeDoc(doc.DocPtr)
}

func (doc *Doc) MetaEncoding() string {
	return C.GoString((*C.char)(unsafe.Pointer(C.htmlGetMetaEncoding(doc.DocPtr))))
}

func (doc *Doc) Dump() string {
	return C.GoString(C.DumpXmlToString(doc.DocPtr))
}

func (doc *Doc) DumpHTML() string {
	return C.GoString(C.DumpHtmlToString(doc.DocPtr))
}

func (doc *Doc) RootNode() Node {
	return NewNode(C.xmlDocGetRootElement(doc.DocPtr), doc)
}

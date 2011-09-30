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

xmlChar *
DumpToXmlChar(xmlDoc *doc) {
  xmlChar *buff;
  int buffersize;
  xmlDocDumpFormatMemory(doc, 
                         &buff,
                         &buffersize, 1);
  return buff;
}

xmlChar *
DumpHTMLToXmlChar(xmlDoc *doc) {
  xmlChar *buff;
  int buffersize;
  htmlDocDumpMemory(doc, &buff, &buffersize);
  return buff;
}

xmlNode * GoXmlCastDocToNode(xmlDoc *doc) { return (xmlNode *)doc; }
*/
import "C"

import . "libxml/help"
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
	s := C.htmlGetMetaEncoding(doc.DocPtr)
	return XmlChar2String(s)
}

func (doc *Doc) Dump() string {
	return XmlChar2String(C.DumpToXmlChar(doc.DocPtr))
}

func (doc *Doc) DumpHTML() string {
	return XmlChar2String(C.DumpHTMLToXmlChar(doc.DocPtr))
}

func (doc *Doc) RootNode() Node {
	return NewNode(C.xmlDocGetRootElement(doc.DocPtr), doc)
}

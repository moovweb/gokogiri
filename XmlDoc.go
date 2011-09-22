package libxml 
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

type XmlDoc struct {
	DocPtr *C.xmlDoc
  *XmlNode
}

func buildXmlDoc(ptr *C.xmlDoc) *XmlDoc {
  doc := buildNode(C.GoXmlCastDocToNode(ptr), nil).(*XmlDoc)
	doc.DocPtr = ptr
	return doc
}

func (doc *XmlDoc) Free() { 
  C.xmlFreeDoc(doc.DocPtr) 
}

func (doc *XmlDoc) MetaEncoding() string { 
  s := C.htmlGetMetaEncoding(doc.DocPtr) 
  return XmlChar2String(s)
}

func (doc *XmlDoc) Dump() string {
  return XmlChar2String(C.DumpToXmlChar(doc.DocPtr))
}

func (doc *XmlDoc) DumpHTML() string {
  return XmlChar2String(C.DumpHTMLToXmlChar(doc.DocPtr))
}

func (doc *XmlDoc) RootNode() Node { 
  return buildNode(C.xmlDocGetRootElement(doc.DocPtr), doc)
}

func (doc *XmlDoc) XPathContext() *XPathContext {
  return &XPathContext{Ptr: C.xmlXPathNewContext(doc.DocPtr), Doc: doc}
}

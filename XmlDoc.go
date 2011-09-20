package libxml 
/* 
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

*/ 
import "C"

type XmlDoc struct { 
  Ptr *C.xmlDoc 
}

func BuildXmlDoc(ptr *C.xmlDoc) *XmlDoc {
  if ptr == nil {
    return nil
  }
  return &XmlDoc{Ptr: ptr}
}

func (doc *XmlDoc) Free() { 
  C.xmlFreeDoc(doc.Ptr) 
}

func (doc *XmlDoc) MetaEncoding() string { 
  s := C.htmlGetMetaEncoding(doc.Ptr) 
  return XmlChar2String(s)
}

func (doc *XmlDoc) Dump() string {
  return XmlChar2String(C.DumpToXmlChar(doc.Ptr))
}

func (doc *XmlDoc) DumpHTML() string {
  return XmlChar2String(C.DumpHTMLToXmlChar(doc.Ptr))
}

func (doc *XmlDoc) RootNode() *XmlNode { 
  return BuildXmlNode(C.xmlDocGetRootElement(doc.Ptr), doc)
}

func (doc *XmlDoc) XPathContext() *XPathContext {
  return &XPathContext{Ptr: C.xmlXPathNewContext(doc.Ptr), Doc: doc}
}

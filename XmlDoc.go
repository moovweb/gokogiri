package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
#include <libxml/xpath.h> 
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

func (doc *XmlDoc) RootElement() *XmlNode { 
  return BuildXmlNode(C.xmlDocGetRootElement(doc.Ptr))
}

func (doc *XmlDoc) XPathContext() *XPathContext {
  return &XPathContext{Ptr: C.xmlXPathNewContext(doc.Ptr)}
}

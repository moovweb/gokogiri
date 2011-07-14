package libxml 
/* 
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
*/ 
import "C"

type XPathContext struct { 
  Ptr *C.xmlXPathContext
}

func (context *XPathContext) RegisterNamespace(prefix, href string) bool {
  result := C.xmlXPathRegisterNs(context.Ptr, String2XmlChar(prefix), String2XmlChar(href))
  return result == 0
}

func (context *XPathContext) ParseXPath(input string) {
  
}

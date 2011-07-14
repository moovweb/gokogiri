package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
char* xmlChar2C(xmlChar* x) { return (char *) x; } 
*/ 
import "C" 

func XmlCheckVersion() int { 
  var v C.int 
  C.xmlCheckVersion(v) 
  return int(v) 
} 

func XmlCleanUpParser() { 
  C.xmlCleanupParser() 
}

func HtmlReadFile(url string, encoding string, opts int) *XmlDoc { 
  return BuildXmlDoc(C.htmlReadFile( C.CString(url), C.CString(encoding), C.int(opts) ))
} 

func HtmlReadDoc(content string, url string, encoding string, opts int) *XmlDoc { 
  c := C.xmlCharStrdup( C.CString(content) ) 
  xmlDocPtr := C.htmlReadDoc( c, C.CString(url), C.CString(encoding), C.int(opts) )
  return &XmlDoc{Ptr: xmlDocPtr}
} 

func HtmlReadDocSimple(content string) *XmlDoc {
  return HtmlReadDoc(content, "", "", HTML_PARSE_COMPACT | HTML_PARSE_NOBLANKS | 
                                      HTML_PARSE_NOERROR | HTML_PARSE_NOWARNING)
}

func XmlChar2String(s *C.xmlChar) string {
  return C.GoString( C.xmlChar2C(s) ) 
}
 
func HtmlTagLookup(name string) *C.htmlElemDesc { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  return C.htmlTagLookup(c) 
} 

func HtmlEntityLookup(name string) *C.htmlEntityDesc { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  return C.htmlEntityLookup(c) 
}

func HtmlEntityValueLookup(value uint) *C.htmlEntityDesc { 
  return C.htmlEntityValueLookup( C.uint(value) ) 
}

//Helpers 

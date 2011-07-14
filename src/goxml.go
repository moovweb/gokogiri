package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
char* xmlChar2C(xmlChar* x) { return (char *) x; } 
xmlNode * NodeNext(xmlNode *node) { return node->next; } 
xmlNode * NodeChildren(xmlNode *node) { return node->children; } 
int NodeType(xmlNode *node) { return (int)node->type; } 
*/ 
import "C" 
import ( 
//      "unsafe" 
//      "os" 
) 
const ( 
        //parser option 
        HTML_PARSE_RECOVER = 1 << 0       //relaxed parsing 
        HTML_PARSE_NOERROR = 1 << 5       //suppress error reports 
        HTML_PARSE_NOWARNING = 1 << 6     //suppress warning reports 
        HTML_PARSE_PEDANTIC = 1 << 7      //pedantic error reporting 
        HTML_PARSE_NOBLANKS = 1 << 8      //remove blank nodes 
        HTML_PARSE_NONET = 1 << 11                //forbid network access 
        HTML_PARSE_COMPACT = 1 << 16      //compact small text nodes 
        //element type 
        XML_ELEMENT_NODE = 1 
        XML_ATTRIBUTE_NODE = 2 
        XML_TEXT_NODE = 3 
        XML_CDATA_SECTION_NODE = 4 
        XML_ENTITY_REF_NODE = 5 
        XML_ENTITY_NODE = 6 
        XML_PI_NODE = 7 
        XML_COMMENT_NODE = 8 
        XML_DOCUMENT_NODE = 9 
        XML_DOCUMENT_TYPE_NODE = 10 
        XML_DOCUMENT_FRAG_NODE = 11 
        XML_NOTATION_NODE = 12 
        XML_HTML_DOCUMENT_NODE = 13 
        XML_DTD_NODE = 14 
        XML_ELEMENT_DECL = 15 
        XML_ATTRIBUTE_DECL = 16 
        XML_ENTITY_DECL = 17 
        XML_NAMESPACE_DECL = 18 
        XML_XINCLUDE_START = 19 
        XML_XINCLUDE_END = 20 
        XML_DOCB_DOCUMENT_NODE = 21 
) 
type XmlNode struct { 
  Ptr *C.xmlNode 
}

type XmlDoc struct { 
  Ptr *C.xmlDoc 
}

func XmlCheckVersion() int { 
  var v C.int 
  C.xmlCheckVersion(v) 
  return int(v) 
} 

func XmlCleanUpParser() { 
  C.xmlCleanupParser() 
} 

func (doc *XmlDoc) Free() { 
  C.xmlFreeDoc(doc.Ptr) 
}

func BuildXmlDoc(ptr *C.xmlDoc) *XmlDoc {
  if ptr == nil {
    return nil
  }
  return &XmlDoc{Ptr: ptr}
}
func BuildXmlNode(ptr *C.xmlNode) *XmlNode {
  if ptr == nil {
    return nil
  }
  return &XmlNode{Ptr: ptr}
}

func HtmlReadFile(url string, encoding string, opts int) *XmlDoc { 
  return BuildXmlDoc(C.htmlReadFile( C.CString(url), C.CString(encoding), C.int(opts) ))
} 

func HtmlReadDoc(content string, url string, encoding string, opts int) *XmlDoc { 
  c := C.xmlCharStrdup( C.CString(content) ) 
  xmlDocPtr := C.htmlReadDoc( c, C.CString(url), C.CString(encoding), C.int(opts) )
  return &XmlDoc{Ptr: xmlDocPtr}
} 

func (doc *XmlDoc) GetMetaEncoding() string { 
  s := C.htmlGetMetaEncoding(doc.Ptr) 
  return C.GoString( C.xmlChar2C(s) ) 
} 

func (doc *XmlDoc) GetRootElement() *XmlNode { 
  return BuildXmlNode(C.xmlDocGetRootElement(doc.Ptr))
} 

func (node *XmlNode) GetProp(name string) string { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  s := C.xmlGetProp(node.Ptr, c) 
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
func NewDoc() (doc *C.xmlDoc) { return } 
func NewNode() (node *C.xmlNode) { return } 
func (node *XmlNode) GetNext() *XmlNode { return BuildXmlNode(C.NodeNext(node.Ptr)) } 
func (node *XmlNode) GetChildren() *XmlNode { return BuildXmlNode(C.NodeChildren(node.Ptr)) } 
func (node *XmlNode) GetName() string { return C.GoString( C.xmlChar2C(node.Ptr.name) ) } 
func (node *XmlNode) GetType() int { return int(C.NodeType(node.Ptr)) }
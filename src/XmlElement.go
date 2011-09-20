package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
*/
import "C"

type XmlElement struct { 
  *XmlBaseNode
}

func BuildXmlElement(ptr *C.xmlNode, doc *XmlDoc) *XmlElement {
  if ptr == nil {
    return nil
  }
  return &XmlElement{Ptr: ptr, Doc: doc}
}

func (node *XmlElement) Attribute(name string) string { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  s := C.xmlGetProp(node.Ptr, c) 
  return XmlChar2String(s)
}

func (node *XmlElement) SetAttribute(name string, value string) {
	c_name  := C.xmlCharStrdup( C.CString(name) ) 
	c_value := C.xmlCharStrdup( C.CString(value) ) 
	C.xmlSetProp(node.Ptr, c_name, c_value)
}

func (node *XmlElement) Name() string { 
  return XmlChar2String(node.Ptr.name)
}

func (node *XmlElement) SetName(name string) {
	C.xmlNodeSetName(node.Ptr, C.xmlCharStrdup( C.CString(name) ))
}

func (node *XmlElement) Type() int { 
  return int(C.NodeType(node.Ptr)) 
}
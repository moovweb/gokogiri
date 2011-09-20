package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
*/
import "C"

type XmlElement struct { 
	*XmlNode
}

func (node *XmlElement) Attribute(name string) string { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  s := C.xmlGetProp(node.Ptr(), c) 
  return XmlChar2String(s)
}

func (node *XmlElement) SetAttribute(name string, value string) {
	c_name  := C.xmlCharStrdup( C.CString(name) ) 
	c_value := C.xmlCharStrdup( C.CString(value) ) 
	C.xmlSetProp(node.Ptr(), c_name, c_value)
}
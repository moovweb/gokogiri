package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
*/ 
import "C"


type XmlNode struct { 
  Ptr *C.xmlNode 
}

func BuildXmlNode(ptr *C.xmlNode) *XmlNode {
  if ptr == nil {
    return nil
  }
  return &XmlNode{Ptr: ptr}
}

func (node *XmlNode) GetProp(name string) string { 
  c := C.xmlCharStrdup( C.CString(name) ) 
  s := C.xmlGetProp(node.Ptr, c) 
  return XmlChar2String(s)
}

func (node *XmlNode) Next() *XmlNode { 
  return BuildXmlNode(C.NodeNext(node.Ptr)) 
}

func (node *XmlNode) Children() *XmlNode { 
  return BuildXmlNode(C.NodeChildren(node.Ptr)) 
}

func (node *XmlNode) Name() string { 
  return C.GoString( C.xmlChar2C(node.Ptr.name) ) 
}

func (node *XmlNode) Type() int { 
  return int(C.NodeType(node.Ptr)) 
}
package tree
/*
#cgo pkg-config: libxml-2.0

#include <libxml/tree.h> 
#include <stdlib.h>
*/
import "C"
import "unsafe"
//import "log"

type Attribute struct {
	*XmlNode
}

func NewAttribute(ptr unsafe.Pointer, node Node) *Attribute {
	return NewNode(ptr, node.Doc()).(*Attribute)
}

func (attr *Attribute) Remove() bool {
	if ! attr.IsValid() {
		return false
	}
	attr.Doc().ClearNodeInMap(attr.ptr())
	C.xmlRemoveProp((*C.xmlAttr)(unsafe.Pointer(attr.ptr())))
	attr.NodePtr = nil
	return true
}

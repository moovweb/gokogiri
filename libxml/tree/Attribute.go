package tree
/*
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
	C.xmlRemoveProp((*C.xmlAttr)(unsafe.Pointer(attr.ptr())))
	attr.Doc().ClearNodeInMap(attr.ptr())
	attr.NodePtr = nil
	return true
}

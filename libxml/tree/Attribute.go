package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h> 
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Attribute struct {
	*XmlNode
}

func NewAttribute(ptr unsafe.Pointer, node Node) *Attribute {
	return NewNode(ptr, node.Doc()).(*Attribute);
}
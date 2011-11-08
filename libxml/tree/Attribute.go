package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
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

/*
func (attr *Attribute) Content() string {
    log.Printf("att.First = %v\n", attr.First())
    log.Printf("att.First = %v\n", attr.First())
	return attr.iContent()
}


func (attr *Attribute) SetContent(value string) {
	attr.First().SetContent(value)
}

func (attr *Attribute) String() string {
}
*/

func (attr *Attribute) Remove() bool {
	C.xmlRemoveProp((*C.xmlAttr)(unsafe.Pointer(attr.ptr())))
	return true
}

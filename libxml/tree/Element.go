package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h>
*/
import "C"
import "unsafe"

type Element struct {
	*XmlNode
}

func (node *Element) new(ptr *C.xmlNode) *Element {
	if ptr == nil {
		return nil
	}
	return NewNode(unsafe.Pointer(ptr), node.Doc()).(*Element)
}

func (node *Element) NextElement() *Element {
	return node.new(C.xmlNextElementSibling(node.NodePtr))
}

func (node *Element) PrevElement() *Element {
	return node.new(C.xmlPreviousElementSibling(node.NodePtr))
}

func (node *Element) FirstElement() *Element {
	return node.new(C.xmlFirstElementChild(node.NodePtr))
}

func (node *Element) LastElement() *Element {
	return node.new(C.xmlLastElementChild(node.NodePtr))
}
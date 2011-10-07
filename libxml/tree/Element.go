package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h>
#include "Chelpers.h"
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

func (node *Element) AppendContent(content string) {
    docPtr := (*C.xmlDoc)(node.Doc().Ptr());
    content_p := C.CString(content)
    content_len := len(content)
    C.xmlElement_append(node.ptr(), docPtr, content_p, C.int(content_len), nil) 
}

func (node *Element) PrependContent(content string) {
    docPtr := (*C.xmlDoc)(node.Doc().Ptr());
    content_p := C.CString(content)
    content_len := len(content)
    C.xmlElement_prepend(node.ptr(), docPtr, content_p, C.int(content_len), nil) 
}

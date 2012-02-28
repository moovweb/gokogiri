package xml

//#include "chelper.h"
import "C"
import (
	"unsafe"
	"errors"
)

var (
	fragmentWrapperStart = []byte("<root>")
	fragmentWrapperEnd   = []byte("</root>")
	
	ErrFailParseFragment = errors.New("failed to parse xml fragment")
)

type DocumentFragment struct {
	*Document
	Children []Node
}

const initChildrenNumber = 4

func ParseFragment(document *Document, content, url []byte, options int) (fragment *DocumentFragment, err error) {
	//deal with trivial cases
	if document == nil || len(content) == 0 { return }
	
	content = append(fragmentWrapperStart, content...)
	content = append(content, fragmentWrapperEnd...)

	var contentPtr, urlPtr unsafe.Pointer
	contentPtr   = unsafe.Pointer(&content[0])
	contentLen   := len(content)
	if len(url) > 0  { urlPtr = unsafe.Pointer(&url[0]) }
	
	rootElementPtr := C.xmlParseFragment(document.DocPtr, contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
	
	//
	if rootElementPtr == nil { err = ErrFailParseFragment; return }
	
	fragment = &DocumentFragment{}
	fragment.Document = document
	fragment.Children = make([]Node, 0, initChildrenNumber)
	
	c := (*C.xmlNode)(unsafe.Pointer(rootElementPtr.children))
	var nextSibling *C.xmlNode
	
	for ; c != nil; c = nextSibling {
		nextSibling = (*C.xmlNode)(unsafe.Pointer(c.next))
		C.xmlUnlinkNode(c)
		fragment.Children = append(fragment.Children, NewNode(c, document))
	}
	//now we have rip all its children nodes, we should release the root node
	C.xmlFreeNode(rootElementPtr)
	return
}

func (f *DocumentFragment) Free() {
	for _, node := range(f.Children) {
		node.Free()
	}
}

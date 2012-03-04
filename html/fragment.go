package html

//#include "helper.h"
import "C"
import (
	"unsafe"
	"os"
	"gokogiri/xml"
)

var (
	fragmentWrapperStart = []byte("<html><body>")
	ErrFailParseFragment = os.NewError("failed to parse html fragment")
)

type DocumentFragment struct {
	*Document
	Children []xml.Node
}

const DefaultDocumentFragmentEncoding = "utf-8"
const initChildrenNumber = 4

var defaultDocumentFragmentEncodingBytes = []byte(DefaultDocumentFragmentEncoding)
var emptyDocContent = []byte("")

func ParseFragment(document *Document, content, url []byte, options int) (fragment *DocumentFragment, err os.Error) {
	//deal with trivial cases
	if len(content) == 0 { return }
	
	if document == nil {
		document, err = Parse(emptyDocContent, url, defaultDocumentFragmentEncodingBytes, options)
		if err != nil {
			return
		}
	} 
	
	content = append(fragmentWrapperStart, content...)

	var contentPtr, urlPtr unsafe.Pointer
	contentPtr   = unsafe.Pointer(&content[0])
	contentLen   := len(content)
	if len(url) > 0  { urlPtr = unsafe.Pointer(&url[0]) }
	
	docPtr := unsafe.Pointer(document.DocPtr)
	rootElementPtr := C.htmlParseFragment(docPtr, contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
	
	//
	if rootElementPtr == nil { err = ErrFailParseFragment; return }
	
	fragment = &DocumentFragment{}
	fragment.Document = document
	fragment.Children = make([]xml.Node, 0, initChildrenNumber)
	
	c := (*C.xmlNode)(unsafe.Pointer(rootElementPtr.children))
	var nextSibling *C.xmlNode
	
	for ; c != nil; c = nextSibling {
		nextSibling = (*C.xmlNode)(unsafe.Pointer(c.next))
		C.xmlUnlinkNode(c)
		fragment.Children = append(fragment.Children, xml.NewNode(c, document))
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

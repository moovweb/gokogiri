package xml

//#include "helper.h"
import "C"
import (
	"unsafe"
	"os"
)

type DocumentFragment struct {
	Node
	Children *NodeSet
}

var (
	fragmentWrapperStart = []byte("<root>")
	fragmentWrapperEnd   = []byte("</root>")
)

var ErrFailParseFragment = os.NewError("failed to parse xml fragment")

const initChildrenNumber = 4

func ParseFragment(document Document, content, encoding, url []byte, options int) (fragment *DocumentFragment, err os.Error) {
	//deal with trivial cases
	if len(content) == 0 { return }
	
	//if a document is not provided, we should create an empty Xml document
	//a fragment must reside in a document
	if document == nil {
		document = CreateEmptyDocument(encoding)
	}
	
	//wrap the content before parsing
	content = append(fragmentWrapperStart, content...)
	content = append(content, fragmentWrapperEnd...)

	//set up pointers before calling the C function
	var contentPtr, urlPtr unsafe.Pointer
	contentPtr   = unsafe.Pointer(&content[0])
	contentLen   := len(content)
	if len(url) > 0  { urlPtr = unsafe.Pointer(&url[0]) }
	
	rootElementPtr := C.xmlParseFragment(document.DocPtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
	
	//Note we've parsed the fragment within the given document 
	//the root is not the root of the document; rather it's the root of the subtree from the fragment
	root := NewNode(unsafe.Pointer(rootElementPtr), document)

	//the fragment was in invalid
	if root == nil {
		err = ErrFailParseFragment
		return
	}
	
	fragment = &DocumentFragment{}
	fragment.Node = root
	
	nodes := make([]Node, 0, initChildrenNumber)
	child := root.FirstChild()
	for ; child != nil; child = child.NextSibling() {
		nodes = append(nodes, child)
	}
	fragment.Children = NewNodeSet(document, nodes)
	document.BookkeepFragment(fragment)
	return
}

func (fragment *DocumentFragment) Remove() {
	fragment.Children.Remove()
	fragment.Node.Remove()
}

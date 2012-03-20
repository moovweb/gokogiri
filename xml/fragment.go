package xml

//#include "helper.h"
import "C"
import (
	"unsafe"
	"os"
)

type DocumentFragment struct {
	Node
}

var (
	fragmentWrapperStart = []byte("<root>")
	fragmentWrapperEnd   = []byte("</root>")
)

var ErrFailParseFragment = os.NewError("failed to parse xml fragment")
var ErrEmptyFragment = os.NewError("empty xml fragment")

const initChildrenNumber = 4

func parsefragment(document Document, content, encoding, url []byte, options int) (fragment *DocumentFragment, err os.Error) {
	//wrap the content before parsing
	println("GO : in xml / parsefragment()")
	content = append(fragmentWrapperStart, content...)
	content = append(content, fragmentWrapperEnd...)

	//set up pointers before calling the C function
	var contentPtr, urlPtr unsafe.Pointer
	contentPtr = unsafe.Pointer(&content[0])
	contentLen := len(content)
	if len(url) > 0 {
		url = append(url, 0)
		urlPtr = unsafe.Pointer(&url[0])
	}

	rootElementPtr := C.xmlParseFragmentAsDoc(document.DocPtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)

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
	document.BookkeepFragment(fragment)
	return
}

func ParseFragment(content, inEncoding, url []byte, options int, outEncoding, outBuffer []byte) (fragment *DocumentFragment, err os.Error) {
	document := CreateEmptyDocument(inEncoding, outEncoding, outBuffer)
	fragment, err = parsefragment(document, content, inEncoding, url, options)
	return
}

func (fragment *DocumentFragment) Remove() {
	fragment.Node.Remove()
}

func (fragment *DocumentFragment) Children() []Node {
	nodes := make([]Node, 0, initChildrenNumber)
	child := fragment.FirstChild()
	for ; child != nil; child = child.NextSibling() {
		nodes = append(nodes, child)
	}
	return nodes
}

//just for now
func (fragment *DocumentFragment) String() string {
	out := ""
	nodes := fragment.Children()
	for _, node := range(nodes) {
		out += node.String()
	}
	return out
}

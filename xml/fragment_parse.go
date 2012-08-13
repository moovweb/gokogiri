package xml

//#include "helper.h"
import "C"
import "unsafe"
import . "gokogiri/util"

type FragmentParser interface {
	ParseFragment(*DocCtx, *XmlNode, []byte, []byte, int) (*DocumentFragment, error) 
}

type XmlFragmentParser struct {
}


func (fp *XmlFragmentParser)ParseFragment(docCtx *DocCtx, node *XmlNode, content, url []byte, options int) (fragment *DocumentFragment, err error) {
	//wrap the content before parsing
	content = append(fragmentWrapperStart, content...)
	content = append(content, fragmentWrapperEnd...)

	//set up pointers before calling the C function
	var contentPtr, urlPtr unsafe.Pointer
	contentPtr = unsafe.Pointer(&content[0])
	contentLen := len(content)
	if len(url) > 0 {
		url = AppendCStringTerminator(url)
		urlPtr = unsafe.Pointer(&url[0])
	}

	var rootElementPtr *C.xmlNode

	if node == nil {
		inEncoding := docCtx.InEncoding
		var encodingPtr unsafe.Pointer
		if len(inEncoding) > 0 {
			encodingPtr = unsafe.Pointer(&inEncoding[0])
		}
		rootElementPtr = C.xmlParseFragmentAsDoc(docCtx.DocPtr, contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)

	} else {
		rootElementPtr = C.xmlParseFragment(node.NodePtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
	}

	//Note we've parsed the fragment within the given document 
	//the root is not the root of the document; rather it's the root of the subtree from the fragment
	root := NewNode(unsafe.Pointer(rootElementPtr), docCtx)

	//the fragment was in invalid
	if root == nil {
		err = ErrFailParseFragment
		return
	}

	fragment = &DocumentFragment{}
	fragment.Node = root
	return
}

package html

//#include "helper.h"
import "C"
import (
	"unsafe"
	"os"
	"bytes"
	"gokogiri/xml"
)

var (
	fragmentWrapperStart = []byte("<html><body>")
	ErrFailParseFragment = os.NewError("failed to parse html fragment")
)

type DocumentFragment struct {
	*xml.DocumentFragment
}

const DefaultDocumentFragmentEncoding = "utf-8"
const initChildrenNumber = 4

var defaultDocumentFragmentEncodingBytes = []byte(DefaultDocumentFragmentEncoding)
var bodySigBytes = []byte("<body")

func ParseFragment(document xml.Document, content, url []byte, options int) (fragment *DocumentFragment, err os.Error) {
	//deal with trivial cases
	if len(content) == 0 { return }
	if document == nil {
		document, err = Parse(nil, url, defaultDocumentFragmentEncodingBytes, options)
		if err != nil {
			return
		}
	} 
	
	containBody := (bytes.Index(content, bodySigBytes) >= 0)
	
	content = append(fragmentWrapperStart, content...)

	var contentPtr, urlPtr unsafe.Pointer
	contentPtr   = unsafe.Pointer(&content[0])
	contentLen   := len(content)
	if len(url) > 0  { urlPtr = unsafe.Pointer(&url[0]) }
	
	htmlPtr := C.htmlParseFragment(document.DocPtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
	if htmlPtr == nil {
		err = ErrFailParseFragment
		return
	}

	defer C.xmlFreeNode(htmlPtr)
	
	fragment = &DocumentFragment{}
	fragment.DocumentFragment = &xml.DocumentFragment{}
	fragment.Document = document
	fragment.Children = make([]xml.Node, 0, initChildrenNumber)
	bodyPtr := (*C.xmlNode)(unsafe.Pointer(htmlPtr.children))
	
	if bodyPtr == nil {
		return
	}
	childPtr := bodyPtr
	if ! containBody {
		//note that the body node will be freed together with its parent
		childPtr = (*C.xmlNode)(bodyPtr.children)
	}
	var nextSibling *C.xmlNode
	
	for ; childPtr != nil; {
		nextSibling = (*C.xmlNode)(unsafe.Pointer(childPtr.next))
		C.xmlUnlinkNode(childPtr)
		fragment.Children = append(fragment.Children, xml.NewNode(unsafe.Pointer(childPtr), document))
		childPtr = nextSibling
	}
	return
}

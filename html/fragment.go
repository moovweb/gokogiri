package html

//#include "helper.h"
import "C"
import (
	"unsafe"
	"os"
	"gokogiri/xml"
)

var fragmentWrapperStart = []byte("<div>")
var fragmentWrapperEnd = []byte("</div>")
var bodySigBytes = []byte("<body")

var ErrFailParseFragment = os.NewError("failed to parse html fragment")
var ErrEmptyFragment = os.NewError("empty html fragment")

const initChildrenNumber = 4

func parsefragment(document xml.Document, content, encoding, url []byte, options int) (fragment *xml.DocumentFragment, err os.Error) {	
	//wrap the content
	content = append(fragmentWrapperStart, content...)
	content = append(content, fragmentWrapperEnd...)

	//set up pointers before calling the C function
	var contentPtr, urlPtr unsafe.Pointer
	contentPtr   = unsafe.Pointer(&content[0])
	contentLen   := len(content)
	if len(url) > 0  { urlPtr = unsafe.Pointer(&url[0]) }
	
	rootPtr := C.htmlParseFragment(document.DocPtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
	

	//Note we've parsed the fragment within the given document 
	//the root is not the root of the document; rather it's the root of the subtree from the fragment
	root := xml.NewNode(unsafe.Pointer(rootPtr), document)
	fragment = &xml.DocumentFragment{}
	fragment.Node = root

	document.BookkeepFragment(fragment)
	return
}

func ParseFragment(content, inEncoding, url []byte, options int, outEncoding, outBuffer []byte) (fragment *xml.DocumentFragment, err os.Error) {
	document := CreateEmptyDocument(inEncoding, outEncoding, outBuffer)
	fragment, err = parsefragment(document, content, inEncoding, url, options)
	return
}

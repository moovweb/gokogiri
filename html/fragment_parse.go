package html

//#include "helper.h"
import "C"
import "unsafe"
import "bytes"
import "gokogiri/xml"

type HtmlFragmentParser struct {
}


func (fp *HtmlFragmentParser)ParseFragment(docCtx *xml.DocCtx, node *xml.XmlNode, content, url []byte, options int) (fragment *xml.DocumentFragment, err error) {
	//set up pointers before calling the C function
	var contentPtr, urlPtr unsafe.Pointer
	if len(url) > 0 {
		urlPtr = unsafe.Pointer(&url[0])
	}

	var root xml.Node
	if node == nil {
		containBody := (bytes.Index(content, bodySigBytes) >= 0)

		content = append(fragmentWrapper, content...)
		contentPtr = unsafe.Pointer(&content[0])
		contentLen := len(content)

		inEncoding := docCtx.InEncoding
		var encodingPtr unsafe.Pointer
		if len(inEncoding) > 0 {
			encodingPtr = unsafe.Pointer(&inEncoding[0])
		}
		htmlPtr := C.htmlParseFragmentAsDoc(docCtx.DocPtr, contentPtr, C.int(contentLen), urlPtr, encodingPtr, C.int(options), nil, 0)

		//Note we've parsed the fragment within the given document 
		//the root is not the root of the document; rather it's the root of the subtree from the fragment
		html := xml.NewNode(unsafe.Pointer(htmlPtr), docCtx)

		if html == nil {
			err = ErrFailParseFragment
			return
		}
		root = html

		if !containBody {
			root = html.FirstChild()
			html.AddPreviousSibling(root)
			html.Remove() //remove html otherwise it's leaked
		}
	} else {
		//wrap the content
		newContent := append(fragmentWrapperStart, content...)
		newContent = append(newContent, fragmentWrapperEnd...)
		contentPtr = unsafe.Pointer(&newContent[0])
		contentLen := len(newContent)
		rootElementPtr := C.htmlParseFragment(node.NodePtr(), contentPtr, C.int(contentLen), urlPtr, C.int(options), nil, 0)
		if rootElementPtr == nil {
			//try to parse it as a doc
			fragment, err = fp.ParseFragment(docCtx, nil, content, url, options)
			return
		}
		if rootElementPtr == nil {
			err = ErrFailParseFragment
			return
		}
		root = xml.NewNode(unsafe.Pointer(rootElementPtr), docCtx)
	}

	fragment = &xml.DocumentFragment{}
	fragment.Node = root
	return
}
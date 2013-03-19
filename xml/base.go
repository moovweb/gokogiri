package xml

/*
#include "helper.h"
*/
import "C"
import "gokogiri/xpath"
import . "gokogiri/util"
import "unsafe"


const initialFragments = 2

//the data structure shared by all nodes in a document
type DocCtx struct {
	InEncoding    []byte
	OutEncoding   []byte
	UnlinkedNodes map[*C.xmlNode]bool
	XPathCtx      *xpath.XPath
	DocPtr        unsafe.Pointer
	Fragments     []*DocumentFragment
	FragmentParser
}

func NewDocCtx(docPtr unsafe.Pointer, inEncoding, outEncoding []byte) (*DocCtx) {
	inEncoding = AppendCStringTerminator(inEncoding)
	outEncoding = AppendCStringTerminator(outEncoding)
	ctx := &DocCtx{InEncoding: inEncoding, OutEncoding: outEncoding, DocPtr: docPtr, FragmentParser: &XmlFragmentParser{}}
	ctx.UnlinkedNodes = make(map[*C.xmlNode]bool)
	ctx.XPathCtx = xpath.NewXPath(docPtr)
	ctx.Fragments = make([]*DocumentFragment, 0, initialFragments)
	return ctx
}

func (ctx *DocCtx) SetFragmentParser(fp FragmentParser) {
	ctx.FragmentParser = fp
}

func (ctx *DocCtx) InputEncoding() []byte {
	return ctx.InEncoding
}

func (ctx *DocCtx) OutputEncoding() []byte {
	return ctx.OutEncoding
}

func (ctx *DocCtx) AddUnlinkedNode(nodePtr unsafe.Pointer) {
	p := (*C.xmlNode)(nodePtr)
	ctx.UnlinkedNodes[p] = true
}

func (ctx *DocCtx) RemoveUnlinkedNode(nodePtr unsafe.Pointer) bool {
	p := (*C.xmlNode)(nodePtr)
	if ctx.UnlinkedNodes != nil && ctx.UnlinkedNodes[p] {
		delete(ctx.UnlinkedNodes, p)
		return true
	}
	return false
}

func (ctx *DocCtx) DocType() int {
	return int(C.getNodeType((*C.xmlNode)(ctx.DocPtr)))
}

func (ctx *DocCtx) Free() {
	if ctx.XPathCtx != nil {
		//no need to call free
		//ctx.XPathCtx.Free()
		ctx.XPathCtx = nil
	}
	//must clear the fragments first
	//because the nodes are put in the unlinked list
	if ctx.Fragments != nil {
		for _, frag := range ctx.Fragments {
			frag.Remove()
		}
		ctx.Fragments = nil
	}

	if ctx.UnlinkedNodes != nil {
		for p, _ := range ctx.UnlinkedNodes {
			C.xmlFreeNode(p)
		}
		ctx.UnlinkedNodes = nil
	}

	if ctx.DocPtr != nil {
		C.xmlFreeDoc((*C.xmlDoc)(ctx.DocPtr))
		ctx.DocPtr = nil
	}
}

func (ctx *DocCtx) ParseFragment(input, url []byte, options int) (fragment *DocumentFragment, err error) {
	rootPtr := C.xmlDocGetRootElement((*C.xmlDoc)(ctx.DocPtr))
	if rootPtr == nil {
		fragment, err = ctx.FragmentParser.ParseFragment(ctx, nil, input, url, options)
	} else {
		root := NewNode(unsafe.Pointer(rootPtr), ctx).(*ElementNode)
		fragment, err = ctx.FragmentParser.ParseFragment(ctx, root.XmlNode, input, url, options)
	}
	if fragment != nil {
		ctx.Fragments = append(ctx.Fragments, fragment)
	}
	return
}

func (ctx *DocCtx) CreateElementNode(tag string) (element *ElementNode) {
	tagBytes := GetCString([]byte(tag))
	tagPtr := unsafe.Pointer(&tagBytes[0])
	newNodePtr := C.xmlNewNode(nil, (*C.xmlChar)(tagPtr))
	newNode := NewNode(unsafe.Pointer(newNodePtr), ctx)
	element = newNode.(*ElementNode)
	return
}

func (ctx *DocCtx) CreateTextNode(data string) (text *TextNode) {
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewText((*C.xmlChar)(dataPtr))
	if nodePtr != nil {
		nodePtr.doc = (*_Ctype_struct__xmlDoc)(ctx.DocPtr)
		text = NewNode(unsafe.Pointer(nodePtr), ctx).(*TextNode)
	}
	return
}

func (ctx *DocCtx) CreateCDataNode(data string) (cdata *CDataNode) {
	dataLen := len(data)
	dataBytes := GetCString([]byte(data))
	dataPtr := unsafe.Pointer(&dataBytes[0])
	nodePtr := C.xmlNewCDataBlock((*C.xmlDoc)(ctx.DocPtr), (*C.xmlChar)(dataPtr), C.int(dataLen))
	if nodePtr != nil {
		cdata = NewNode(unsafe.Pointer(nodePtr), ctx).(*CDataNode)
	}
	return
}

func (ctx *DocCtx) Root() (element *ElementNode) {
	nodePtr := C.xmlDocGetRootElement((*C.xmlDoc)(ctx.DocPtr))
	if nodePtr != nil {
		element = NewNode(unsafe.Pointer(nodePtr), ctx).(*ElementNode)
	}
	return
}

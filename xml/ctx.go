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
	
	if ctx.XPathCtx != nil {
		ctx.XPathCtx.Free()
		ctx.XPathCtx = nil
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
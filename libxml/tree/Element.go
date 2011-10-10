package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h>
*/
import "C"
import "unsafe"

type Element struct {
	*XmlNode
}

func (node *Element) new(ptr *C.xmlNode) *Element {
	if ptr == nil {
		return nil
	}
	return NewNode(unsafe.Pointer(ptr), node.Doc()).(*Element)
}

func (node *Element) NextElement() *Element {
	return node.new(C.xmlNextElementSibling(node.NodePtr))
}

func (node *Element) PrevElement() *Element {
	return node.new(C.xmlPreviousElementSibling(node.NodePtr))
}

func (node *Element) FirstElement() *Element {
	return node.new(C.xmlFirstElementChild(node.NodePtr))
}

func (node *Element) LastElement() *Element {
	return node.new(C.xmlLastElementChild(node.NodePtr))
}

func (node *Element) Clear() {
	// Remember, as we delete them, the last one moves to the front
	child := node.First()
	for child != nil {
		child.Remove()
		child.Free()
		child = node.First()
	}
}

func (node *Element) Content() string {
	child := node.First()
	output := ""
	for child != nil {
		output = output + child.String()
		child = child.Next()
	}
	return output
}

func (node *Element) SetContent(content string) {
	node.Clear()
	node.AppendContent(content)
}

func (node *Element) AppendContent(content string) {
	child := node.Doc().ParseFragment(content)
	for child != nil {
		node.AppendChildNode(child)
		child = child.Next()
	}
}

func (node *Element) PrependContent(content string) {
	child := node.Doc().ParseFragment(content).Parent().Last()
	for child != nil {
		node.PrependChildNode(child)
		child = child.Prev()
	}
}

package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/tree.h> 
#include <libxml/xmlstring.h> 
#include <libxml/HTMLtree.h>

int NodeType(xmlNode *node) { return (int)node->type; }

char *
DumpNodeToXmlChar(xmlNode *node, xmlDoc *doc) {
  xmlBuffer *buff = xmlBufferCreate();
  xmlNodeDump(buff, doc, node, 0, 0);
  return (char*)buff->content;
}
*/
import "C"
import "unsafe"

func xmlNodeType(node *C.xmlNode) int {
	return int(C.NodeType(node))
}

func (node *XmlNode) Doc() *Doc {
	return node.DocRef
}
func (node *XmlNode) SetDoc(doc *Doc) {
	node.DocRef = doc
}

func (node *XmlNode) Ptr() unsafe.Pointer {
	return unsafe.Pointer(node.NodePtr)
}

func (node *XmlNode) ptr() *C.xmlNode {
	return node.NodePtr
}

func (node *XmlNode) Type() int {
	return int(C.NodeType(node.ptr()))
}

func (node *XmlNode) Free() {
	C.xmlFreeNode(node.ptr())
}

// Used internally to the XmlNode to quickly create nodes
func (node *XmlNode) new(ptr *_Ctype_struct__xmlNode) Node {
	if ptr == nil {
		return nil
	}
	return NewNode(unsafe.Pointer(ptr), node.Doc())
}

func (node *XmlNode) Parent() Node {
	return node.new(node.ptr().parent)
}

func (node *XmlNode) Next() Node {
	return node.new(node.ptr().next)
}

func (node *XmlNode) Prev() Node {
	return node.new(node.ptr().prev)
}

func (node *XmlNode) First() Node {
	// xmlNode->children actually points to the first 
	// element in an array, so we can just use this as the ptr for first
	return node.new(node.ptr().children)
}

func (node *XmlNode) Last() Node {
	return node.new(node.ptr().last)
}

func (node *XmlNode) Remove() bool {
	C.xmlUnlinkNode(node.ptr())
	return true // TODO: Return false if it was previously unlinked
}

func (node *XmlNode) IsLinked() bool {
	return node.ptr().parent != nil
}

func (node *XmlNode) Duplicate() Node {
	copy := C.xmlCopyNode(node.ptr(), 1)
	return NewNode(unsafe.Pointer(copy), node.Doc())
}

func (node *XmlNode) Size() int {
	return int(C.xmlChildElementCount(node.ptr()))
}

func (node *XmlNode) Name() string {
	return XmlChar2String(node.ptr().name)
}

func (node *XmlNode) SetName(name string) {
	C.xmlNodeSetName(node.ptr(), C.xmlCharStrdup(C.CString(name)))
}

func (node *XmlNode) Content() string {
	return XmlChar2String(C.xmlNodeGetContent(node.ptr()))
}

func (node *XmlNode) SetContent(content string) {
	docPtr := (*C.xmlDoc)(node.Doc().Ptr())
	xmlChar := C.xmlEncodeSpecialChars(docPtr, String2XmlChar(content))
	C.xmlNodeSetContent(node.ptr(), xmlChar)
}

func (node *XmlNode) String() string {
	if node.ptr() == nil {
		return ""
	}
	return C.GoString(C.DumpNodeToXmlChar(node.ptr(), node.Doc().DocPtr))
}

func (node *XmlNode) DumpHTML() string {
	cBuffer := C.xmlBufferCreate()
	C.htmlNodeDump(cBuffer, node.Doc().DocPtr, node.ptr())
	defer C.free(unsafe.Pointer(cBuffer))
	if cBuffer.content == nil {
		return ""
	}
	cString := unsafe.Pointer(cBuffer.content)
	return C.GoString((*C.char)(cString))
}

func (node *XmlNode) Attribute(name string) (*Attribute, bool) {
	cName := String2XmlChar(name)
	xmlAttrPtr := C.xmlHasProp(node.NodePtr, cName)
	didCreate := false
	if xmlAttrPtr == nil {
		didCreate = true
		xmlAttrPtr = C.xmlNewProp(node.NodePtr, cName, C.xmlCharStrdup(C.CString("")))
	}
	attribute := NewNode(unsafe.Pointer(xmlAttrPtr), node.Doc()).(*Attribute)
	return attribute, didCreate
}

func (node *XmlNode) AppendChildNode(child Node) {
	C.xmlAddChild(node.ptr(), (*C.xmlNode)(child.Ptr()))
}
func (node *XmlNode) PrependChildNode(child Node) {
	if node.Size() >= 1 {
		node.First().AddNodeBefore(child)
	} else {
		node.AppendChildNode(child)
	}
}

func (node *XmlNode) AddNodeAfter(sibling Node) {
	C.xmlAddNextSibling(node.ptr(), (*C.xmlNode)(sibling.Ptr()))
}
func (node *XmlNode) AddNodeBefore(sibling Node) {
	C.xmlAddPrevSibling(node.ptr(), (*C.xmlNode)(sibling.Ptr()))
}

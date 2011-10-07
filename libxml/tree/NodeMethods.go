package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/tree.h> 
#include <libxml/xmlstring.h> 

xmlNode * GoNodeNext(xmlNode *node) { return node->next; } 
xmlNode * GoNodePrev(xmlNode *node) { return node->prev; } 
xmlNode * GoNodeChildren(xmlNode *node) { return node->children; }
xmlNode * GoNodeLast(xmlNode *node) { return node->last; } 
xmlNode * GoNodeParent(xmlNode *node) { return node->parent; } 
int NodeType(xmlNode *node) { return (int)node->type; }

char * GoNodeName(xmlNode *node) { return (char*)node->name; }

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

func (node *XmlNode) Ptr() unsafe.Pointer {
	return unsafe.Pointer(node.NodePtr)
}

func (node *XmlNode) ptr() *C.xmlNode {
	return node.NodePtr
}

func (node *XmlNode) Type() int {
	return int(C.NodeType(node.ptr()))
}

func (node *XmlNode) Parent() Node {
	return NewNode(unsafe.Pointer(C.GoNodeParent(node.ptr())), node.Doc())
}

func (node *XmlNode) Next() Node {
	return NewNode(unsafe.Pointer(C.GoNodeNext(node.ptr())), node.Doc())
}

func (node *XmlNode) Prev() Node {
	return NewNode(unsafe.Pointer(C.GoNodePrev(node.ptr())), node.Doc())
}

func (node *XmlNode) First() Node {
	return NewNode(unsafe.Pointer(C.GoNodeChildren(node.ptr())), node.Doc())
}

func (node *XmlNode) Last() Node {
	return NewNode(unsafe.Pointer(C.GoNodeLast(node.ptr())), node.Doc())
}

func (node *XmlNode) Remove() {
	C.xmlUnlinkNode(node.ptr())
}

func (node *XmlNode) Name() string {
	return C.GoString(C.GoNodeName(node.ptr()))
}

func (node *XmlNode) Size() int {
	return int(C.xmlChildElementCount(node.ptr()))
}

func (node *XmlNode) SetName(name string) {
	C.xmlNodeSetName(node.ptr(), C.xmlCharStrdup(C.CString(name)))
}

func (node *XmlNode) String() string {
	return C.GoString(C.DumpNodeToXmlChar(node.ptr(), node.Doc().DocPtr))
}

func (node *XmlNode) Attribute(name string) (*Attribute, bool) {
	cName := C.xmlCharStrdup(C.CString(name))
	xmlAttrPtr := C.xmlHasProp(node.NodePtr, cName)
	didCreate := false
	if xmlAttrPtr == nil {
		didCreate = true;
		xmlAttrPtr = C.xmlNewProp(node.NodePtr, cName, C.xmlCharStrdup(C.CString("")))
	}
	attribute := NewNode(unsafe.Pointer(xmlAttrPtr), node.Doc()).(*Attribute)
	return attribute, didCreate;
}

func (node *XmlNode) AttributeValue(name string) string {
	c := C.xmlCharStrdup(C.CString(name))
	s := C.xmlGetProp(node.ptr(), c)
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

func (node *XmlNode) SetAttributeValue(name string, value string) {
	c_name := C.xmlCharStrdup(C.CString(name))
	c_value := C.xmlCharStrdup(C.CString(value))
	C.xmlSetProp(node.ptr(), c_name, c_value)
}

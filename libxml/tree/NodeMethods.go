package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
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

func (node *XmlNode) Ptr() *C.xmlNode {
	return node.NodePtr
}

/*
	In order to properly cast the pointer to a C.xmlDoc later, we must 
	provide this "anonymous" pointer so that the external function can 
	use a Type Assertion on it.
*/
func (node *XmlNode) AnonPtr() interface{} {
	return node.Ptr()
}

func (node *XmlNode) Type() int {
	return int(C.NodeType(node.Ptr()))
}

func (node *XmlNode) Parent() Node {
	return NewNode(C.GoNodeParent(node.Ptr()), node.Doc())
}

func (node *XmlNode) Next() Node {
	return NewNode(C.GoNodeNext(node.Ptr()), node.Doc())
}

func (node *XmlNode) Prev() Node {
	return NewNode(C.GoNodePrev(node.Ptr()), node.Doc())
}

func (node *XmlNode) First() Node {
	return NewNode(C.GoNodeChildren(node.Ptr()), node.Doc())
}

func (node *XmlNode) Last() Node {
	return NewNode(C.GoNodeLast(node.Ptr()), node.Doc())
}

func (node *XmlNode) Remove() {
	C.xmlUnlinkNode(node.Ptr())
}

func (node *XmlNode) Name() string {
	return C.GoString(C.GoNodeName(node.Ptr()))
}

func (node *XmlNode) Size() int {
	return int(C.xmlChildElementCount(node.Ptr()))
}

func (node *XmlNode) SetName(name string) {
	C.xmlNodeSetName(node.Ptr(), C.xmlCharStrdup(C.CString(name)))
}

func (node *XmlNode) Dump() string {
	return C.GoString(C.DumpNodeToXmlChar(node.Ptr(), node.Doc().DocPtr))
}

func (node *XmlNode) AttributeValue(name string) string {
	c := C.xmlCharStrdup(C.CString(name))
	s := C.xmlGetProp(node.Ptr(), c)
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

func (node *XmlNode) SetAttributeValue(name string, value string) {
	c_name := C.xmlCharStrdup(C.CString(name))
	c_value := C.xmlCharStrdup(C.CString(value))
	C.xmlSetProp(node.Ptr(), c_name, c_value)
}

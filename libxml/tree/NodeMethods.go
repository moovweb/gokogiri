package tree
/*
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/tree.h> 
#include <libxml/xmlstring.h> 

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

func (node *XmlNode) Size() int {
	return int(C.xmlChildElementCount(node.ptr()))
}

func (node *XmlNode) Name() string {
	return C.GoString(C.GoNodeName(node.ptr()))
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

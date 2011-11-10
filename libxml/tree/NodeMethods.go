package tree
/*
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/tree.h> 
#include <libxml/xmlstring.h> 
#include <libxml/HTMLtree.h>

int NodeType(xmlNode *node) { return (int)node->type; }

xmlBufferPtr DumpNodeToXml(xmlNode *node, xmlDoc *doc) {
  xmlBuffer *buff = xmlBufferCreate();
  xmlNodeDump(buff, doc, node, 0, 0);
  return buff;
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
  content := C.xmlNodeGetContent(node.ptr())
  contentPtr := (*C.char)(unsafe.Pointer(content))
  defer XmlFreeChars(contentPtr)
	return XmlChar2String(content)
}

func (node *XmlNode) SetContent(content string) {
	docPtr := (*C.xmlDoc)(node.Doc().Ptr())
  contentXmlChar := String2XmlChar(content)
  defer XmlFreeXmlChars(contentXmlChar)
	xmlChar := C.xmlEncodeSpecialChars(docPtr, contentXmlChar)
  defer XmlFreeXmlChars(xmlChar)
	C.xmlNodeSetContent(node.ptr(), xmlChar)
}

func (node *XmlNode) String() string {
	if node.ptr() == nil {
		return ""
	}
  buffer := C.DumpNodeToXml(node.ptr(), node.Doc().DocPtr)
  defer C.xmlBufferFree(buffer)
  bufferPtr := (*C.char)(unsafe.Pointer(buffer.content))
	return C.GoString(bufferPtr)
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
  childPtr := (*C.xmlNode)(child.Ptr())
  C.xmlUnlinkNode(childPtr);
  copiedChild := C.xmlDocCopyNode(childPtr, node.Doc().DocPtr, 1);
  C.xmlAddChild(node.ptr(), copiedChild);
  C.xmlFreeNode(childPtr); //this is a must; otherwise it would leak memory on text nodes
}
func (node *XmlNode) PrependChildNode(child Node) {
	if node.Size() >= 1 {
		node.First().AddNodeBefore(child)
	} else {
		node.AppendChildNode(child)
	}
}

func (node *XmlNode) AddNodeAfter(sibling Node) {
  siblingPtr := (*C.xmlNode)(sibling.Ptr())
  C.xmlUnlinkNode(siblingPtr);
  copiedSibling := C.xmlDocCopyNode(siblingPtr, node.Doc().DocPtr, 1);
	C.xmlAddNextSibling(node.ptr(), copiedSibling)
  C.xmlFreeNode(siblingPtr)
}

func (node *XmlNode) AddNodeBefore(sibling Node) {
  siblingPtr := (*C.xmlNode)(sibling.Ptr())
  C.xmlUnlinkNode(siblingPtr);
  copiedSibling := C.xmlDocCopyNode(siblingPtr, node.Doc().DocPtr, 1);
	C.xmlAddPrevSibling(node.ptr(), copiedSibling)
  C.xmlFreeNode(siblingPtr)
}

func (node *XmlNode) NewChild(elementName, content string) *Element {
	newCNode := C.xmlNewChild(node.ptr(), nil, String2XmlChar(elementName), String2XmlChar(content))
	return NewNode(unsafe.Pointer(newCNode), node.Doc()).(*Element)
}

func (node *XmlNode) Wrap(elementName string) (wrapperNode *Element) {
	// Build the wrapper
	wrapperNode = node.Parent().NewChild(elementName, "")
	// Add it after me
	node.AddNodeAfter(wrapperNode)
	// Add me as its child
	wrapperNode.AppendChildNode(node)
	println("ABOUT TO SEGFAULT")
	println(wrapperNode.String())
	println("See, never made it!")
	return
}

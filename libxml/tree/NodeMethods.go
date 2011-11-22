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

//export invalidNode
func invalidNode(nodePtr unsafe.Pointer, docPtr unsafe.Pointer) {
	doc := (*Doc)(docPtr)
	node := (*C.xmlNode)(nodePtr)
	doc.ClearNodeInMap(node)
}

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
	if ! node.IsValid() {
		return NIL_NODE
	}
	return int(C.NodeType(node.ptr()))
}

func (node *XmlNode) Free() {
	if node.IsValid() {
		child := node.First()
		for child != nil {
			child.Free()
			child = node.First()
		}
		node.Remove()
		C.xmlFreeNode(node.ptr())
		node.Doc().ClearNodeInMap(node.ptr())
		node.NodePtr = nil
	}
}

func (node *XmlNode) IsValid() bool {
	_, xmlNodeInMap := node.Doc().LookupNodeInMap(node.NodePtr)
    return (node.NodePtr != nil && xmlNodeInMap == node)
}

// Used internally to the XmlNode to quickly create nodes
func (node *XmlNode) new(ptr *_Ctype_struct__xmlNode) Node {
	if ptr == nil {
		return nil
	}
	return NewNode(unsafe.Pointer(ptr), node.Doc())
}

func (node *XmlNode) Parent() Node {
	if ! node.IsValid() {
		return nil
	}
	return node.new(node.ptr().parent)
}

func (node *XmlNode) Next() Node {
	if ! node.IsValid() {
		return nil
	}
	return node.new(node.ptr().next)
}

func (node *XmlNode) Prev() Node {
	if ! node.IsValid() {
		return nil
	}
	return node.new(node.ptr().prev)
}

func (node *XmlNode) First() Node {
	// xmlNode->children actually points to the first 
	// element in an array, so we can just use this as the ptr for first
	if ! node.IsValid() {
		return nil
	}
	return node.new(node.ptr().children)
}

func (node *XmlNode) Last() Node {
	if ! node.IsValid() {
		return nil
	}
	return node.new(node.ptr().last)
}

func (node *XmlNode) Remove() bool {
	if ! node.IsValid() {
		return false
	}
	if ! node.IsLinked() {
		return false
	}
	
	C.xmlUnlinkNode(node.ptr())
	return true // TODO: Return false if it was previously unlinked
}

func (node *XmlNode) IsLinked() bool {
	if ! node.IsValid() {
		return false
	}
	return node.ptr().parent != nil
}

func (node *XmlNode) Duplicate() Node {
	if ! node.IsValid() {
		return nil
	}
	copy := C.xmlCopyNode(node.ptr(), 1)
	return NewNode(unsafe.Pointer(copy), node.Doc())
}

func (node *XmlNode) Size() int {
	if ! node.IsValid() {
		return 0
	}
	return int(C.xmlChildElementCount(node.ptr()))
}

func (node *XmlNode) Name() string {
	if ! node.IsValid() {
		return ""
	}
	return XmlChar2String(node.ptr().name)
}

func (node *XmlNode) SetName(name string) {
	if node.IsValid() {
		nameXmlCharPtr := String2XmlChar(name)
		defer XmlFreeChars(unsafe.Pointer(nameXmlCharPtr))
		C.xmlNodeSetName(node.ptr(), nameXmlCharPtr)
	}
}

func (node *XmlNode) Content() string {
	if ! node.IsValid() {
		return ""
	}
	contentXmlCharPtr := C.xmlNodeGetContent(node.ptr())
	defer XmlFreeChars(unsafe.Pointer(contentXmlCharPtr))
	return XmlChar2String(contentXmlCharPtr)
}

func (node *XmlNode) encodeSpecialChars(content string) *C.xmlChar {
	contentXmlCharPtr := String2XmlChar(content)
	defer XmlFreeChars(unsafe.Pointer(contentXmlCharPtr))
	docPtr := (*C.xmlDoc)(node.Doc().Ptr())
	encodedXmlCharPtr := C.xmlEncodeSpecialChars(docPtr, contentXmlCharPtr)
	return encodedXmlCharPtr
}

func (node *XmlNode) SetCDataContent(content string) {
	if node.IsValid() {
    	encodedXmlCharPtr := node.encodeSpecialChars(content)
    	defer XmlFreeChars(unsafe.Pointer(encodedXmlCharPtr))
		C.xmlNodeSetContent(node.ptr(), encodedXmlCharPtr)
	}
}

// This is overriden in some subclasses... by default use the CData content method
func (node *XmlNode) SetContent(content string) {
	if node.IsValid() {
    	node.SetCDataContent(content)
	}
}

func (node *XmlNode) String() string {
	if ! node.IsValid() {
		return ""
	}
	buffer := C.DumpNodeToXml(node.ptr(), node.Doc().DocPtr)
	defer C.xmlBufferFree(buffer)
	contentCharPtr := (*C.char)(unsafe.Pointer(buffer.content))
	return C.GoString(contentCharPtr)
}

func (node *XmlNode) DumpHTML() string {
	if ! node.IsValid() {
		return ""
	}
	buffer := C.xmlBufferCreate()
	C.htmlNodeDump(buffer, node.Doc().DocPtr, node.ptr())
	defer C.xmlBufferFree(buffer)
	contentCharPtr := (*C.char)(unsafe.Pointer(buffer.content))
	return C.GoString(contentCharPtr)
}

func (node *XmlNode) Attribute(name string) (*Attribute, bool) {
	if ! node.IsValid() {
		return nil, false
	}
	nameXmlCharPtr := String2XmlChar(name)
	defer XmlFreeChars(unsafe.Pointer(nameXmlCharPtr))
	xmlAttrPtr := C.xmlHasProp(node.NodePtr, nameXmlCharPtr)
	didCreate := false
	if xmlAttrPtr == nil {
		didCreate = true
		emptyCharPtr := C.CString("")
		defer C.free(unsafe.Pointer(emptyCharPtr))
		emptyXmlCharPtr := C.xmlCharStrdup(emptyCharPtr)
		defer XmlFreeChars(unsafe.Pointer(emptyXmlCharPtr))
		xmlAttrPtr = C.xmlNewProp(node.NodePtr, nameXmlCharPtr, emptyXmlCharPtr)
	}
	attribute := NewNode(unsafe.Pointer(xmlAttrPtr), node.Doc()).(*Attribute)
	return attribute, didCreate
}

func (node *XmlNode) AppendChildNode(child Node) {
	if node.IsValid() && child.IsValid() {
		childPtr := (*C.xmlNode)(child.Ptr())
		C.xmlUnlinkNode(childPtr)
		if node.Doc().DocPtr != child.Doc().DocPtr {
			copiedChildPtr := C.xmlDocCopyNode(childPtr, node.Doc().DocPtr, 1)
			C.xmlAddChild(node.ptr(), copiedChildPtr)
			child.Free() //this is a must; otherwise it would leak memory on text nodes
		} else {
			C.xmlAddChild(node.ptr(), childPtr)
		}
	}
}

func (node *XmlNode) PrependChildNode(child Node) {
	if node.IsValid() && child.IsValid() {
		if node.Size() >= 1 {
			node.First().AddNodeBefore(child)
		} else {
			node.AppendChildNode(child)
		}
	}
}

func (node *XmlNode) AddNodeAfter(sibling Node) {
	if node.IsValid() && sibling.IsValid() {
		siblingPtr := (*C.xmlNode)(sibling.Ptr())
		C.xmlUnlinkNode(siblingPtr)
		if node.Doc().DocPtr != sibling.Doc().DocPtr {
			copiedSibling := C.xmlDocCopyNode(siblingPtr, node.Doc().DocPtr, 1)
			C.xmlAddNextSibling(node.ptr(), copiedSibling)
			sibling.Free()
		} else {
			C.xmlAddNextSibling(node.ptr(), siblingPtr)
		}
	}
}

func (node *XmlNode) AddNodeBefore(sibling Node) {
	if node.IsValid() && sibling.IsValid() {
		siblingPtr := (*C.xmlNode)(sibling.Ptr())
		C.xmlUnlinkNode(siblingPtr)
		if node.Doc().DocPtr != sibling.Doc().DocPtr {
			copiedSibling := C.xmlDocCopyNode(siblingPtr, node.Doc().DocPtr, 1)
			C.xmlAddPrevSibling(node.ptr(), copiedSibling)
			sibling.Free()
		} else {
			C.xmlAddPrevSibling(node.ptr(), siblingPtr)
		}
	}
}

// If you get nil back from NewChild, make sure that the element type can have children
func (node *XmlNode) NewChild(elementName, content string) *Element {
	if ! node.IsValid() {
		return nil
	}
	nameXmlCharPtr := String2XmlChar(elementName)
	defer XmlFreeChars(unsafe.Pointer(nameXmlCharPtr))
	contentXmlCharPtr := String2XmlChar(content)
	defer XmlFreeChars(unsafe.Pointer(contentXmlCharPtr))
	newCNode := C.xmlNewChild(node.ptr(), nil, nameXmlCharPtr, contentXmlCharPtr)
	return NewNode(unsafe.Pointer(newCNode), node.Doc()).(*Element)
}

func (node *XmlNode) Wrap(elementName string) (wrapperNode *Element) {
	if ! node.IsValid() {
		return nil
	}
	// Build the wrapper
	wrapperNode = node.Doc().NewElement(elementName)
	// Add it after me
	node.AddNodeBefore(wrapperNode)
	// Add me as its child
	node.Remove()
	wrapperNode.AppendChildNode(node)
	return
}
/*
func (node *XmlNode) Children() (children []Node) {
	children = make([]*XmlNode, node.Size())
	currentNode := node.First()
	for i := 0; currentNode != nil; i++ {
		children[i] = currentNode
		currentNode = currentNode.Next()
	}
}*/

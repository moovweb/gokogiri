package tree
/* 
#include <libxml/xmlversion.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
#include <libxml/xpath.h> 

char *
DumpXmlToString(xmlDoc *doc) {
  xmlChar *buff;
  int buffersize;
  xmlDocDumpFormatMemory(doc, 
                         &buff,
                         &buffersize, 1);
  return (char *)buff;
}

char *
DumpHtmlToString(xmlDoc *doc) {
  xmlChar *buff;
  int buffersize;
  htmlDocDumpMemory(doc, &buff, &buffersize);
  return (char *)buff;
}
*/
import "C"
import "unsafe"
import "fmt"

type PtrPair struct {
	node Node
	xmlNode *XmlNode
}

type Doc struct {
	DocPtr *C.xmlDoc
	*XmlNode
	nodeMap map[string]*PtrPair
}

func NewDoc(ptr unsafe.Pointer) *Doc {
	doc := NewNode(ptr, nil).(*Doc)
	doc.DocPtr = (*C.xmlDoc)(ptr)
	doc.InitDocNodeMap()
	return doc
}

func CreateHtmlDoc() *Doc {
	cDoc := C.htmlNewDoc(String2XmlChar(""), String2XmlChar(""))
	return NewNode(unsafe.Pointer(cDoc), nil).(*Doc)
}

func (doc *Doc) NewElement(name string) *Element {
	nameXmlCharPtr := String2XmlChar(name)
	defer XmlFreeChars(unsafe.Pointer(nameXmlCharPtr))
	xmlNode := C.xmlNewNode(nil, nameXmlCharPtr)
	node := NewNode(unsafe.Pointer(xmlNode), doc)
	return node.(*Element)
}

func (doc *Doc) InitDocNodeMap() {
	if doc.nodeMap == nil {
		doc.nodeMap = make(map[string]*PtrPair)
	}
}

func (doc *Doc) LookupNodeInMap(nodePtr *C.xmlNode) (node Node, xmlNode *XmlNode) {
	id := fmt.Sprintf("%d", unsafe.Pointer(nodePtr))
	pair := doc.nodeMap[id]
	if pair == nil {
		return nil, nil
	} else {
		return pair.node, pair.xmlNode
	}
	return
}

func (doc *Doc) SaveNodeInMap(nodePtr *C.xmlNode, node Node, xmlNode *XmlNode) {
	id := fmt.Sprintf("%d", unsafe.Pointer(nodePtr))
	pair := &PtrPair{node: node, xmlNode: xmlNode}
	doc.nodeMap[id] = pair, true
}

func (doc *Doc) ClearNodeInMap(nodePtr *C.xmlNode) {
	id := fmt.Sprintf("%d", unsafe.Pointer(nodePtr))
	doc.nodeMap[id] = nil, false
}

func (doc *Doc) Free() {
	if doc.DocPtr != nil {
		C.xmlFreeDoc(doc.DocPtr)
		doc.DocPtr = nil
	}
}

func (doc *Doc) IsValid() bool {
	return (doc.DocPtr != nil)
}

func (doc *Doc) MetaEncoding() string {
	if ! doc.IsValid() {
		return ""
	}
	metaEncodingXmlCharPtr := C.htmlGetMetaEncoding(doc.DocPtr)
	return C.GoString((*C.char)(unsafe.Pointer(metaEncodingXmlCharPtr)))
}

func (doc *Doc) String() string {
	if ! doc.IsValid() {
		return ""
	}
	// TODO: Decide what type of return to do HTML or XML
	dumpCharPtr := C.DumpXmlToString(doc.DocPtr)
	defer XmlFreeChars(unsafe.Pointer(dumpCharPtr))
	return C.GoString(dumpCharPtr)
}

func (doc *Doc) DumpHTML() string {
	if ! doc.IsValid() {
		return ""
	}
	dumpCharPtr := C.DumpHtmlToString(doc.DocPtr)
	defer XmlFreeChars(unsafe.Pointer(dumpCharPtr))
	return C.GoString(dumpCharPtr)
}

func (doc *Doc) DumpXML() string {
	if ! doc.IsValid() {
		return ""
	}
	dumpCharPtr := C.DumpXmlToString(doc.DocPtr)
	defer XmlFreeChars(unsafe.Pointer(dumpCharPtr))
	return C.GoString(dumpCharPtr)
}

func (doc *Doc) RootElement() *Element {
	if ! doc.IsValid() {
		return nil
	}
	return NewNode(unsafe.Pointer(C.xmlDocGetRootElement(doc.DocPtr)), doc).(*Element)
}

func (doc *Doc) NewCData(content string) *CData {
	if ! doc.IsValid() {
		return nil
	}
	length := C.int(len([]byte(content)))
	cData := C.xmlNewCDataBlock(doc.DocPtr, String2XmlChar(content), length)
	return NewNode(unsafe.Pointer(cData), doc).(*CData)
}

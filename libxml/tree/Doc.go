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

type Doc struct {
	DocPtr *C.xmlDoc
	*XmlNode
	nodeMap map[string]*XmlNode
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
	return NewNode(unsafe.Pointer(C.xmlNewNode(nil, nameXmlCharPtr)), doc).(*Element)
}

func (doc *Doc) InitDocNodeMap() {
	if doc.nodeMap == nil {
		doc.nodeMap = make(map[string]*XmlNode)
	}
}

func (doc *Doc) LookupNode(nodePtr *C.xmlNode) (node *XmlNode) {
	id := fmt.Sprintf("%d", unsafe.Pointer(nodePtr))
	node = doc.nodeMap[id]
	if node == nil {
		node = &XmlNode{NodePtr: nodePtr, DocRef: doc}
		doc.nodeMap[id] = node
	}
	return
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

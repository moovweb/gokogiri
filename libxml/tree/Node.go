package tree
/*
#include <libxml/tree.h>
*/
import "C"
import "unsafe"
//import "runtime"
//import "fmt"

type Node interface {
	ptr() *C.xmlNode
	Ptr() unsafe.Pointer // Used to access the C.Ptr's externally
	Doc() *Doc           // reference to doc
	SetDoc(doc *Doc)     // set the reference to the doc go object

	String() string
	DumpHTML() string
	Remove() bool
	Duplicate() Node // Copy this node
	Free()           // free the libxml structs. BE CAREFUL!
	IsLinked() bool

	// Element Traversal Methods
	Parent() Node // child->parent link
	First() Node  // first child link
	Last() Node   // last child link
	Next() Node   // next sibling link
	Prev() Node   // previous sibling link

	// Informational Methods
	Size() int
	Type() int

	// Node Getters and Setters
	Name() string
	SetName(name string)

	//IsValid checks if the node has been freed with libxml
	IsValid() bool
	
	//Path
	Path() string

	Content() string
	SetContent(content string)
	SetCDataContent(content string)

	NewChild(elementName, content string) *Element
	Wrap(elementName string) *Element

	AppendChildNode(child Node)
	PrependChildNode(child Node)
	AddNodeAfter(sibling Node)
	AddNodeBefore(sibling Node)

	Attribute(name string) (*Attribute, bool) // First, the attribute, then if it is new or not
}

type XmlNode struct {
	NodePtr *C.xmlNode
	DocRef  *Doc
}

func GetNodeType(node Node) (nodeType int) {
	switch v := node.(type) {
	case *Doc:
	    nodeType = XML_DOCUMENT_NODE
	case *Element:
	    nodeType = XML_ELEMENT_NODE
	case *Attribute:
	    nodeType = XML_ATTRIBUTE_NODE
	case *CData:
		nodeType = XML_CDATA_SECTION_NODE
	case *Text:
		nodeType = XML_TEXT_NODE
	default:
	    nodeType = NIL_NODE
	}
	return
}

func NewNode(ptr unsafe.Pointer, doc *Doc) (node Node) {
	cPtr := (*C.xmlNode)(ptr)
	if cPtr == nil {
		return nil
	}
	if doc == nil {
		doc = &Doc{}
		doc.InitDocNodeMap()
	}
	node, _ = doc.LookupNodeInMap(cPtr)
	nodeType := xmlNodeType(cPtr)
	if node == nil || GetNodeType(node) != nodeType {
		xmlNode := &XmlNode{NodePtr: cPtr, DocRef: doc}
		if nodeType == C.XML_DOCUMENT_NODE || nodeType == C.XML_HTML_DOCUMENT_NODE {
			doc.XmlNode = xmlNode
			// If we are a doc, then we reference ourselves
			doc.XmlNode.DocRef = doc
			doc.DocRef = doc
			doc.DocPtr = (*C.xmlDoc)(ptr)
			node = doc
		} else if nodeType == C.XML_ELEMENT_NODE {
			node = &Element{XmlNode: xmlNode}
		} else if nodeType == C.XML_ATTRIBUTE_NODE {
			node = &Attribute{XmlNode: xmlNode}
		} else if nodeType == C.XML_CDATA_SECTION_NODE {
			node = &CData{XmlNode: xmlNode}
		} else if nodeType == C.XML_TEXT_NODE {
			node = &Text{XmlNode: xmlNode}
		} else {
			node = xmlNode
		}
		doc.SaveNodeInMap(cPtr, node, xmlNode)
		/*
		println("create node. pointer:", ptr, " type:", nodeType, node)
		if nodeType == 3  {
			fmt.Printf("node: %q\n", node.String())
			//if len(node.String()) == 0 {
				pc, file, line, ok := runtime.Caller(0)
				msg := ""
				for i := 1; ok == true; i++ {
					msg += fmt.Sprintf("        from %s:%d:in `<%s>'\n", file, line, runtime.FuncForPC(pc).Name())
					pc, file, line, ok = runtime.Caller(i)
				}
				println("\n\n", msg, "\n")
			//}
		}
	} else {
		println("cached node. pointer:", ptr, " type:", nodeType, node)
		*/
	}
	return
}

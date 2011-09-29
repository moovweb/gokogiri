package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>
#include <stdlib.h> 

xmlNode ** 
NodeSetArray(xmlNodeSet *nodes) { 
  if(nodes != NULL) {
    return (xmlNode **)nodes->nodeTab; }
  return NULL; }

xmlNode * 
FetchNode(xmlNodeSet *nodes, int index) { 
  return nodes->nodeTab[index]; }

int
SizeOfSet(xmlNodeSet *set) {
  return set->nodeNr;
}

*/
import "C"

type NodeSet struct {
	Ptr *C.xmlNodeSet
	Doc *Doc
}

func NewNodeSet(aPtr interface{}, doc *Doc) *NodeSet {
	ptr := aPtr.(*C.xmlNodeSet)
	if ptr == nil {
		return nil
	}
	return &NodeSet{Ptr: ptr, Doc: doc}
}

func (nodeSet *NodeSet) Size() int {
	return int(C.SizeOfSet(nodeSet.Ptr))
}

func (nodeSet *NodeSet) NodeAt(pos int) Node {
	node := C.FetchNode(nodeSet.Ptr, C.int(pos))
	return NewNode(node, nodeSet.Doc)
}

func (nodeSet *NodeSet) First() Node {
	return nodeSet.NodeAt(0)
}

func (nodeSet *NodeSet) Slice() []Node {
	list := make([]Node, 0, 0)

	for i := 0; i < nodeSet.Size(); i++ {
		node := nodeSet.NodeAt(i)
		if node != nil {
			list = append(list, node)
		}
	}

	return list
}

func (nodeSet *NodeSet) RemoveAll() {
	for i := 0; i < nodeSet.Size(); i++ {
		node := nodeSet.NodeAt(i)
		if node != nil {
			node.Remove()
		}
	}
}

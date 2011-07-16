package libxml 
/* 
#include <libxml/xpath.h> 
#include <libxml/xpathInternals.h>


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
import("C")

type XmlNodeSet struct {
  Ptr *C.xmlNodeSet
  Doc *XmlDoc
}

func BuildXmlNodeSet(ptr *C.xmlNodeSet, doc *XmlDoc) *XmlNodeSet {
  if ptr == nil {
    return nil
  }
  return &XmlNodeSet{Ptr: ptr, Doc: doc}
}

func (nodeSet *XmlNodeSet) Size() int {
  return int(C.SizeOfSet(nodeSet.Ptr))
}

func (nodeSet *XmlNodeSet) NodeAt(pos int) *XmlNode {
  node := C.FetchNode(nodeSet.Ptr, C.int(pos))
  return BuildXmlNode(node, nodeSet.Doc)
}

func (nodeSet *XmlNodeSet) Slice() []XmlNode {
  list := make([]XmlNode, 0, 100)

  for i := 0; i < nodeSet.Size(); i++ {
    node := nodeSet.NodeAt(i);
    if node != nil {
      list = append(list, *node)
    }
  }

  return list
}
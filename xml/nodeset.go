package xml

type NodeSet struct {
	Document
	Nodes []Node
}

func NewNodeSet(document Document, ptrs []unsafe.Pointer) (set *NodeSet) {
	set = &NodeSet{}
	set.Document = document
	if num := len(ptrs); num > 0 {
		set.Nodes = make([]Node, num)
		for i, p := range(ptrs) {
			set.Nodes[i] = NewNode(p, document)
		}
	}
	return
}
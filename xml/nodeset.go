package xml

import "unsafe"

type NodeSet struct {
	Nodes []Node
}

func NewNodeSet(docCtx *DocCtx, nodes interface{}) (set *NodeSet) {
	switch t := nodes.(type) {
	case []Node:
		set.Nodes = t
	case []unsafe.Pointer:
		if num := len(t); num > 0 {
			set.Nodes = make([]Node, num)
			for i, p := range(t) {
				set.Nodes[i] = NewNode(p, docCtx)
			}
		}
	default:
		//unexpected param type
		//ignore the data
	}
	return
}

func (set *NodeSet) Length() int {
	return len(set.Nodes)
}

func (set *NodeSet) Remove() {
	for _, node := range(set.Nodes) {
		node.Remove()
	}
	set.Nodes = set.Nodes[:0]
}

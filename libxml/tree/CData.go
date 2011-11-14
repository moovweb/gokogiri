package tree
/* 
#include <libxml/tree.h>
*/
import "C"

type CData struct {
	*XmlNode
}

func (node *CData) DumpHTML() string {
	return node.String()
}

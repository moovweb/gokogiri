package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h>
*/
import "C"

type CData struct {
	*XmlNode
}

func (node *CData) DumpHTML() string {
	return node.String()
}
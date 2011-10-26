package tree
/* 
#cgo LDFLAGS: -lxml2
#cgo CFLAGS: -I/usr/include/libxml2
#include <libxml/tree.h>
*/
import "C"

type TextNode struct {
	*XmlNode
}

func (node *TextNode) Content() string {
	return XmlChar2String(node.ptr().content)
}

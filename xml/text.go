package xml

/*
#include "helper.h"
#include <string.h>
const xmlChar xmlStringTextNoenc[] = { 't', 'e', 'x', 't', 'n', 'o', 'e', 'n', 'c', 0 };
const xmlChar* disableEscaping = xmlStringTextNoenc;
*/
import "C"

type TextNode struct {
	*XmlNode
}

// DisableOutputEscaping disables the usual safeguards against creating invalid XML and allows the
// characters '<', '>', and '&' to be written out verbatim. Normally they are safely escaped as entities.
//
// This API is intended to provide support for XSLT processors and similar XML manipulation libraries that
// may need to output unsupported entity references or use the XML API for non-XML output. It should never
// be used in the normal course of XML processing.
func (node *TextNode) DisableOutputEscaping() {
	node.Ptr.name = C.disableEscaping
}

package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
*/
import "C"

type XmlElement struct { 
	*XmlNode
}
package libxml 
/* 
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/HTMLparser.h> 
#include <libxml/HTMLtree.h> 
#include <libxml/xmlstring.h> 
xmlNode * NodeNext(xmlNode *node) { return node->next; } 
xmlNode * NodeChildren(xmlNode *node) { return node->children; } 
int NodeType(xmlNode *node) { return (int)node->type; }
*/ 
import "C"

type XmlElement struct { 
  XmlNode
}


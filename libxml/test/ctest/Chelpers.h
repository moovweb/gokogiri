#include <libxml/parser.h>
#include <libxml/tree.h>

extern int xmlElement_append(xmlNodePtr node, xmlDocPtr doc, const char* content, int content_length, const char* encoding);
extern int xmlElement_prepend(xmlNodePtr node, xmlDocPtr doc, const char* content, int content_length, const char* encoding);

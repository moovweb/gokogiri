#ifndef __PARSE_H__
#define __PARSE_H__

#include <libxml/tree.h>
#include <libxml/parser.h>
#include <libxml/HTMLtree.h>
#include <libxml/HTMLparser.h>
#include <libxml/xmlsave.h>

xmlDoc* native_parse(void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int errror_buffer_len);
xmlNode* native_parse_fragment(xmlDoc* doc, void *buffer, int buffer_len, void *url, int options, void *error_buffer, int error_buffer_len);
void native_save_document(void *obj, void *node, void *encoding, int options);

xmlDoc* newEmptyXmlDoc();
xmlElementType getNodeType(xmlNode *node);
char *xmlDocDumpToString(xmlDoc *doc, void *encoding, int format);
char *htmlDocDumpToString(xmlDoc *doc, int format);
void xmlFreeChars(char *buffer);

#endif //__PARSE_H__
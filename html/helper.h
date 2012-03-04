#ifndef __CHELPER_H__
#define __CHELPER_H__

#include <libxml/tree.h>
#include <libxml/parser.h>
#include <libxml/HTMLtree.h>
#include <libxml/HTMLparser.h>
#include <libxml/xmlsave.h>

htmlDocPtr html_native_parse(void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int errror_buffer_len);
xmlNode* html_native_parse_fragment(xmlDoc* doc, void *buffer, int buffer_len, void *url, int options, void *error_buffer, int error_buffer_len);

#endif //__CHELPER_H__
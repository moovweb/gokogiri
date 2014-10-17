#include "helper.h"
#include "../xml/helper.h"
#include <string.h>

void* htmlParse(void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int error_buffer_len, dispatchTable* libxml_symbols) {
	const char *c_buffer       = (char*)buffer;
	const char *c_url          = (char*)url;
	const char *c_encoding     = (char*)encoding;
	shim_xmlDoc *doc = NULL;

	libxml_symbols->xmlResetLastError();
	doc = libxml_symbols->htmlReadMemory(c_buffer, buffer_len, c_url, c_encoding, options);

	return doc;
}

void* htmlParseFragment(void *doc, void *buffer, int buffer_len, void *url, int options, void *error_buffer, int error_buffer_len, dispatchTable* libxml_symbols) {
	shim_xmlNode* root_element = NULL;
	shim_xmlParserErrors errCode;
	errCode = libxml_symbols->xmlParseInNodeContext((xmlNodePtr)doc, buffer, buffer_len, options, root_element);
	if (errCode != XML_ERR_OK) {
		return NULL;
	}
	return root_element;
}

void* htmlParseFragmentAsDoc(void *doc, void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int error_buffer_len, dispatchTable* libxml_symbols) {
	shim_xmlDoc* tmpDoc = NULL;
	shim_xmlNode* tmpRoot = NULL;
	tmpDoc = libxml_symbols->htmlReadMemory((char*)buffer, buffer_len, (char*)url, (char*)encoding, options);
	if (tmpDoc == NULL) {
		return NULL;
	}
	tmpRoot = libxml_symbols->xmlDocGetRootElement(tmpDoc);
	if (tmpRoot == NULL) {
		return NULL;
	}
	tmpRoot = libxml_symbols->xmlDocCopyNode(tmpRoot, doc, 1);
	libxml_symbols->xmlFreeDoc(tmpDoc);
	return tmpRoot;
}

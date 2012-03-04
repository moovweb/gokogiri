#include <string.h>
#include "helper.h"

//internal callback functions
int xml_write_callback(void *document, char *buffer, int len) {
	xmlNodeWriteCallback(document, buffer, len);
  	return len;
}

int close_callback(void * ctx) {
  return 0;
}

xmlDoc* newEmptyXmlDoc() {
	//why does xmlNewDoc NOT call xmlInitParser like other parse functions?
	xmlInitParser();
	return xmlNewDoc(BAD_CAST XML_DEFAULT_VERSION); 
}

xmlElementType getNodeType(xmlNode *node) { return node->type; }

void xmlFreeChars(char *buffer) { 
	if (buffer) {
		xmlFree((xmlChar*)buffer); 
	}
}

char *xmlDocDumpToString(xmlDoc *doc, void *encoding, int format) {
	xmlChar *buff;
	int buffersize;
	xmlDocDumpFormatMemoryEnc(doc, &buff, &buffersize, (char*)encoding, format);
	return (char*)buff;
}

char *htmlDocDumpToString(htmlDocPtr doc, int format) {
	xmlChar *buff;
	int buffersize;
	htmlDocDumpMemoryFormat(doc, &buff, &buffersize, format);
	return (char*)buff;
}

xmlDoc* xmlParse(void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int error_buffer_len) {
	const char *c_buffer       = (char*)buffer;
	const char *c_url          = (char*)url;
	const char *c_encoding     = (char*)encoding;
	xmlDoc *doc = NULL;
	
	xmlResetLastError();
	doc = xmlReadMemory(c_buffer, buffer_len, c_url, c_encoding, options);

	if(doc == NULL) {
		xmlErrorPtr error;
	    xmlFreeDoc(doc);
	    error = xmlGetLastError();
		if(error != NULL && error_buffer != NULL && error->level >= XML_ERR_ERROR) {
			char *c_error_buffer = (char*)error_buffer;
			if (error->message != NULL) {
				strncpy(c_error_buffer, error->message, error_buffer_len-1);
				c_error_buffer[error_buffer_len-1] = '\0';
			}
			else {
				snprintf(c_error_buffer, error_buffer_len, "xml parsing error:%d", error->code);
			}
		}
	}
	return doc;
}

xmlNode* xmlParseFragment(void *doc, void *buffer, int buffer_len, void *url, int options, void *error_buffer, int error_buffer_len) {
	xmlNodePtr root_element = NULL;
	xmlParserErrors errCode;
	
	errCode = xmlParseInNodeContext((xmlNodePtr)doc, buffer, buffer_len, options, &root_element);
	if (errCode != XML_ERR_OK) {
		char *c_error_buffer = (char*)error_buffer;
		snprintf(c_error_buffer, error_buffer_len, "xml fragemnt parsing error (xmlParserErrors):%d", errCode);
		return NULL;
	} 
	return root_element;
}

void xmlSetContent(void *n, void *content) {
	xmlNode *node = (xmlNode*)n;
	xmlNode *child = node->children;
	xmlNode *next = NULL;
	char *encoded = xmlEncodeSpecialChars(node->doc, content);
	if (encoded) {
		while (child) {
			next = child->next ;
			xmlUnlinkNode(child);
			xmlFreeNode(child);
			child = next ;
	  	}
	  	xmlNodeSetContent(node, (xmlChar*)encoded);
		xmlFree(encoded);
	}
}

void xmlSaveNode(void *obj, void *node, void *encoding, int options) {
	xmlSaveCtxtPtr savectx;
	const char *c_encoding = (char*)encoding;
	
	savectx = xmlSaveToIO(
	      (xmlOutputWriteCallback)xml_write_callback,
	      (xmlOutputCloseCallback)close_callback,
	      (void *)obj,
	      encoding,
	      options
	  );
	xmlSaveTree(savectx, (xmlNode*)node);
	xmlSaveClose(savectx);
}
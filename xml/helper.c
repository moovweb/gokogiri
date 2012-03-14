#include <string.h>
#include "helper.h"

//internal callback functions
int xml_write_callback(void *ctx, char *buffer, int len) {
	XmlBufferContext *xmlBufferCtx = (XmlBufferContext*)ctx;
	if (len > 0 && xmlBufferCtx != NULL && xmlBufferCtx->buffer_len > xmlBufferCtx->data_size) {
		int bytesToWrite = xmlBufferCtx->buffer_len - xmlBufferCtx->data_size;
		char *start = NULL;

		if (bytesToWrite > len) {
			bytesToWrite = len;
		}
		start =  xmlBufferCtx->buffer + sizeof(char)*xmlBufferCtx->data_size;
		strncpy(start, buffer, bytesToWrite);
		xmlBufferCtx->data_size = xmlBufferCtx->data_size + bytesToWrite;
		return bytesToWrite;
	}
  	return 0;
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
	printf("buffer_len = %d\n", buffer_len);
	errCode = xmlParseInNodeContext((xmlNodePtr)doc, buffer, buffer_len, options, &root_element);
	if (errCode != XML_ERR_OK) {
		if (error_buffer != NULL && error_buffer_len > 0) {
			char *c_error_buffer = (char*)error_buffer;
			snprintf(c_error_buffer, error_buffer_len, "xml fragemnt parsing error (xmlParserErrors):%d", errCode);
		}
		printf("errorcode %d %d\n", errCode, root_element);
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

int xmlSaveNode(void *buffer, int buffer_len, void *node, void *encoding, int options) {
	xmlSaveCtxtPtr savectx;
	const char *c_encoding = (char*)encoding;
	XmlBufferContext xmlBufferCtx;
	int ret;

	xmlBufferCtx.buffer = (char*)buffer;
	xmlBufferCtx.buffer_len = buffer_len;
	xmlBufferCtx.data_size = 0;
	
	savectx = xmlSaveToIO(
	      (xmlOutputWriteCallback)xml_write_callback,
	      (xmlOutputCloseCallback)close_callback,
	      (void *)&xmlBufferCtx,
	      encoding,
	      options
	  );
	xmlSaveTree(savectx, (xmlNode*)node);
	ret = xmlSaveClose(savectx);
	if (ret <  0) {
		return ret;
	}
	return xmlBufferCtx.data_size;
}


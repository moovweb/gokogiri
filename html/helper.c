#include "helper.h"
#include "../xml/helper.h"
#include <string.h>

htmlDocPtr htmlParse(void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int error_buffer_len) {
	const char *c_buffer       = (char*)buffer;
	const char *c_url          = (char*)url;
	const char *c_encoding     = (char*)encoding;
	xmlDoc *doc = NULL;
	
	xmlResetLastError();
	doc = htmlReadMemory(c_buffer, buffer_len, c_url, c_encoding, options);

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

xmlNode* htmlParseFragment(void *doc, void *buffer, int buffer_len, void *url, int options, void *error_buffer, int error_buffer_len) {
	return xmlParseFragment((xmlDoc*)doc, buffer, buffer_len, url, options, error_buffer, error_buffer_len);
}

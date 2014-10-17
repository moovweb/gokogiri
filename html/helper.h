#ifndef __CHELPER_H__
#define __CHELPER_H__

#include "../libxml_shim.h"

void* htmlParse(void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int errror_buffer_len, dispatchTable* libxml_symbols);
void* htmlParseFragment(void* doc, void *buffer, int buffer_len, void *url, int options, void *error_buffer, int error_buffer_len, dispatchTable* libxml_symbols);
void* htmlParseFragmentAsDoc(void *doc, void *buffer, int buffer_len, void *url, void *encoding, int options, void *error_buffer, int error_buffer_len, dispatchTable* libxml_symbols);

#endif //__CHELPER_H__

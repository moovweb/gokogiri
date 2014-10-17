#ifndef LIBXML_SHIM
#define LIBXML_SHIM

#include "libxml_defs.h"

// wrapper functions
shim_htmlDoc* shim_htmlNewDoc(const shim_xmlChar* URI, const shim_xmlChar* ExternalID, dispatchTable* libxml_symbols);
shim_xmlChar* shim_htmlGetMetaEncoding(shim_htmlDoc* doc, dispatchTable* libxml_symbols);
int shim_htmlSetMetaEncoding(shim_htmlDoc* doc, const shim_xmlChar* encoding, dispatchTable* libxml_symbols);


//package html
typedef void xmlResetLastError_fnPtr();
typedef struct shim_xmlDoc* htmlReadMemory_fnPtr(const char* buffer, int size, const char* URL, const char* encoding, int options);
typedef struct shim_htmlDoc* htmlNewDoc_fnPtr(const shim_xmlChar* URI, const shim_xmlChar* ExternalID);
typedef struct shim_xmlChar* htmlGetMetaEncoding_fnPtr(struct shim_htmlDoc* doc);
typedef int htmlSetMetaEncoding_fnPtr(struct shim_htmlDoc* doc, const shim_xmlChar* encoding);
typedef struct shim_xmlNode* htmlParseFragmentAsDoc_fnPtr(void* doc, void* buffer, int buffer_len, void* url, void* encoding, int options, void* error_buffer, int error_buffer_len);
typedef struct shim_xmlNode* xmlDocCopyNode_fnPtr(const struct shim_xmlNode* node, struct shim_xmlDoc* doc, int extended);

//package xml
// typedef void* xmlResetLastError_fnPtr();
// typedef shim_xmlDoc* xmlReadMemory_fnPtr(const char* buffer, int size, const char* URL, const char* encoding, int options);
typedef void* xmlFreeDoc_fnPtr(struct shim_xmlDoc* doc);
// typedef shim_xmlError* xmlGetLastError_fnPtr();
// typedef shim_xmlDoc* xmlReadFile_fnPtr(const char* filename, const char* encoding, int options);
typedef struct shim_xmlDoc* xmlNewDoc_fnPtr(const shim_xmlChar* version);
typedef struct shim_xmlNode* xmlDocGetRootElement_fnPtr(struct shim_xmlDoc* doc);
// typedef shim_xmlNode* xmlGetID_fnPtr(shim_xmlDoc* doc, shim_xmlChar* data);
typedef struct shim_xmlNode* xmlNewNode_fnPtr(void* ns, const shim_xmlChar* name);
// typedef shim_xmlNode* xmlNewText_fnPtr(const shim_xmlChar* content);
// typedef shim_xmlNode* xmlNewCDataBlock(shim_xmlDoc* doc, const shim_xmlChar* content, int len);
// typedef shim_xmlNode* xmlNewComment_fnPtr(const shim_xmlChar* content);
// typedef shim_xmlNode* xmlNewDocPI_fnPtr(shim_xmlDoc* doc, const shim_xmlChar* name, const shim_xmlChar* content)
// typedef shim_xmlEntity* xmlGetDocEntity_fnPtr(shim_xmlDoc* doc, const shim_xmlChar* name);
// typedef void xmlFreeNode_fnPtr(shim_xmlNode* node);
// typedef void xmlFreeDoc_fnPtr(shim_xmlDoc* doc);
typedef shim_xmlParserErrors xmlParseInNodeContext_fnPtr(xmlNodePtr node, const char *data, int datalen, int options, xmlNodePtr lst);


typedef struct dispatchTable {

  //package html
  xmlResetLastError_fnPtr* xmlResetLastError;
  htmlReadMemory_fnPtr* htmlReadMemory;
  htmlNewDoc_fnPtr* htmlNewDoc;
  htmlGetMetaEncoding_fnPtr* htmlGetMetaEncoding;
  htmlSetMetaEncoding_fnPtr* htmlSetMetaEncoding;
  htmlParseFragmentAsDoc_fnPtr* htmlParseFragmentAsDoc;
  xmlDocCopyNode_fnPtr* xmlDocCopyNode;
  xmlParseInNodeContext_fnPtr* xmlParseInNodeContext;

  // package xml
  // xmlResetLastError_fnPtr xmlResetLastError;
  // xmlReadMemory_fnPtr xmlReadMemory;
  xmlFreeDoc_fnPtr* xmlFreeDoc;
  // xmlGetLastError_fnPtr xmlGetLastError;
  // xmlReadFile_fnPtr xmlReadFile;
  xmlNewDoc_fnPtr* xmlNewDoc;
  xmlDocGetRootElement_fnPtr* xmlDocGetRootElement;
  // xmlGetID_fnPtr xmlGetID;
  xmlNewNode_fnPtr* xmlNewNode;
  // xmlNewText_fnPtr xmlNewText;
  // xmlNewCDataBlock_fnPtr xmlNewCDataBlock;
  // xmlNewComment_fnPtr xmlNewComment;
  // xmlNewDocPI_fnPtr xmlNewDocPI;
  // xmlGetDocEntity_fnPtr xmlGetDocEntity;
  // xmlFreeNode_fnPtr xmlFreeNode;
  // xmlFreeDoc_fnPtr xmlFreeDoc;

} dispatchTable;

dispatchTable* init(const char* filename);
void loadSymbols(void* handle, dispatchTable* symbols);

#endif
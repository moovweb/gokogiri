#ifndef LIBXML_SHIM
#define LIBXML_SHIM

void init(const char* filename);
void loadSymbols(void* handle);

//package html
typedef void* htmlDoc_objPtr;
typedef void xmlResetLastError_fnPtr();
typedef xmlDoc_objPtr* htmlReadMemory_fnPtr(const char* buffer, int size, const char* URL, const char* encoding, int options);
typedef htmlDoc_objPtr* htmlNewDoc(const xmlChar_objPtr* URI, const xmlChar_objPtr* ExternalID);
typedef xmlChar_objPtr* htmlGetMetaEncoding_fnPtr(htmlDoc_objPtr* doc);
typedef int htmlSetMetaEncoding_fnPtr(htmlDoc_objPtr* doc, const xmlChar_objPtr* encoding);
typedef xmlNode_objPtr* htmlParseFragmentAsDoc_fnPtr(void* doc, void* buffer, int buffer_len, void* url, void* encoding, int options, void* error_buffer, int error_buffer_len);
typedef xmlNode_objPtr* xmlDocCopyNode_fnPtr(const xmlNode_objPtr* node, xmlDoc_objPtr* doc, int extended);

//package xml
typedef void* xmlNode_objPtr; 
typedef void* xmlDoc_objPtr;
// typedef void* xmlResetLastError_fnPtr();
// typedef xmlDoc_objPtr* xmlReadMemory_fnPtr(const char* buffer, int size, const char* URL, const char* encoding, int options);
// typedef void* xmlErrorPtr_objPtr;
typedef void* xmlFreeDoc_fnPtr(xmlDoc_objPtr* doc);
// typedef xmlErrorPtr_objPtr* xmlGetLastError_fnPtr();
// typedef xmlDoc_objPtr* xmlReadFile_fnPtr(const char* filename, const char* encoding, int options);
typedef void* xmlChar_objPtr;
// typedef xmlDoc_objPtr* xmlNewDoc_fnPtr(const xmlChar_objPtr* version);
typedef xmlNode_objPtr* xmlDocGetRootElement_fnPtr(xmlDoc_objPtr* doc);
// typedef xmlNode_objPtr* xmlGetID_fnPtr(xmlDoc_objPtr* doc, xmlChar_objPtr* data);
// typedef xmlNode_objPtr* xmlNewNode_fnPtr(void* ns, const xmlChar_objPtr* name);
// typedef xmlNode_objPtr* xmlNewText_fnPtr(const xmlChar_objPtr* content);
// typedef xmlNode_objPtr* xmlNewCDataBlock(xmlDoc_objPtr* doc, const xmlChar_objPtr* content, int len);
// typedef xmlNode_objPtr* xmlNewComment_fnPtr(const xmlChar_objPtr* content);
// typedef xmlNode_objPtr* xmlNewDocPI_fnPtr(xmlDoc_objPtr* doc, const xmlChar_objPtr* name, const xmlChar_objPtr* content)
// typedef void* xmlEntity_objPtr;
// typedef xmlEntity_objPtr* xmlGetDocEntity_fnPtr(xmlDoc_objPtr* doc, const xmlChar_objPtr* name);
// typedef void xmlFreeNode_fnPtr(xmlNode_objPtr* node);
// typedef void xmlFreeDoc_fnPtr(xmlDoc_objPtr* doc);


typedef struct {

  //package html
  htmlDoc_objPtr htmlDoc; //htmlDocPtr
  xmlResetLastError_fnPtr xmlResetLastError;
  htmlReadMemory_fnPtr htmlReadMemory;
  htmlNewDoc_fnPtr htmlNewDoc;
  htmlGetMetaEncoding_fnPtr htmlGetMetaEncoding;
  htmlSetMetaEncoding_fnPtr htmlSetMetaEncoding;
  htmlParseFragmentAsDoc_fnPtr htmlParseFragmentAsDoc;
  xmlDocCopyNode_fnPtr xmlDocCopyNode;

  // package xml
  xmlNode_objPtr xmlNode;
  xmlDoc_objPtr xmlDoc;
  // xmlResetLastError_fnPtr xmlResetLastError;
  // xmlReadMemory_fnPtr xmlReadMemory;
  // xmlError_objPtr xmlErrorPtr;
  xmlFreeDoc_fnPtr xmlFreeDoc;
  // xmlGetLastError_fnPtr xmlGetLastError;
  // xmlReadFile_fnPtr xmlReadFile;
  xmlChar_objPtr xmlChar;
  // xmlNewDoc_fnPtr xmlNewDoc;
  xmlDocGetRootElement_fnPtr xmlDocGetRootElement;
  // xmlGetID_fnPtr xmlGetID;
  // xmlNewNode_fnPtr xmlNewNode;
  // xmlNewText_fnPtr xmlNewText;
  // xmlNewCDataBlock_fnPtr xmlNewCDataBlock;
  // xmlNewComment_fnPtr xmlNewComment;
  // xmlNewDocPI_fnPtr xmlNewDocPI;
  // xmlEntity_objPtr xmlEntity;
  // xmlGetDocEntity_fnPtr xmlGetDocEntity;
  // xmlFreeNode_fnPtr xmlFreeNode;
  // xmlFreeDoc_fnPtr xmlFreeDoc;

} dispatchTable;

#endif
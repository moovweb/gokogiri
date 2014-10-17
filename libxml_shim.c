#include "libxml_shim.h"

dispatchTable* init(const char* version) {
  // Open the dynamic library
  void *handle = nullptr;
  switch (version) {
    case "legacy":
      // for now even though this isn't legacy
      handle = dlopen("~/dev/clibs/output/libxml2/lib/libxml2.2.dylib", RTLD_LAZY);
      break;
    default:
      break;
  }

  if (!handle) {
    fprintf(stderr, "%s\n", dlerror());
    exit(1);
  }
  dispatchTable* libxml_symbols = nullptr;
  loadSymbols(handle, libxml_symbols);

  return libxml_symbols;
}

void loadSymbols(void* handle, dispatchTable* libxml_symbols) {

  libxml_symbols->xmlResetLastError = (xmlResetLastError_fnPtr)dlsym(handle, "xmlResetLastError");
  if (libxml_symbols->xmlResetLastError == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlResetLastError.\n");
    return;
  }
  libxml_symbols->htmlReadMemory = (htmlReadMemory_fnPtr)dlsym(handle, "htmlReadMemory");
  if (libxml_symbols->htmlReadMemory == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->htmlReadMemory.\n");
    return;
  }
  libxml_symbols->htmlNewDoc = (htmlNewDoc_fnPtr)dlsym(handle, "htmlNewDoc");
  if (libxml_symbols->htmlNewDoc == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->htmlNewDoc.\n");
    return;
  }
  libxml_symbols->htmlGetMetaEncoding = (htmlGetMetaEncoding_fnPtr)dlsym(handle, "htmlGetMetaEncoding");
  if (libxml_symbols->htmlGetMetaEncoding == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->htmlGetMetaEncoding.\n");
    return;
  }
  libxml_symbols->htmlSetMetaEncoding = (htmlSetMetaEncoding_fnPtr)dlsym(handle, "htmlSetMetaEncoding");
  if (libxml_symbols->htmlSetMetaEncoding == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->htmlSetMetaEncoding.\n");
    return;
  }
  libxml_symbols->htmlParseFragmentAsDoc = (htmlParseFragmentAsDoc_fnPtr)dlsym(handle, "htmlParseFragmentAsDoc");
  if (libxml_symbols->htmlParseFragmentAsDoc == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->htmlParseFragmentAsDoc.\n");
    return;
  }
  libxml_symbols->xmlDocCopyNode = (xmlDocCopyNode_fnPtr)dlsym(handle, "xmlDocCopyNode");
  if (libxml_symbols->xmlDocCopyNode == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlDocCopyNode.\n");
    return;
  }
  libxml_symbols->xmlFreeDoc = (xmlFreeDoc_fnPtr)dlsym(handle, "xmlFreeDoc");
  if (libxml_symbols->xmlFreeDoc == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlFreeDoc.\n");
    return;
  }
  libxml_symbols->xmlDocGetRootElement = (xmlDocGetRootElement_fnPtr)dlsym(handle, "xmlDocGetRootElement");
  if (libxml_symbols->xmlDocGetRootElement == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlDocGetRootElement.\n");
    return;
  }
  libxml_symbols->xmlNewNode = (xmlNewNode_fnPtr)dlsym(handle, "xmlNewNode");
  if (libxml_symbols->xmlNewNode == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlNewNode.\n");
    return;
  }
  libxml_symbols->xmlNewDoc = (xmlNewDoc_fnPtr)dlsym(handle, "xmlNewDoc");
  if (libxml_symbols->xmlNewDoc == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlNewDoc.\n");
    return;
  }
  libxml_symbols->xmlParseInNodeContext = (xmlParseInNodeContext_fnPtr)dlsym(handle, "xmlParseInNodeContext");
  if (libxml_symbols->xmlParseInNodeContext == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlParseInNodeContext.\n");
    return;
  }
  //etc etc

  //structs
  // shim_htmlDoc = libxml_symbols->htmlNewDoc(NULL, NULL);
  // shim_xmlNode = libxml_symbols->xmlNewNode(NULL, "div");
  // shim_xmlDoc = libxml_symbols->xmlNewDoc(NULL);

  return;
}

// need wrapper functions because cgo can't handle function pointers :(

shim_htmlDoc* shim_htmlNewDoc(const shim_xmlChar* URI, const shim_xmlChar* ExternalID, dispatchTable* libxml_symbols) {
  return libxml_symbols->htmlNewDoc(URI, ExternalID);
}

shim_xmlChar* shim_htmlGetMetaEncoding(shim_htmlDoc* doc, dispatchTable* libxml_symbols) {
  return libxml_symbols->htmlGetMetaEncoding(doc);
}

int shim_htmlSetMetaEncoding(shim_htmlDoc* doc, const shim_xmlChar* encoding, dispatchTable* libxml_symbols) {
  return libxml_symbols->htmlSetMetaEncoding(doc, encoding);
}



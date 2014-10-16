#include "libxml_shim.h"

dispatchTable* libxml_symbols = nullptr;

void init(const char* version) {
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

  loadSymbols(handle);

  return;
}

void loadSymbols(void* handle) {
  //reset it
  libxml_symbols = nullptr;

  libxml_symbols->xmlDoc = (xmlDoc_objPtr)dlsym(handle, "xmlDoc");
  if (libxml_symbols->xmlDoc == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlDoc.\n");
    return;
  }
  libxml_symbols->xmlNode = (xmlNode_objPtr)dlsym(handle, "xmlNode");
  if (libxml_symbols->xmlNode == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlNode.\n");
    return;
  }
  libxml_symbols->htmlDoc = (htmlDoc_objPtr)dlsym(handle, "htmlDoc");
  if (libxml_symbols->htmlDoc == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->htmlDoc.\n");
    return;
  }
  libxml_symbols->xmlChar = (xmlChar_objPtr)dlsym(handle, "xmlChar");
  if (libxml_symbols->xmlChar == nullptr) {
    fprintf(stderr,"Error: couldn't initialize libxml_symbols->xmlChar.\n");
    return;
  }
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
  //etc etc

  return;
}
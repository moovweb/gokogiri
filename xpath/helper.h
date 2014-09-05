#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <libxml/parser.h>

xmlNode* fetchNode(xmlNodeSet *nodeset, int index);

xmlXPathObjectPtr go_resolve_variables(void* ctxt, char* name, char* ns);
int go_can_resolve_function(void* ctxt, char* name, char* ns);
void exec_xpath_function(xmlXPathParserContextPtr ctxt, int nargs);

xmlXPathFunction go_resolve_function(void* ctxt, char* name, char* ns);

void set_var_lookup(xmlXPathContext* c, void* data);

void set_function_lookup(xmlXPathContext* c, void* data);

int getXPathObjectType(xmlXPathObject* o);

xmlXPathObjectPtr xmlXPathEvalWithTimeout(xmlXPathCompExprPtr comp, xmlXPathContextPtr ctx, int timeout);
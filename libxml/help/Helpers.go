package help
/* 
#cgo pkg-config: libxml-2.0
#include <stdio.h>
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/xmlstring.h> 
char* xmlChar2C(xmlChar* x) { return (char *) x; }
xmlChar* C2xmlChar(char* x) { return (xmlChar *) x; }
void printMemoryLeak() { xmlMemDisplay(stdout); }
*/
import "C"

func XmlCheckVersion() int {
	var v C.int
	C.xmlCheckVersion(v)
	return int(v)
}

func XmlInitParser() {
	C.xmlInitParser()
}

func XmlCleanUpParser() {
	C.xmlCleanupParser()
}

func XmlMemoryAllocation() int {
	return (int)(C.xmlMemBlocks())
}

func XmlMemoryLeakReport() {
	C.printMemoryLeak()
}

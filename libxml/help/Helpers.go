package help
/* 
#include <stdio.h>
#include <libxml/xmlversion.h> 
#include <libxml/parser.h> 
#include <libxml/xmlstring.h> 
#include "XmlMem.h"

char* xmlChar2C(xmlChar* x) { return (char *) x; }
xmlChar* C2xmlChar(char* x) { return (xmlChar *) x; }
*/
import "C"
import "unsafe"
import "log"

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

func InitMemFreeCallback() {
	C.initMemFreeCallback()
}

func XmlMemoryLeakReport() {
	C.xmlMemDisplay(C.stdout)
}

//export XmlNodeFreedByLibXml
func XmlNodeFreedByLibXml(ptr unsafe.Pointer) {
	log.Printf("XmlNodeFreedByLibXml called %d", ptr)
	doc.BookkeepNode(ptr, node)
}
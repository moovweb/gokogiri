package mem

/*
#cgo pkg-config: libxml-2.0

#include <libxml/xmlversion.h>
#include "libxml.h"
*/
import "C"

const (
	LIBXML_VERSION         = C.LIBXML_DOTTED_VERSION
	LIBXML_NUMERIC_VERSION = C.LIBXML_VERSION
)

func init() {
	C.xmlCheckVersion(LIBXML_NUMERIC_VERSION)
	C.libxmlGoInit()
}

func AllocSize() int {
	return int(C.libxmlGoAllocSize())
}

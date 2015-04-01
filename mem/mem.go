package mem

/*
#cgo CFLAGS: -I../../../../../../clibs/include/libxml2
#cgo LDFLAGS: -L../../../../../../clibs/lib -lxml2

#include <libxml/xmlversion.h>
#include "libxml.h"
*/
import "C"

const LIBXML_VERSION = C.LIBXML_DOTTED_VERSION

func init() {
	C.libxmlGoInit()
}

func AllocSize() int {
	return int(C.libxmlGoAllocSize())
}

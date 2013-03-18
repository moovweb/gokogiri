package html

//#include "helper.h"
import "C"
import (
	. "gokogiri/util"
	"gokogiri/xml"
)

var fragmentWrapperStart = []byte("<div>")
var fragmentWrapperEnd = []byte("</div>")
var fragmentWrapper = []byte("<html><body>")
var bodySigBytes = []byte("<body")

const initChildrenNumber = 4

func ParseFragment(content, inEncoding, url []byte, options int, outEncoding []byte) (fragment *xml.DocumentFragment, err error) {
	inEncoding = AppendCStringTerminator(inEncoding)
	outEncoding = AppendCStringTerminator(outEncoding)
	document := CreateEmptyDocument(inEncoding, outEncoding)
	fragment, err = document.ParseFragment(content, url, options)
	return
}

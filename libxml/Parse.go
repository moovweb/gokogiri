package libxml
import . "libxml/tree"

func HtmlParseString(content string) *Doc {
	doc := HtmlParseStringWithOptions(content, "", "",
		HTML_PARSE_RECOVER|
			HTML_PARSE_NONET|
			HTML_PARSE_NOERROR|
			HTML_PARSE_NOWARNING)
	if doc == nil {
		return HtmlParseString("<html />")
	}
	return doc
}

func XmlParseString(content string) *Doc {
	return XmlParseWithOptions(content, "", "", 
        XML_PARSE_RECOVER | 
        XML_PARSE_NONET|
        XML_PARSE_NOERROR|
        XML_PARSE_NOWARNING)

}

func XmlParseFragment(content string) *Doc {
	return XmlParseFragmentWithOptions(content, "", "", 
        XML_PARSE_RECOVER | 
        XML_PARSE_NONET|
        XML_PARSE_NOERROR|
        XML_PARSE_NOWARNING)

}



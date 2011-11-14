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
	return XmlParseWithOption(content, "", "", 1)
}

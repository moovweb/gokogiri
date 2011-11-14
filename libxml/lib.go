package libxml

import "libxml/tree"

func HtmlParseString(input string) *tree.Doc {
	return tree.HtmlParseString(input)
}

func XmlParseString(input string) *tree.Doc {
	return tree.XmlParseString(input)
}
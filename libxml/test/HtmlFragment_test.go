package test

import (
	"libxml"
	"testing"
	"libxml/tree"
	"libxml/help"
	"strings"
)

func TestHtmlFragment(t *testing.T) {
	children := ParseHtmlFragment("<body><div></body>")
	Equal(t, children[0].Name(), "body")
	Equal(t, children[0].First().Name(), "div")
	Equal(t, children[0].String(), "<body><div></body>")
}
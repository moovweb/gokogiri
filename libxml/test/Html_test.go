package test

import "libxml/html"
import(
	"testing"
)

func TestSimpleParse(t *testing.T) {
	doc := html.ParseString("<html />")
	if doc.Size() != 1 {
		t.Error("Should be 1 big!")
	}
}
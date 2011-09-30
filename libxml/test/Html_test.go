package test

import "libxml"
import(
	"testing"
)

func TestSimpleParse(t *testing.T) {
	doc := libxml.HtmlParseString("<html><head /><body /></html>")
	if doc.Size() != 1 {
		t.Error("Incorrect size")
	}
	htmlTag := doc.First()
	if htmlTag.Size() != 2 {
		t.Error(htmlTag.Name())
		t.Error("Two tags are inside of <html>")
	}
	
}
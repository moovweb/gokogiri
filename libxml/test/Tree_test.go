package test

import(
	"libxml"
	"testing"
	"libxml/tree"
)

func TestTree(t *testing.T) {
	doc := libxml.XmlParseString("<root><parent><child /><child>Content</child></parent><aunt /><catlady/></root>")
	Equal(t, doc.Size(), 1);

	root := doc.First()
	if root.Name() != "root" {
		t.Error("Should have returned root element")
	}

	Equal(t, root.Size(), 3);

	// If we are on root, and we go "next", we should get
	// nothing, as root has no siblings. Should return nil
	// error
	AssertNil(t, root.Next(), "root next")
	AssertNil(t, root.Prev(), "root prev")
	AssertNil(t, doc.Parent(), "doc parent")
	parent := Assert(t, root.First(), "first is a node").(tree.Node)
	Equal(t, parent.Name(), "parent")

	catLady := Assert(t, root.Last(), "root last node exists").(tree.Node)
	AssertNil(t, catLady.First(), "catlady first")
	AssertNil(t, catLady.Next(), "catlady has no siblings")

	// See if we get <aunt /> for both of these
	// TODO: implement it so that they are ACTUALLY equal to each other.
	Equal(t, parent.Next().String(), catLady.Prev().String())
}

func AssertNil(t *testing.T, value interface{}, what string) {
	if value != nil {
		t.Error(what, "should be nil")
	}
}
func Equal(t *testing.T, value, expected interface{}) {
	if value != expected {
		t.Error("Expected: ", expected, "\nBut got: ", value)
	}
}
func Assert(t *testing.T, value interface{}, what string) interface{} {
	if value == nil {
		t.Error("Assertion failed: ", what)
	}
	return value
}
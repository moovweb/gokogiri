package test

import (
	"libxml"
	"testing"
	"libxml/tree"
	"strings"
)

func TestTree(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><child /><child>Text</child></parent><aunt /><catlady/></root>")
	Equal(t, doc.Size(), 1)
	Equal(t, doc.Content(), "hiText")

	root := doc.First().(*tree.Element)
	if root.Name() != "root" {
		t.Error("Should have returned root element")
	}

	Equal(t, root.Size(), 3)

	// If we are on root, and we go "next", we should get
	// nothing, as root has no siblings. Should return nil
	// error
	AssertNil(t, root.Next(), "root next")
	AssertNil(t, root.Prev(), "root prev")
	AssertNil(t, doc.Parent(), "doc parent")
	rootText := Assert(t, root.First(), "first is a text node").(tree.Node)
	Equal(t, rootText.Content(), "hi")
	parent := Assert(t, root.FirstElement(), "first is a element").(*tree.Element)
	Equal(t, parent.Name(), "parent")

	lastChild := Assert(t, parent.Last(), "parent last").(tree.Node)
	childText := Assert(t, lastChild.First(), "lastChild's first").(tree.Node)
	Equal(t, childText.Content(), "Text")

	catLady := Assert(t, root.Last(), "root last node exists").(tree.Node)
	AssertNil(t, catLady.First(), "catlady first")
	AssertNil(t, catLady.Next(), "catlady has no siblings")

	// See if we get <aunt /> for both of these
	// TODO: implement it so that they are ACTUALLY equal to each other.
	Equal(t, parent.Next().String(), catLady.Prev().String())
}

func TestAddingChildLast(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	childDoc := tree.Parse("<child/>")
	child := childDoc.First()
	doc.RootElement().FirstElement().AppendChildNode(child)
	if !strings.Contains(doc.String(), "<brother/><child/>") {
		t.Error("Should have new last child")
	}
}

func TestAddingChildFirst(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	childDoc := tree.Parse("<child/>")
	child := childDoc.First()
	doc.RootElement().FirstElement().PrependChildNode(child)
	if !strings.Contains(doc.String(), "<child/><brother/>") {
		t.Fail()
	}
}

func TestAddingBefore(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	childDoc := tree.Parse("<child/>")
	child := childDoc.First()
	doc.RootElement().FirstElement().AddNodeBefore(child)
	if !strings.Contains(doc.String(), "<child/><parent") {
		t.Error("Should have new sibling before")
	}
}

func TestAddingAfter(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
	childDoc := tree.Parse("<child/>")
	child := childDoc.First()
	doc.RootElement().FirstElement().AddNodeAfter(child)
	if !strings.Contains(doc.String(), "</parent><child/></root>") {
		t.Error("Should have new sibling after")
	}
}

func TestNodeDuplicate(t *testing.T) {
	doc := libxml.XmlParseString("<root><parent><brother>hi</brother></parent></root>")
	parent := doc.RootElement().FirstElement()
	brother := parent.FirstElement()
	dupBrother := brother.Duplicate()
	dupBrother.First().SetContent("bye")
	parent.AppendChildNode(dupBrother)
	if !strings.Contains(doc.String(), "<brother>hi</brother>") {
		t.Error("Should have original sibling")
	}
	if !strings.Contains(doc.String(), "<brother>bye</brother>") {
		t.Error("Should have new sibling too!")
	}
}

func TestSetContent(t *testing.T) {
	doc := libxml.XmlParseString("<root>hi</root>")
	root := doc.RootElement()
	text := root.First()
	Equal(t, text.Content(), "hi")
	text.SetContent("bye")
	Equal(t, text.Content(), "bye")
	if !strings.Contains(doc.String(), "<root>bye</root>") {
		t.Fail()
	}
	root.SetContent("world")
	if !strings.Contains(doc.String(), "world") {
		t.Fail()
	}
}

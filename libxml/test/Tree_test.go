package test

import (
	"libxml"
	"testing"
	"libxml/tree"
	"libxml/help"
	"strings"
)
/*
func TestMem(t *testing.T) {
    doc := libxml.XmlParseString("<root>hi<parent><child /><child>Text</child></parent><aunt /><catlady/></root>")
    doc.Free()
    help.XmlCleanUpParser()
    if help.XmlMemoryAllocation() != 0 {
      t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
      help.XmlMemoryLeakReport()
    }
}*/

func TestParallelTree(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root>hi<parent><child /><child>Text</child></parent><aunt /><catlady/></root>")
		done <- false

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

		doc.Free()
		done <- true
	}

	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestParallelAddingChildLast(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
		done <- false

		childDoc := tree.Parse("<child/>")
		child := childDoc.First()
		doc.RootElement().FirstElement().AppendChildNode(child)
		if !strings.Contains(doc.String(), "<brother/><child/>") {
			t.Error("Should have new last child")
			t.Error("%q\n", child.String())
		}

		childDoc.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

func TestParallelAddingChildFirst(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root>hi<parent><foo/><bar/></parent></root>")
		done <- false

		childDoc := tree.Parse("<child/>")
		child := childDoc.First()
		doc.RootElement().FirstElement().PrependChildNode(child)
		if !strings.Contains(doc.String(), "<child/><foo/>") {
			t.Error("Should have new first child: %q\n", doc.String())
		}
		childDoc.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

func TestParallelAddingBefore(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
		done <- false
		childDoc := tree.Parse("<child/>")
		child := childDoc.First()
		doc.RootElement().FirstElement().AddNodeBefore(child)
		if !strings.Contains(doc.String(), "<child/><parent") {
			t.Error("Should have new sibling before")
		}
		childDoc.Free()
		doc.Free()
		done <- true
	}

	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

func TestParallelAddingAfter(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root>hi<parent><brother/></parent></root>")
		done <- false
		childDoc := tree.Parse("<child/>")
		child := childDoc.First()
		doc.RootElement().FirstElement().AddNodeAfter(child)
		if !strings.Contains(doc.String(), "</parent><child/></root>") {
			t.Error("Should have new sibling after")
		}
		childDoc.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestParallelNodeDuplicate(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root><parent><brother>hi</brother></parent></root>")
		done <- false
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
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

func TestParallelSetContent(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root>hi</root>")
		done <- false
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
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestParallelNodeIsLinked(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root><child /></root>")
		done <- false
		child := doc.RootElement().FirstElement()
		if child.IsLinked() != true {
			t.Error("Children start off linked")
		}
		child.Remove()
		if child.IsLinked() != false {
			t.Error("Children should report being unlinked")
		}
		child.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

func TestParallelPath(t *testing.T) {
	testFunc := func(done chan bool) {
		doc := libxml.XmlParseString("<root><child><a/><b/></child></root>")
		done <- false
		child := doc.RootElement().FirstElement()
		if child.Path() != "/root/child" {
			t.Error("path wrong")
		}
		b := child.First().Next()
		if b.Path() != "/root/child/b" {
			t.Error("path wrong")
		}
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)
	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}

}

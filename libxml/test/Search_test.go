package test

import (
	"libxml"
	"libxml/help"
	"libxml/xpath"
	"testing"
)

func TestSearch(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false
		doc := libxml.HtmlParseString("<html><body><div>Hi<div>Mom</div></div></body></html>")
		xp := xpath.NewXPath(doc)
		divs := xp.Search(doc, "//div")

		// Doctype gets returned as the first child!
		if divs.Size() != 2 {
			t.Errorf("Returned the two divs!: %d", divs.Size())
		}

		div := divs.NodeAt(0)
		if div.Size() != 1 {
			t.Error("Only has one element in it!")
		}
		textChild := div.First()
		if textChild.Name() != "text" {
			t.Error("Should return a text child")
		}

		xp.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestMultiSearch(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false

		doc := libxml.HtmlParseString("<html><body><div>Hi<div>Mom</div></div></body></html>")
		xp := xpath.NewXPath(doc)
		divs := xp.Search(doc, "//div")

		// Doctype gets returned as the first child!
		if divs.Size() != 2 {
			t.Errorf("Returned the two divs!: %d", divs.Size())
		}

		div := divs.NodeAt(0)
		if div.Size() != 1 {
			t.Error("Only has one element in it!")
		}
		textChild := div.First()
		if textChild.Name() != "text" {
			t.Error("Should return a text child")
		}

		body := xp.Search(textChild, "//body")
		Equal(t, body.Size(), 1)
		Equal(t, body.NodeAt(0).Name(), "body")

		xp.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

// What if we remove a node we will soon match?
func TestSearchRemoval(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false

		doc := libxml.XmlParseString("<root><parent><child /></parent></root>")
		xp := xpath.NewXPath(doc)
		root := doc.RootElement()
		nodeSet := xp.Search(root, "//*")
		nodes := nodeSet.Slice()
		for i := range nodes {
			node := nodes[i]
			Equal(t, node.Type(), 1)
		}

		parent := root.FirstElement()
		parent.SetContent("empty")
		if parent.IsLinked() != true {
			t.Error("Parent starts off linked")
		}
		parent.Remove()
		if parent.IsLinked() != false {
			t.Error("Parent should report being unlinked")
		}
		parent.Free()
		xp.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

// What if we remove a node that may have been removed
func TestSearchRemoval2(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false
		doc := libxml.XmlParseString("<root><div class=\"a\"><div class=\"b\"><div class=\"c\"/></div></div></root>")
		xp := xpath.NewXPath(doc)
		root := doc.RootElement()
		nodeSet := xp.Search(root, "//div")
		nodes := nodeSet.Slice()
		for i := range nodes {
			node := nodes[i]
			node.SetContent("")
		}

		xp.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

//what if a search returns a nil pointer?
func TestNilSearch(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false

		doc := libxml.XmlParseString("<root id=\"foo\"><h1></h1></root>")
		xp := xpath.NewXPath(doc)
		nodeSet := xp.Search(doc, "//*[@id = 'foo1']//*")
		nodes := nodeSet.Slice()
		if len(nodes) != 0 {
			t.Error("Should return zero size node set")
		}
		xp.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestSearchComments(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false
		doc := libxml.HtmlParseString("<html><head><!-- COMMENT --><title>Mom</title></html>")
		xp := xpath.NewXPath(doc)
		comments := xp.Search(doc, "//comment()")

		// Doctype gets returned as the first child!
		if comments.Size() != 1 {
			t.Errorf("# of comments!: %d", comments.Size())
		}

		xp.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

func TestSelectParent(t *testing.T) {
	testFunc := func(done chan bool) {
		done <- false

		doc := libxml.XmlParseString("<div id='a'><div id='b'>Im a b</div></div>")
		xp := xpath.NewXPath(doc)
		child1 := xp.Search(doc, ".//*[@id='b']").NodeAt(0)
		child2 := xp.Search(child1, "..").NodeAt(0)
		child3 := xp.Search(child2, "/").NodeAt(0)
		if (child3 != doc) {
			t.Errorf("should get doc\n")
		}

		xp.Free()
		doc.Free()
		done <- true
	}
	runParallel(testFunc, numConcurrentRuns)

	if help.XmlMemoryAllocation() != 0 {
		t.Errorf("Memeory leaks %d!!!", help.XmlMemoryAllocation())
		help.XmlMemoryLeakReport()
	}
}

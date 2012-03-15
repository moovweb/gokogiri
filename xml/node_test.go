package xml

import "testing"
import "fmt"


func TestAddChild(t *testing.T) {

	docAssertion := func (doc *XmlDocument) (string, string, string) {
		expectedDocAfterAdd :=
		`<?xml version="1.0" encoding="utf-8"?>
<foo>
  <bar/>
</foo>
`
		doc.Root().AddChild("<bar></bar>")

		return doc.String(), expectedDocAfterAdd, "output of the xml doc after AddChild does not match"
	}

	nodeAssertion := func (doc *XmlDocument) (string, string, string) {
		expectedNodeAfterAdd :=
		`<foo>
  <bar/>
</foo>`

		return doc.Root().String(), expectedNodeAfterAdd, "the output of the xml root after AddChild does not match"
	}


	RunTest(t, "node", "add_child", nil, docAssertion, nodeAssertion)

}

func TestAddPreviousSibling(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		err := doc.Root().AddPreviousSibling("<bar></bar><cat></cat>")

		if err != nil {
			t.Errorf("Error adding previous sibling:\n%v\n", err.String())
		}
	}


	RunTest(t, "node", "add_previous_sibling", testLogic)
}

func TestAddNextSibling(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		doc.Root().AddNextSibling("<bar></bar><baz></baz>")
	}

	RunTest(t, "node", "add_next_sibling", testLogic)
}

func TestSetContent(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.SetContent("<fun></fun>")
	}

	RunTest(t, "node", "set_content", testLogic)
}

func TestSetChildren(t *testing.T) {
	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		root.SetChildren("<fun></fun>")
	}

	RunTest(t, "node", "set_children", testLogic)
}

func TestReplace(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {	
		root := doc.Root()
		root.Replace("<fun></fun><cool/>")
	}


	rootAssertion := func(doc *XmlDocument) (string, string, string) {
		root := doc.Root()
		return root.String(), "<fun/>", "the output of the xml root does not match"
	}

	RunTest(t, "node", "replace", testLogic, rootAssertion)
}

func TestAttributes(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {	

		root := doc.Root()
		attributes := root.Attributes()

		if len(attributes) != 2 || attributes["myname"].String() != "ff" {
			fmt.Printf("%v, %q\n", attributes, attributes["myname"].String())
			t.Error("root's attributes do not match")
		}

		child := root.FirstChild()
		childAttributes := child.Attributes()

		if len(childAttributes) != 1 || childAttributes["class"].String() != "shine" {
			t.Error("child's attributes do not match")
		}
	}

	RunTest(t, "node", "attributes", testLogic)

}

func TestSearch(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		result, _ := root.Search(".//*[@class]")
		if len(result) != 2 {
			t.Error("search at root does not match")
		}
		result, _ = root.Search("//*[@class]")
		if len(result) != 3 {
			t.Error("search at root does not match")
		}
		result, _ = doc.Search(".//*[@class]")
		if len(result) != 3 {
			t.Error("search at doc does not match")
		}
		result, _ = doc.Search(".//*[@class='shine']")
		if len(result) != 2 {
			t.Error("search with value at doc does not match")
		}
	}

	RunTest(t, "node", "search", testLogic)
}

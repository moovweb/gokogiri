package xml

import "testing"
import "fmt"
import "gokogiri/help"

func TestSetValue(t *testing.T) {
	defer help.CheckXmlMemoryLeaks(t)
	doc, err := Parse([]byte("<foo id=\"a\" myname=\"ff\"><bar class=\"shine\"/></foo>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)
	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}
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
	attributes["myname"].SetValue("new")
	expected :=
		`<foo id="a" myname="new">
  <bar class="shine"/>
</foo>`
	if root.String() != expected {
		println("got:\n", root.String())
		println("expected:\n", expected)
		t.Error("root's new attr do not match")
	}
	attributes["id"].Remove()
	expected =
		`<foo myname="new">
  <bar class="shine"/>
</foo>`

	if root.String() != expected {
		println("got:\n", root.String())
		println("expected:\n", expected)
		t.Error("root's remove attr do not match")
	}
	doc.Free()
}

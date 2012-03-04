package xml

import "testing"
import "fmt"
import "gokogiri/help"

func TestAddChild(t *testing.T) {
	doc, err := Parse([]byte("<foo></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	expectedDoc :=
`<?xml version="1.0" encoding="utf-8"?>
<foo/>
`
	expectedDocAfterAdd :=
`<?xml version="1.0" encoding="utf-8"?>
<foo><bar/></foo>
`
	expectedNodeAfterAdd :=
`<foo><bar/></foo>`

	if err != nil {
		t.Error("Parsing has error:", err)
	}
	if doc.String() != expectedDoc {
		t.Error("the output of the xml doc does not match")
	}
	doc.Root().AddChild("<bar></bar>")
	if doc.String() != expectedDocAfterAdd {
		println(doc.String())
		t.Error("the output of the xml doc after AddChild does not match")
	}
	if doc.Root().String() != expectedNodeAfterAdd {
		t.Error("the output of the xml root after AddChild does not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestAddPreviousSibling(t *testing.T) {
	doc, err := Parse([]byte("<foo></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	expected :=
`<?xml version="1.0" encoding="utf-8"?>
<bar/>
<cat/>
<foo/>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	err = doc.Root().AddPreviousSibling("<bar></bar><cat></cat>")
	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestAddNextSibling(t *testing.T) {
	doc, err := Parse([]byte("<foo></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	expected :=
`<?xml version="1.0" encoding="utf-8"?>
<foo/>
<bar/>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	doc.Root().AddNextSibling("<bar></bar>")
	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestSetContent(t *testing.T) {
	doc, err := Parse([]byte("<foo><bar/></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	expected :=
`<?xml version="1.0" encoding="utf-8"?>
<foo>&lt;fun&gt;&lt;/fun&gt;</foo>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	root := doc.Root()
	root.SetContent("<fun></fun>")
	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}


func TestSetChildren(t *testing.T) {
	doc, err := Parse([]byte("<foo><bar1/><bar2/></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	expected :=
`<?xml version="1.0" encoding="utf-8"?>
<foo><fun/></foo>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	
	root := doc.Root()
	root.SetChildren("<fun></fun>")
	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestReplace(t *testing.T) {
	doc, err := Parse([]byte("<foo><bar/></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	expected :=
`<?xml version="1.0" encoding="utf-8"?>
<fun/>
<cool/>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	root := doc.Root()
	root.Replace("<fun></fun><cool/>")
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match")
	}
	root = doc.Root()
	if root.String() != "<fun/>" {
		t.Error("the output of the xml root does not match")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestAttributes(t *testing.T) {
	doc, err := Parse([]byte("<foo id=\"a\" myname=\"ff\"><bar class=\"shine\"/></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		t.Error("Parsing has error:", err)
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
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestSearch(t *testing.T) {
	doc, err := Parse([]byte("<foo id=\"a\" class=\"shine\"><bar class=\"shine\"/><vic class=\"dim\"></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	
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
	
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

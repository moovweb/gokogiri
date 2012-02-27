package xml

import "testing"

func TestAddChild(t *testing.T) {
	doc, err := Parse([]byte("<foo></foo>"), nil, []byte("utf-8"), DefaultParseOption)
	expected :=
`<?xml version="1.0" encoding="utf-8"?>
<foo/>
<bar/>
`
	if err != nil {
		t.Error("Parsing has error:", err)
	}
	doc.AddChild("<bar></bar>")
	if doc.String() != expected {
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	CheckXmlMemoryLeaks(t)
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
	err = doc.RootElement().AddPreviousSibling("<bar></bar><cat></cat>")
	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	CheckXmlMemoryLeaks(t)
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
	doc.RootElement().AddNextSibling("<bar></bar>")
	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()
	CheckXmlMemoryLeaks(t)
}
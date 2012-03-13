package html

import (
	"testing"
	"gokogiri/help"
)

func TestParseDocument(t *testing.T) {
	expected :=
		`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">
<html><body><div><h1></h1></div></body></html>
`
	expected_xml :=
		`<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">
<html>
  <body>
    <div>
      <h1/>
    </div>
  </body>
</html>
`
	doc, err := Parse([]byte("<html><body><div><h1></div>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	if doc.String() != expected {
		println("got:\n", doc.String())
		println("expected:\n", expected)
		t.Error("the output of the html doc does not match")
	}

	if string(doc.ToXml(nil)) != expected_xml {
		println("got:\n", string(doc.ToXml(nil)))
		println("expected:\n", expected_xml)
		t.Error("the xml output of the html doc does not match")
	}

	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

func TestEmptyDocument(t *testing.T) {
	expected :=
		`<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.0 Transitional//EN" "http://www.w3.org/TR/REC-html40/loose.dtd">

`
	doc, err := Parse(nil, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	if doc.String() != expected {
		println(doc.String())
		t.Error("the output of the html doc does not match the empty xml")
	}
	doc.Free()
	help.CheckXmlMemoryLeaks(t)
}

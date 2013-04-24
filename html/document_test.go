package html

import "testing"

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

	s, _ := doc.ToXml(nil, nil)
	if string(s) != expected_xml {
		println("got:\n", string(s))
		println("expected:\n", expected_xml)
		t.Error("the xml output of the html doc does not match")
	}

	doc.Free()
	CheckXmlMemoryLeaks(t)
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
	CheckXmlMemoryLeaks(t)
}
func TestRemoveNamespace(t *testing.T) {
	xml := "<SOAP-ENV:Envelope xmlns:SOAP-ENV=\"http://schemas.xmlsoap.org/soap/envelope/\"><SOAP-ENV:Body><m:setPresence xmlns:m=\"http://schemas.microsoft.com/winrtc/2002/11/sip\"><m:presentity m:uri=\"test\"><m:availability m:aggregate=\"300\" m:description=\"online\"/><m:activity m:aggregate=\"400\" m:description=\"Active\"/><deviceName xmlns=\"http://schemas.microsoft.com/2002/09/sip/client/presence\" name=\"WIN-0DDABKC1UI8\"/></m:presentity></m:setPresence></SOAP-ENV:Body></SOAP-ENV:Envelope>"
	xml_no_namespace := "<Envelope><Body><setPresence><presentity uri=\"test\"><availability aggregate=\"300\" description=\"online\"/><activity aggregate=\"400\" description=\"Active\"/><deviceName name=\"WIN-0DDABKC1UI8\"/></presentity></setPresence></Body></Envelope>"

	doc, _ := Parse([]byte(xml), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)
	doc.Root().RecursivelyRemoveNamespaces()
	doc2, _ := Parse([]byte(xml_no_namespace), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	output := fmt.Sprintf("%v", doc)
	output_no_namespace := fmt.Sprintf("%v", doc2)
	if output != output_no_namespace {
		t.Errorf("Xml namespaces not removed!")
	}
}

/*
func TestHTMLFragmentEncoding(t *testing.T) {
	defer CheckXmlMemoryLeaks(t)

	input, output, error := getTestData(filepath.Join("tests", "document", "html_fragment_encoding"))

	if len(error) > 0 {
		t.Errorf("Error gathering test data for %v:\n%v\n", "html_fragment_encoding", error)
		t.FailNow()
	}

	expected := string(output)

	inputEncodingBytes := []byte("utf-8")

	buffer := make([]byte, 100)
	fragment, err := ParseFragment([]byte(input), inputEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes, buffer)

	if err != nil {
		println("WHAT")
		t.Error(err.Error())
	}

	if fragment.String() != expected {
		badOutput(fragment.String(), expected)
		t.Error("the output of the xml doc does not match")
	}

	fragment.Node.MyDocument().Free()
}
*/

package xml

import (
	"testing"
	"os"
	"gokogiri/help"
)

func TestDocuments(t *testing.T) {
	tests = collectTests()
	
	for _, test := range(tests) {
		input, err := ioutil.Readfile(filepath.Join(test, "input.txt"))
		
		if err != nil {
			t.Errorf("Couldn't read test (%v) input:\n%v\n", test, err.String())
		}
		
		output, err := ioutil.Readfile(filepath.Join(test, "output.txt"))
		
		if err != nil {
			t.Errorf("Couldn't read test (%v) output:\n%v\n", test, err.String())
		}

		RunDocumentParseTest(t, test, input, output)
	}
	
}

func RunDocumentParseTest(t *testing.T, name string, input []byte, output []byte) error string {
	
	defer help.CheckXmlMemoryLeaks(t)

	doc, err := Parse(input, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("parsing error:", err)
	}
	
	if doc.String() != string(output) {
		t.Error("the output of the xml doc does not match")
	}
	doc.Free()	
	
}

func collectTests() names []string {
	entries, err := ioutil.ReadDir(".")

	if err != nil {
		t.Errorf("Couldn't read tests:\n%v\n", err.String())
	}

	for _, entry := range(entries) {
		if entry.IsDirectory() {
			names = append(names, entry.Name)
		}
	}

	return
}



func TestBufferedDocuments(t *testing.T) {
}
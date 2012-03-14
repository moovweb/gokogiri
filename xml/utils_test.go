package xml

import (
	"io/ioutil"
	"fmt"
	"path/filepath"
	"testing"
	"strings"
	"gokogiri/help"
	)

func badOutput(actual string, expected string) {
	fmt.Printf("Got:\n[%v]\n", actual)
	fmt.Printf("Expected:\n[%v]\n", expected)
}

func RunTest(t *testing.T, suite string, name string, specificLogic func(doc *XmlDocument), extraAssertions ...func(doc *XmlDocument) (string, string, string) ) {
	defer help.CheckXmlMemoryLeaks(t)

	input, output, error := getTestData(filepath.Join("tests", suite, name))

	if len(error) > 0 {
		t.Errorf("Error gathering test data for %v:\n%v\n", name, error)
		t.FailNow()
	}

	expected := string(output)

	doc := parseInput(t, input)	

	if specificLogic != nil {
		specificLogic(doc)
	}

	if doc.String() != expected {
		badOutput(doc.String(), expected)
		t.Error("the output of the xml doc does not match")
	}

	for _, extraAssertion := range(extraAssertions) {
		actual, expected, message := extraAssertion(doc)

		if actual != expected {
			badOutput(actual, expected)
			t.Error(message)
		}
	}
	
	doc.Free()
}


func parseInput(t *testing.T, input interface{}) *XmlDocument {
	var realInput []byte

	switch thisInput := input.(type) {
	case []byte:
		realInput = thisInput
	case string:
		realInput = []byte(thisInput)
	default:
		t.Errorf("Unrecognized parsing input!")
	}

	doc, err := Parse(realInput, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("parsing error:", err.String())
		return nil
	}
	
	return doc
}

func getTestData(name string) (input []byte, output []byte, error string) {
	var errorMessage string
	offset := "\t"
	inputFile := filepath.Join(name, "input.txt")

	input, err := ioutil.ReadFile(inputFile)

	if err != nil {
		errorMessage += fmt.Sprintf("%vCouldn't read test (%v) input:\n%v\n", offset, name, offset+err.String())
	}

	output, err = ioutil.ReadFile(filepath.Join(name, "output.txt"))

	if err != nil {
		errorMessage += fmt.Sprintf("%vCouldn't read test (%v) output:\n%v\n", offset, name, offset+err.String())
	}

	return input, output, errorMessage
}

func collectTests(suite string) (names []string, error string) {
	testPath := filepath.Join("tests", suite)
	entries, err := ioutil.ReadDir(testPath)

	if err != nil {
		return nil, fmt.Sprintf("Couldn't read tests:\n%v\n", err.String())
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name, "_") || strings.HasPrefix(entry.Name, ".") {
			continue
		}

		if entry.IsDirectory() {
			names = append(names, filepath.Join(testPath, entry.Name))
		}
	}

	return
}
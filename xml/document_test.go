package xml

import (
	"testing"
	"gokogiri/help"
	"io/ioutil"
	"path/filepath"
	"strings"
	"fmt"
)

func TestDocuments(t *testing.T) {
	tests := collectTests(t)
	
	errors := make([]string, 0)

	print("\nTesting: [")

	for _, test := range(tests) {
		error := RunDocumentParseTest(t, test)

		if error != nil {
			errors = append(errors, fmt.Sprintf("Test %v failed:\n%v\n", test, *error))
			print("F")
		} else {
			print(".")
		}
	}
	
	println("]")

	if len(errors) > 0 {
		errorMessage := "\t" + strings.Join( strings.Split(strings.Join(errors, "\n\n"), "\n"), "\n\t")
		t.Errorf("\nSome tests failed! (%d passed / %d total) :\n%v", len(tests) - len(errors), len(tests), errorMessage)
	} else {
		fmt.Printf("\nAll (%d) tests passed!\n", len(tests))
	}
}

func RunDocumentParseTest(t *testing.T, name string) (error *string) {

	var errorMessage string
	offset := "\t"

	defer help.CheckXmlMemoryLeaks(t)

	input, err := ioutil.ReadFile(filepath.Join(name, "input.txt"))
	
	if err != nil {
		errorMessage += fmt.Sprintf("%vCouldn't read test (%v) input:\n%v\n", offset, name, offset + err.String())
	}
	
	output, err := ioutil.ReadFile(filepath.Join(name, "output.txt"))
	
	if err != nil {
		errorMessage += fmt.Sprintf("%vCouldn't read test (%v) output:\n%v\n", offset, name, offset + err.String())
	}


	doc, err := Parse(input, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		errorMessage = fmt.Sprintf("parsing error:%v\n", err)
	}
	
	if doc.String() != string(output) {
		formattedOutput := offset + strings.Join(strings.Split("[" + doc.String() + "]", "\n"), "\n" + offset)
		formattedExpectedOutput := offset + strings.Join(strings.Split("[" + string(output) + "]", "\n"), "\n" + offset)
		errorMessage = fmt.Sprintf("%v-- Got --\n%v\n%v-- Expected --\n%v\n", offset, formattedOutput, offset, formattedExpectedOutput)
	}
	doc.Free()	

	if len(errorMessage) > 0 {
		return &errorMessage	
	} 
	return nil

}

func collectTests(t *testing.T) (names []string) {
	testPath := "tests"
	entries, err := ioutil.ReadDir(testPath)

	if err != nil {
		t.Errorf("Couldn't read tests:\n%v\n", err.String())
	}

	for _, entry := range(entries) {
		if strings.HasPrefix(entry.Name, "_") || strings.HasPrefix(entry.Name, ".") {
			continue
		}

		if entry.IsDirectory() {
			names = append(names, filepath.Join(testPath, entry.Name) )
		}
	}

	return
}



func TestBufferedDocuments(t *testing.T) {
}
package xml

import (
	"testing"
	"github.com/moovweb/gokogiri/help"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"
	"fmt"
)

func TestDocuments(t *testing.T) {
	tests, err := collectTests("document")

	if len(err) > 0 {
		t.Errorf(err)
	}

	errors := make([]string, 0)

	print("\nTesting: Basic Parsing [")

	for _, test := range tests {
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
		errorMessage := "\t" + strings.Join(strings.Split(strings.Join(errors, "\n\n"), "\n"), "\n\t")
		t.Errorf("\nSome tests failed! (%d passed / %d total) :\n%v", len(tests)-len(errors), len(tests), errorMessage)
	} else {
		fmt.Printf("\nAll (%d) tests passed!\n", len(tests))
	}
}

func TestBufferedDocuments(t *testing.T) {
	tests, err := collectTests("document")

	if len(err) > 0 {
		t.Errorf(err)
	}

	errors := make([]string, 0)

	print("\nTesting: Buffered Parsing [")

	for _, test := range tests {
		error := RunParseDocumentWithBufferTest(t, test)

		if error != nil {
			errors = append(errors, fmt.Sprintf("Test %v failed:\n%v\n", test, *error))
			print("F")
		} else {
			print(".")
		}
	}

	println("]")

	if len(errors) > 0 {
		errorMessage := "\t" + strings.Join(strings.Split(strings.Join(errors, "\n\n"), "\n"), "\n\t")
		t.Errorf("\nSome tests failed! (%d passed / %d total) :\n%v", len(tests)-len(errors), len(tests), errorMessage)
	} else {
		fmt.Printf("\nAll (%d) tests passed!\n", len(tests))
	}
}

func RunParseDocumentWithBufferTest(t *testing.T, name string) (error *string) {
	var errorMessage string
	offset := "\t"

	defer help.CheckXmlMemoryLeaks(t)

	input, output, dataError := getTestData(name)

	if len(dataError) > 0 {
		errorMessage += dataError
	}

	buffer := make([]byte, 500000)

	doc, err := Parse(input, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		errorMessage = fmt.Sprintf("parsing error:%v\n", err)
	}

	if string(doc.ToBuffer(buffer)) != string(output) {
		formattedOutput := offset + strings.Join(strings.Split("["+doc.String()+"]", "\n"), "\n"+offset)
		formattedExpectedOutput := offset + strings.Join(strings.Split("["+string(output)+"]", "\n"), "\n"+offset)
		errorMessage = fmt.Sprintf("%v-- Got --\n%v\n%v-- Expected --\n%v\n", offset, formattedOutput, offset, formattedExpectedOutput)
	}
	doc.Free()

	if len(errorMessage) > 0 {
		return &errorMessage
	}
	return nil

}

func RunDocumentParseTest(t *testing.T, name string) (error *string) {

	var errorMessage string
	offset := "\t"

	defer help.CheckXmlMemoryLeaks(t)

	input, output, dataError := getTestData(name)

	if len(dataError) > 0 {
		errorMessage += dataError
	}

	doc, err := Parse(input, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		errorMessage = fmt.Sprintf("parsing error:%v\n", err)
	}

	if doc.String() != string(output) {
		formattedOutput := offset + strings.Join(strings.Split("["+doc.String()+"]", "\n"), "\n"+offset)
		formattedExpectedOutput := offset + strings.Join(strings.Split("["+string(output)+"]", "\n"), "\n"+offset)
		errorMessage = fmt.Sprintf("%v-- Got --\n%v\n%v-- Expected --\n%v\n", offset, formattedOutput, offset, formattedExpectedOutput)
		testOutput := filepath.Join(name, "test_output.txt")
		ioutil.WriteFile(testOutput, []byte(doc.String()), os.FileMode(0666))
		errorMessage += fmt.Sprintf("%v Output test output to: %v\n", offset, testOutput)
	}
	doc.Free()

	if len(errorMessage) > 0 {
		return &errorMessage
	}
	return nil

}

func BenchmarkDocOutput(b *testing.B) {
	b.StopTimer()

	tests, err := collectTests("document")

	if len(err) > 0 {
		fmt.Printf(err)
		return
	}

	docs := make([]*XmlDocument, 0)

	for _, testName := range tests {

		input, _, dataError := getTestData(testName)

		if len(dataError) > 0 {
			fmt.Printf(dataError)
			return
		}
		doc, err := Parse(input, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

		if err != nil {
			fmt.Printf("parsing error:%v\n", err)
			return
		}
		docs = append(docs, doc)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for index, _ := range tests {
			_ = docs[index].String()
		}
	}

}

func BenchmarkDocOutputToBuffer(b *testing.B) {
	b.StopTimer()

	tests, err := collectTests("document")

	if len(err) > 0 {
		fmt.Printf(err)
		return
	}

	docs := make([]*XmlDocument, 0)

	for _, testName := range tests {

		input, _, dataError := getTestData(testName)

		if len(dataError) > 0 {
			fmt.Printf(dataError)
			return
		}
		doc, err := Parse(input, DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

		if err != nil {
			fmt.Printf("parsing error:%v\n", err)
			return
		}
		docs = append(docs, doc)
	}

	buffer := make([]byte, 500*1024)

	b.StartTimer()

	for i := 0; i < b.N; i++ {

		for index, _ := range tests {

			_ = docs[index].ToBuffer(buffer)

		}
	}

}

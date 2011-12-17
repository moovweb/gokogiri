package test

import (
	"libxml"
	"io/ioutil"
	"testing"
	"log"
)

func BenchmarkHtmlParsingStringVeryShort(b *testing.B) {
    b.StopTimer()
    input := "<html><body><h1>Hello World</h1></body></html>"
    b.StartTimer()
    for i := 0; i < b.N; i++ {
		doc := libxml.HtmlParseString(input)
		doc.Free()
    }
}

func BenchmarkHtmlParsingStringNormalPage(b *testing.B) {
    b.StopTimer()
    inputBytes, _ := ioutil.ReadFile("bn_body.html")
	input := string(inputBytes)
	log.Printf("\ninput html size: %d\n", len(input))
    b.StartTimer()
    for i := 0; i < b.N; i++ {
		doc := libxml.HtmlParseString(input)
		doc.Free()
    }
}

func BenchmarkHtmlParsingBytesNormalPage(b *testing.B) {
    b.StopTimer()
    input, _ := ioutil.ReadFile("bn_body.html")
	log.Printf("\ninput html size: %d\n", len(input))
    b.StartTimer()
    for i := 0; i < b.N; i++ {
		doc := libxml.HtmlParseBytes(input)
		doc.Free()
    }
}
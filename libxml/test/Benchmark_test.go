package test

import (
	"libxml"
	"testing"
)

func BenchmarkSimpleXmlParsing(b *testing.B) {
    b.StopTimer()
    input := "<root>hi</root>"
    b.StartTimer()
    for i := 0; i < b.N; i++ {
		doc := libxml.XmlParseString(input)
		doc.Free()
    }
}
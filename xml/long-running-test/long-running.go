package main

import (
    "gokogiri/xml"
)

func parseHtmlSnippet()(bool) {
    _, _ = xml.Parse([]byte(""), xml.DefaultEncodingBytes, nil, xml.DefaultParseOption, xml.DefaultEncodingBytes)
    //defer doc.Free()
    return true
}

func main(){
    for {
        _ = parseHtmlSnippet()
    }
}

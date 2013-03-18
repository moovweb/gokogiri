package main

import (
    "gokogiri/html"
)

func parseHtml()(bool) {
    _, _ = html.Parse([]byte(""), html.DefaultEncodingBytes, nil, html.DefaultParseOption, html.DefaultEncodingBytes)
    //defer doc.Free()
    return true
}

func parseHtmlFrag()(bool) {
    _, _ = html.ParseFragment([]byte("<div><a/>"), html.DefaultEncodingBytes, nil, html.DefaultParseOption, html.DefaultEncodingBytes)
    //defer doc.Free()
    return true
}

func main(){
    for {
        _ = parseHtmlFrag()
    }
}

Gokogiri
========
LibXML bindings for the Go programming language.
------------------------------------------------
By Zhigang Chen and Hampton Catlin


This is a major rewrite from v0 in the following places:

- Separation of XML and HTML
- Put more burden of memory allocation/deallocation on Go
- Fragment parsing -- no more deep-copy
- Serialization
- Some API adjustment

To install:

- sudo apt-get install libxml2-dev
- go get github.com/moovweb/gokogiri

To run test:

- go test github.com/moovweb/gokogiri/html
- go test github.com/moovweb/gokogiri/xml

Basic example:

    package main

    import (
      "net/http"
      "io/ioutil"
      "github.com/moovweb/gokogiri"
    )

    func main() {
      // fetch and read a web page
      resp, _ := http.Get("http://www.google.com")
      page, _ := ioutil.ReadAll(resp.Body)

      // parse the web page
      doc, _ := gokogiri.ParseHtml(page)

      // perform operations on the parsed page -- consult the tests for examples

      // important -- don't forget to free the resources when you're done!
      doc.Free()
    }
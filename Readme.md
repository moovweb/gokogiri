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

## Installation

### OSX

If your environment doesn't have "pkg-config" and "libxml2".

```bash
brew install libxml2
brew install pkg-config
```

By the way, don't forget to add these paths either ~/.bashrc or  ~/.zshrc.

```bash
export LDFLAGS="-L/usr/local/opt/libxml2/lib"
export CPPFLAGS="-I/usr/local/opt/libxml2/include"
export PKG_CONFIG_PATH="/usr/local/opt/libxml2/lib/pkgconfig"
```

Finally,

```bash
go get github.com/moovweb/gokogiri
```

### Linux

```bash
sudo apt-get install libxml2-dev
go get github.com/moovweb/gokogiri
```

## Running tests

```bash
go test github.com/moovweb/gokogiri/...
```

## Basic example

```go
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
  defer doc.Free()
}
```

include $(GOROOT)/src/Make.inc
TARG=libxml
CGOFILES=\
          src/constants.go\
          src/libxml.go\
          src/XmlDoc.go\
					src/XmlNode.go\
          src/XmlElement.go\
          src/XPathContext.go\
          src/XmlNodeSet.go
CGO_LDFLAGS=-lxml2
CGO_CFLAGS=-I/usr/include/libxml2
include $(GOROOT)/src/Make.pkg
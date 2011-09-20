include $(GOROOT)/src/Make.inc
TARG=libxml
CGOFILES=\
          constants.go\
          libxml.go\
          XmlDoc.go\
          Node.go\
          XmlElement.go\
          XPathContext.go\
          NodeSet.go
CGO_LDFLAGS=-lxml2
CGO_CFLAGS=-I/usr/include/libxml2
include $(GOROOT)/src/Make.pkg
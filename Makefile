include $(GOROOT)/src/Make.inc
TARG=goxml
CGOFILES=goxml.go 
CGO_LDFLAGS=-lxml2
CGO_CFLAGS=-I/usr/include/libxml2
include $(GOROOT)/src/Make.pkg
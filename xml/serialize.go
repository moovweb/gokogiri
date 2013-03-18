package xml

/*
#include "helper.h"
#include <string.h>
*/
import "C"
import "unsafe"

type WriteBuffer struct {
	Node   *XmlNode
	Buffer []byte
	Offset int
}

func serialize(xmlNode *XmlNode, format int, encoding, outputBuffer []byte) ([]byte, int) {
	nodePtr := unsafe.Pointer(xmlNode.Ptr)
	var encodingPtr unsafe.Pointer
	if len(encoding) == 0 {
		encoding = xmlNode.OutEncoding
	}
	if len(encoding) > 0 {
		encodingPtr = unsafe.Pointer(&(encoding[0]))
	} else {
		encodingPtr = nil
	}

	wbuffer := &WriteBuffer{Node: xmlNode, Buffer: outputBuffer}
	wbufferPtr := unsafe.Pointer(wbuffer)

	format |= XML_SAVE_FORMAT
	ret := int(C.xmlSaveNode(wbufferPtr, nodePtr, encodingPtr, C.int(format)))
	if ret < 0 {
		println("output error!!!")
		return nil, 0
	}

	return wbuffer.Buffer, wbuffer.Offset
}

//export xmlNodeWriteCallback
func xmlNodeWriteCallback(wbufferObj unsafe.Pointer, data unsafe.Pointer, data_len C.int) {
	wbuffer := (*WriteBuffer)(wbufferObj)
	offset := wbuffer.Offset

	if offset > len(wbuffer.Buffer) {
		panic("fatal error in xmlNodeWriteCallback")
	}

	buffer := wbuffer.Buffer[:offset]
	dataLen := int(data_len)

	if dataLen > 0 {
		if len(buffer)+dataLen > cap(buffer) {
			newBuffer := grow(buffer, dataLen)
			wbuffer.Buffer = newBuffer
		}
		destBufPtr := unsafe.Pointer(&(wbuffer.Buffer[offset]))
		C.memcpy(destBufPtr, data, C.size_t(dataLen))
		wbuffer.Offset += dataLen
	}
}

func grow(buffer []byte, n int) (newBuffer []byte) {
	newBuffer = makeSlice(2*cap(buffer) + n)
	copy(newBuffer, buffer)
	return
}

func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}
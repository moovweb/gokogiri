package xml

import "errors"

//for node functions
var (
	ERR_UNDEFINED_COERCE_PARAM               = errors.New("unexpected parameter type in coerce")
	ERR_UNDEFINED_SET_CONTENT_PARAM          = errors.New("unexpected parameter type in SetContent")
	ERR_UNDEFINED_SEARCH_PARAM               = errors.New("unexpected parameter type in Search")
	ERR_CANNOT_MAKE_DUCMENT_AS_CHILD         = errors.New("cannot add a document node as a child")
	ERR_CANNOT_COPY_TEXT_NODE_WHEN_ADD_CHILD = errors.New("cannot copy a text node when adding it")
)

//parse
var ERR_FAILED_TO_PARSE_XML                  = errors.New("failed to parse xml input")

//serialize
var	ErrTooLarge                              = errors.New("Output buffer too large")

//fragment
var ErrFailParseFragment = errors.New("failed to parse xml fragment")
var ErrEmptyFragment = errors.New("empty xml fragment")
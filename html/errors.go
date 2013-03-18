package html

import "errors"

var ErrSetMetaEncoding = errors.New("Set Meta Encoding failed")
var ERR_FAILED_TO_PARSE_HTML = errors.New("failed to parse html input")

var ErrFailParseFragment = errors.New("failed to parse html fragment")
var ErrEmptyFragment = errors.New("empty html fragment")
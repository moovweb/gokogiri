package libxml

const (
	//parser option 
	HTML_PARSE_RECOVER   = 1 << 0  //relaxed parsing 
	HTML_PARSE_NOERROR   = 1 << 5  //suppress error reports 
	HTML_PARSE_NOWARNING = 1 << 6  //suppress warning reports 
	HTML_PARSE_PEDANTIC  = 1 << 7  //pedantic error reporting 
	HTML_PARSE_NOBLANKS  = 1 << 8  //remove blank nodes 
	HTML_PARSE_NONET     = 1 << 11 //forbid network access 
	HTML_PARSE_COMPACT   = 1 << 16 //compact small text nodes

    XML_PARSE_RECOVER    = 1 << 0 //relaxed parsing
    XML_PARSE_NOERROR   = 1 << 5  //suppress error reports 
    XML_PARSE_NOWARNING = 1 << 6  //suppress warning reports 
    XML_PARSE_NONET     = 1 << 11 //forbid network access 

)

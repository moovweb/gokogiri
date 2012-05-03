package html

import (
	"testing"
	"github.com/moovweb/gokogiri/help"
)

func TestInnerScript(t *testing.T) {
	defer help.CheckXmlMemoryLeaks(t)

	doc, err := Parse([]byte("<html><body><div><h1></div>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	h1 := doc.Root().FirstChild().FirstChild().FirstChild()
	h1.SetInnerHtml("<script>if (suppressReviews !== 'true' && app == 'PRR') { ok = true; }</script>")
	if h1.String() != "<h1><script>if (suppressReviews !== 'true' && app == 'PRR') { ok = true; }</script></h1>" {
		t.Error("script does not match")
	}
	doc.Free()
}

func TestInnerScript2(t *testing.T) {
	defer help.CheckXmlMemoryLeaks(t)
	script := `<script>try {
var productNAPage = "",
suppressReviews = "false";
var bvtoken = MACYS.util.Cookie.get("BazaarVoiceToken","GCs");
//bvtoken=bvtoken.substring(0,bvtoken.length-1);
$BV.configure("global", {
userToken: bvtoken,
productId: '531726',
submissionUI: 'LIGHTBOX',
submissionContainerUrl: window.location.href,
allowSamePageSubmission: true,
doLogin: function(callback, success_url) {
MACYS.util.Cookie.set("FORWARDPAGE_KEY",success_url);
window.location = 'https://www.macys.com/signin/index.ognc?fromPage=pdpReviews';
},
doShowContent: function(app, dc, sub, sr) {
if (suppressReviews !== 'true' && app == "PRR") {
MACYS.pdp.showReviewsTab();
} else if (productNAPage !== 'true' && app == "QA") {
MACYS.pdp.showQATab();
}
}
});
if (suppressReviews !== 'true') {
$BV.ui('rr', 'show_reviews', {
});
}
$BV.ui("qa", "show_questions", {
subjectType: 'product'
});
} catch ( e ) { }</script>`

	doc, err := Parse([]byte("<html><body><div><h1></div>"), DefaultEncodingBytes, nil, DefaultParseOption, DefaultEncodingBytes)

	if err != nil {
		t.Error("Parsing has error:", err)
		return
	}

	h1 := doc.Root().FirstChild().FirstChild().FirstChild()
	h1.SetInnerHtml(script)
	if h1.String() != "<h1>"+script+"</h1>" {
		t.Error("script does not match")
	}
	doc.Free()
}

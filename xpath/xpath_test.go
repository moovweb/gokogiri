package xpath

//please check the search tests in gokogiri/xml and gokogiri/html
import "testing"

func TestCompileGoodExpr(t *testing.T) {
	defer CheckXmlMemoryLeaks(t)
	e := Compile(`./*`)
	if e == nil {
		t.Error("expr should be good")
	}
}

func TestCompileBadExpr(t *testing.T) {
	defer CheckXmlMemoryLeaks(t)
	e := Compile("./")
	if e != nil {
		t.Error("expr should be bad")
	}
}

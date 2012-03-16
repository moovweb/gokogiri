package xml

import "testing"

func TestSearch(t *testing.T) {

	testLogic := func(t *testing.T, doc *XmlDocument) {
		root := doc.Root()
		result, _ := root.Search(".//*[@class]")
		if len(result) != 2 {
			t.Error("search at root does not match")
		}
		result, _ = root.Search("//*[@class]")
		if len(result) != 3 {
			t.Error("search at root does not match")
		}
		result, _ = doc.Search(".//*[@class]")
		if len(result) != 3 {
			t.Error("search at doc does not match")
		}
		result, _ = doc.Search(".//*[@class='shine']")
		if len(result) != 2 {
			t.Error("search with value at doc does not match")
		}
	}

	RunTest(t, "node", "search", testLogic)
}

func BenchmarkSearch(b *testing.B) {

	benchmarkLogic := func(b *testing.B, doc *XmlDocument) {
		root := doc.Root()

		for i := 0; i < b.N; i++ {
			root.Search(".//*[@class]")
		}
	}

	RunBenchmark(b, "node", "search", benchmarkLogic)
}

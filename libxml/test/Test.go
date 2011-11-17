package test

import "testing"
import "io/ioutil"
import "runtime"
import "libxml/help"

func AssertNil(t *testing.T, value interface{}, what string) {
	if value != nil {
		t.Error(what, "should be nil")
	}
}
func Equal(t *testing.T, value, expected interface{}) {
	if value != expected {
		t.Error("Expected: ", expected, "\nBut got: ", value)
	}
}
func Assert(t *testing.T, value interface{}, what string) interface{} {
	if value == nil {
		t.Error("Assertion failed: ", what)
	}
	return value
}

func LoadFile(name string) string {
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		print(err.String())
	}
	return string(contents)
}

func runParallel(testFunc func(chan bool), concurrency int) {
	help.XmlInitParser()

	runtime.GOMAXPROCS(4)
	done := make(chan bool, concurrency)
	for i := 0; i < concurrency; i++ {
		go testFunc(done)
	}
	for i := 0; i < concurrency; i++ {
		<-done
		<-done
	}
	runtime.GOMAXPROCS(1)

	help.XmlCleanUpParser()
}

const numConcurrentRuns = 100

package help

import "testing"

func CheckXmlMemoryLeaks(t *testing.T) {
	LibxmlCleanUpParser()
	if ! LibxmlCheckMemoryLeak() {
		t.Errorf("Memeory leaks: %d!!!", LibxmlGetMemoryAllocation())
		LibxmlReportMemoryLeak()
	}
}
package xml

import "testing"
import "gokogiri/help"


func CheckXmlMemoryLeaks(t *testing.T) {
	help.LibxmlCleanUpParser()
	if ! help.LibxmlCheckMemoryLeak() {
		t.Errorf("Memeory leaks: %d!!!", help.LibxmlGetMemoryAllocation())
		help.LibxmlReportMemoryLeak()
	}
}
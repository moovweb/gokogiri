package mem

import "testing"

func TestLibxml(t *testing.T) {
	if AllocSize() != 0 {
		t.Fatal(AllocSize(), "remaining allocations")
	}
}

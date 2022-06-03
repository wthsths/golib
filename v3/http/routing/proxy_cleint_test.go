package gl_routing

import (
	"fmt"
	"testing"
)

func Test_Len_Of_Slice(t *testing.T) {
	var slice []int

	if slice != nil {
		t.Fatalf("slice must be nil")
	}

	// Must not panic.
	sliceLen := len(slice)
	fmt.Println(sliceLen)
}

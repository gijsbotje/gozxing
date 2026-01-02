package encoder

import (
	"testing"
)

func TestDimensions(t *testing.T) {
	dims := NewDimensions(2, 10, 3, 15)
	if dims == nil {
		t.Fatalf("NewDimensions returns nil")
	}

	if dims.GetMinCols() != 2 {
		t.Fatalf("GetMinCols() = %d, expect 2", dims.GetMinCols())
	}
	if dims.GetMaxCols() != 10 {
		t.Fatalf("GetMaxCols() = %d, expect 10", dims.GetMaxCols())
	}
	if dims.GetMinRows() != 3 {
		t.Fatalf("GetMinRows() = %d, expect 3", dims.GetMinRows())
	}
	if dims.GetMaxRows() != 15 {
		t.Fatalf("GetMaxRows() = %d, expect 15", dims.GetMaxRows())
	}
}


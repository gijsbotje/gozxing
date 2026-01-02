package encoder

import (
	"testing"
)

func TestBarcodeMatrix(t *testing.T) {
	matrix := NewBarcodeMatrix(5, 10)
	if matrix == nil {
		t.Fatalf("NewBarcodeMatrix returns nil")
	}

	matrix.StartRow()
	row := matrix.GetCurrentRow()
	if row == nil {
		t.Fatalf("GetCurrentRow returns nil")
	}

	row.AddBar(true, 5)
	row.AddBar(false, 5)

	matrix.StartRow()
	row2 := matrix.GetCurrentRow()
	row2.AddBar(false, 10)

	scaled := matrix.GetScaledMatrix(1, 1)
	if len(scaled) != 5 {
		t.Fatalf("Scaled matrix height = %d, expect 5", len(scaled))
	}
	if len(scaled[0]) == 0 {
		t.Fatalf("Scaled matrix width should be > 0")
	}
}

func TestBarcodeMatrix_GetScaledMatrix(t *testing.T) {
	matrix := NewBarcodeMatrix(3, 5)
	matrix.StartRow()
	matrix.GetCurrentRow().AddBar(true, 10)

	scaled := matrix.GetScaledMatrix(2, 3)
	if len(scaled) != 9 { // 3 * 3
		t.Fatalf("Scaled matrix height = %d, expect 9", len(scaled))
	}
}


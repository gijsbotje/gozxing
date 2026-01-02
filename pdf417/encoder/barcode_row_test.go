package encoder

import (
	"testing"
)

func TestBarcodeRow_AddBar(t *testing.T) {
	row := NewBarcodeRow(100)
	row.AddBar(true, 5)
	row.AddBar(false, 3)
	row.AddBar(true, 2)

	scaled := row.GetScaledRow(1)
	if len(scaled) != 100 {
		t.Fatalf("Scaled row length = %d, expect 100", len(scaled))
	}

	// Check first 5 should be 1 (black)
	for i := 0; i < 5; i++ {
		if scaled[i] != 1 {
			t.Fatalf("Scaled[%d] = %d, expect 1", i, scaled[i])
		}
	}
	// Next 3 should be 0 (white)
	for i := 5; i < 8; i++ {
		if scaled[i] != 0 {
			t.Fatalf("Scaled[%d] = %d, expect 0", i, scaled[i])
		}
	}
	// Next 2 should be 1 (black)
	for i := 8; i < 10; i++ {
		if scaled[i] != 1 {
			t.Fatalf("Scaled[%d] = %d, expect 1", i, scaled[i])
		}
	}
}

func TestBarcodeRow_GetScaledRow(t *testing.T) {
	row := NewBarcodeRow(10)
	row.AddBar(true, 5)
	row.AddBar(false, 5)

	scaled := row.GetScaledRow(2)
	if len(scaled) != 20 {
		t.Fatalf("Scaled row length = %d, expect 20", len(scaled))
	}

	// First 10 should be black
	for i := 0; i < 10; i++ {
		if scaled[i] != 1 {
			t.Fatalf("Scaled[%d] = %d, expect 1", i, scaled[i])
		}
	}
	// Next 10 should be white
	for i := 10; i < 20; i++ {
		if scaled[i] != 0 {
			t.Fatalf("Scaled[%d] = %d, expect 0", i, scaled[i])
		}
	}
}


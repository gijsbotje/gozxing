package encoder

import (
	"testing"
)

func TestPDF417_GenerateBarcodeLogic(t *testing.T) {
	pdf417 := NewPDF417()
	err := pdf417.GenerateBarcodeLogic("Hello", 2)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogic returns error: %v", err)
	}

	matrix := pdf417.GetBarcodeMatrix()
	if matrix == nil {
		t.Fatalf("GetBarcodeMatrix returns nil")
	}
}

func TestPDF417_SetDimensions(t *testing.T) {
	pdf417 := NewPDF417()
	pdf417.SetDimensions(10, 2, 15, 3)

	err := pdf417.GenerateBarcodeLogic("test", 2)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogic returns error: %v", err)
	}
}

func TestPDF417_SetCompaction(t *testing.T) {
	pdf417 := NewPDF417()
	pdf417.SetCompaction(Compaction_TEXT)

	err := pdf417.GenerateBarcodeLogic("test", 2)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogic returns error: %v", err)
	}
}

func TestPDF417_SetCompact(t *testing.T) {
	pdf417 := NewPDF417()
	pdf417.SetCompact(true)

	err := pdf417.GenerateBarcodeLogic("test", 2)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogic returns error: %v", err)
	}
}

func TestPDF417_GenerateBarcodeLogicWithAutoECI(t *testing.T) {
	pdf417 := NewPDF417()
	err := pdf417.GenerateBarcodeLogicWithAutoECI("Hello", 2, false)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogicWithAutoECI returns error: %v", err)
	}

	matrix := pdf417.GetBarcodeMatrix()
	if matrix == nil {
		t.Fatalf("GetBarcodeMatrix returns nil")
	}
}

func TestPDF417_GenerateBarcodeLogicLongMessage(t *testing.T) {
	pdf417 := NewPDF417()
	longMsg := "This is a longer message to test PDF417 encoding with more data. " +
		"It should be able to handle messages of various lengths."
	err := pdf417.GenerateBarcodeLogic(longMsg, 2)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogic returns error: %v", err)
	}
}

func TestPDF417_GenerateBarcodeLogicNumeric(t *testing.T) {
	pdf417 := NewPDF417()
	pdf417.SetCompaction(Compaction_NUMERIC)
	err := pdf417.GenerateBarcodeLogic("12345678901234567890", 2)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogic returns error: %v", err)
	}
}

func TestPDF417_GenerateBarcodeLogicByte(t *testing.T) {
	pdf417 := NewPDF417()
	pdf417.SetCompaction(Compaction_BYTE)
	err := pdf417.GenerateBarcodeLogic("test", 2)
	if err != nil {
		t.Fatalf("GenerateBarcodeLogic returns error: %v", err)
	}
}


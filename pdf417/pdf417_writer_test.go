/*
 * Copyright 2016 ZXing authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pdf417

import (
	"testing"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/pdf417/encoder"
)

func TestPDF417Writer_Encode(t *testing.T) {
	writer := NewPDF417Writer()
	hints := make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_MARGIN] = 0
	size := 64
	matrix, err := writer.Encode("Hello Google", gozxing.BarcodeFormat_PDF_417, size, size, hints)
	if err != nil {
		t.Fatalf("Encode returns error: %v", err)
	}
	if matrix == nil {
		t.Fatalf("Encode returns nil matrix")
	}

	// Verify dimensions - PDF417 may not produce exact size, but should be reasonable
	w, h := matrix.GetWidth(), matrix.GetHeight()
	if w <= 0 || h <= 0 {
		t.Fatalf("Matrix size = %vx%v, should be positive", w, h)
	}
	if w > size*2 || h > size*2 {
		t.Fatalf("Matrix size = %vx%v, seems too large", w, h)
	}

	// Verify it's not all empty
	hasBlack := false
	for y := 0; y < matrix.GetHeight(); y++ {
		for x := 0; x < matrix.GetWidth(); x++ {
			if matrix.Get(x, y) {
				hasBlack = true
				break
			}
		}
		if hasBlack {
			break
		}
	}
	if !hasBlack {
		t.Fatalf("Matrix should have some black pixels")
	}
}

func TestPDF417Writer_EncodeFailed(t *testing.T) {
	writer := NewPDF417Writer()
	var err error

	_, err = writer.EncodeWithoutHint("", gozxing.BarcodeFormat_PDF_417, 30, 30)
	if err == nil {
		t.Fatalf("Encode must be Error for empty string")
	}

	_, err = writer.EncodeWithoutHint("test", gozxing.BarcodeFormat_QR_CODE, 30, 30)
	if err == nil {
		t.Fatalf("Encode must be Error for wrong format")
	}
}

func TestPDF417Writer_EncodeWithHints(t *testing.T) {
	writer := NewPDF417Writer()
	hints := make(map[gozxing.EncodeHintType]interface{})

	// Test with compact mode
	hints[gozxing.EncodeHintType_PDF417_COMPACT] = true
	matrix, err := writer.Encode("test", gozxing.BarcodeFormat_PDF_417, 100, 100, hints)
	if err != nil {
		t.Fatalf("Encode with compact returns error: %v", err)
	}
	if matrix == nil {
		t.Fatalf("Encode returns nil matrix")
	}

	// Test with compaction mode
	hints = make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_PDF417_COMPACTION] = "TEXT"
	matrix, err = writer.Encode("test", gozxing.BarcodeFormat_PDF_417, 100, 100, hints)
	if err != nil {
		t.Fatalf("Encode with TEXT compaction returns error: %v", err)
	}

	// Test with dimensions
	hints = make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_PDF417_DIMENSIONS] = encoder.NewDimensions(2, 10, 2, 10)
	matrix, err = writer.Encode("test", gozxing.BarcodeFormat_PDF_417, 100, 100, hints)
	if err != nil {
		t.Fatalf("Encode with dimensions returns error: %v", err)
	}

	// Test with error correction level
	hints = make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_ERROR_CORRECTION] = 3
	matrix, err = writer.Encode("test", gozxing.BarcodeFormat_PDF_417, 100, 100, hints)
	if err != nil {
		t.Fatalf("Encode with error correction returns error: %v", err)
	}

	// Test with margin
	hints = make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_MARGIN] = 10
	matrix, err = writer.Encode("test", gozxing.BarcodeFormat_PDF_417, 100, 100, hints)
	if err != nil {
		t.Fatalf("Encode with margin returns error: %v", err)
	}
}

func TestPDF417Writer_EncodeLongMessage(t *testing.T) {
	writer := NewPDF417Writer()
	longMsg := "This is a longer message to test PDF417 encoding with more data. " +
		"It should be able to handle messages of various lengths and encode them properly."
	matrix, err := writer.EncodeWithoutHint(longMsg, gozxing.BarcodeFormat_PDF_417, 200, 200)
	if err != nil {
		t.Fatalf("Encode long message returns error: %v", err)
	}
	if matrix == nil {
		t.Fatalf("Encode returns nil matrix")
	}
}

func TestPDF417Writer_EncodeNumeric(t *testing.T) {
	writer := NewPDF417Writer()
	hints := make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_PDF417_COMPACTION] = encoder.Compaction_NUMERIC
	matrix, err := writer.Encode("12345678901234567890", gozxing.BarcodeFormat_PDF_417, 200, 200, hints)
	if err != nil {
		t.Fatalf("Encode numeric returns error: %v", err)
	}
	if matrix == nil {
		t.Fatalf("Encode returns nil matrix")
	}
}


/*
 * Copyright 2012 ZXing authors
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
	"strconv"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/pdf417/encoder"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// PDF417Writer implements the Writer interface for PDF417 barcodes
type PDF417Writer struct{}

// NewPDF417Writer creates a new PDF417Writer
func NewPDF417Writer() *PDF417Writer {
	return &PDF417Writer{}
}

const (
	// default white space (margin) around the code
	WHITE_SPACE = 30
	// default error correction level
	DEFAULT_ERROR_CORRECTION_LEVEL = 2
)

// EncodeWithoutHint encodes a barcode using the default settings
func (w *PDF417Writer) EncodeWithoutHint(contents string, format gozxing.BarcodeFormat, width, height int) (*gozxing.BitMatrix, error) {
	return w.Encode(contents, format, width, height, nil)
}

// Encode encodes a barcode with hints
func (w *PDF417Writer) Encode(contents string, format gozxing.BarcodeFormat, width, height int, hints map[gozxing.EncodeHintType]interface{}) (*gozxing.BitMatrix, error) {
	if format != gozxing.BarcodeFormat_PDF_417 {
		return nil, gozxing.NewWriterException("Can only encode PDF_417, but got %v", format)
	}

	pdf417 := encoder.NewPDF417()
	margin := WHITE_SPACE
	errorCorrectionLevel := DEFAULT_ERROR_CORRECTION_LEVEL
	autoECI := false

	if hints != nil {
		if compact, ok := hints[gozxing.EncodeHintType_PDF417_COMPACT]; ok {
			if compactStr, ok := compact.(string); ok {
				if compactBool, err := strconv.ParseBool(compactStr); err == nil {
					pdf417.SetCompact(compactBool)
				}
			} else if compactBool, ok := compact.(bool); ok {
				pdf417.SetCompact(compactBool)
			}
		}
		if compaction, ok := hints[gozxing.EncodeHintType_PDF417_COMPACTION]; ok {
			if compactionStr, ok := compaction.(string); ok {
				if comp, err := encoder.Compaction_ValueOf(compactionStr); err == nil {
					pdf417.SetCompaction(comp)
				}
			} else if comp, ok := compaction.(encoder.Compaction); ok {
				pdf417.SetCompaction(comp)
			}
		}
		if dimensions, ok := hints[gozxing.EncodeHintType_PDF417_DIMENSIONS]; ok {
			if dims, ok := dimensions.(*encoder.Dimensions); ok {
				pdf417.SetDimensions(dims.GetMaxCols(), dims.GetMinCols(), dims.GetMaxRows(), dims.GetMinRows())
			}
		}
		if m, ok := hints[gozxing.EncodeHintType_MARGIN]; ok {
			if marginInt, ok := m.(int); ok {
				margin = marginInt
			} else if marginStr, ok := m.(string); ok {
				if marginParsed, err := strconv.Atoi(marginStr); err == nil {
					margin = marginParsed
				}
			}
		}
		if ec, ok := hints[gozxing.EncodeHintType_ERROR_CORRECTION]; ok {
			if ecInt, ok := ec.(int); ok {
				errorCorrectionLevel = ecInt
			} else if ecStr, ok := ec.(string); ok {
				if ecParsed, err := strconv.Atoi(ecStr); err == nil {
					errorCorrectionLevel = ecParsed
				}
			}
		}
		if charset, ok := hints[gozxing.EncodeHintType_CHARACTER_SET]; ok {
			if charsetStr, ok := charset.(string); ok {
				// Try to get encoding from name
				if enc, err := getEncodingFromName(charsetStr); err == nil {
					pdf417.SetEncoding(enc)
				}
			} else if enc, ok := charset.(encoding.Encoding); ok {
				pdf417.SetEncoding(enc)
			}
		}
		if autoECIHint, ok := hints[gozxing.EncodeHintType_PDF417_AUTO_ECI]; ok {
			if autoECIStr, ok := autoECIHint.(string); ok {
				if autoECIBool, err := strconv.ParseBool(autoECIStr); err == nil {
					autoECI = autoECIBool
				}
			} else if autoECIBool, ok := autoECIHint.(bool); ok {
				autoECI = autoECIBool
			}
		}
	}

	return bitMatrixFromEncoder(pdf417, contents, errorCorrectionLevel, width, height, margin, autoECI)
}

// bitMatrixFromEncoder takes encoder, accounts for width/height, and retrieves bit matrix
func bitMatrixFromEncoder(pdf417 *encoder.PDF417, contents string, errorCorrectionLevel, width, height, margin int, autoECI bool) (*gozxing.BitMatrix, error) {
	err := pdf417.GenerateBarcodeLogicWithAutoECI(contents, errorCorrectionLevel, autoECI)
	if err != nil {
		return nil, err
	}

	aspectRatio := 4
	originalScale := pdf417.GetBarcodeMatrix().GetScaledMatrix(1, aspectRatio)
	rotated := false
	if (height > width) != (len(originalScale[0]) < len(originalScale)) {
		originalScale = rotateArray(originalScale)
		rotated = true
	}

	scaleX := width / len(originalScale[0])
	scaleY := height / len(originalScale)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	if scale > 1 {
		scaledMatrix := pdf417.GetBarcodeMatrix().GetScaledMatrix(scale, scale*aspectRatio)
		if rotated {
			scaledMatrix = rotateArray(scaledMatrix)
		}
		return bitMatrixFromBitArray(scaledMatrix, margin)
	}
	return bitMatrixFromBitArray(originalScale, margin)
}

// bitMatrixFromBitArray takes an array holding the values of the PDF 417
// input: a byte array of information with 0 is black, and 1 is white
// margin: border around the barcode
// Returns BitMatrix of the input
func bitMatrixFromBitArray(input [][]byte, margin int) (*gozxing.BitMatrix, error) {
	// Creates the bit matrix with extra space for whitespace
	output, err := gozxing.NewBitMatrix(len(input[0])+2*margin, len(input)+2*margin)
	if err != nil {
		return nil, err
	}
	output.Clear()
	for y, yOutput := 0, output.GetHeight()-margin-1; y < len(input); y, yOutput = y+1, yOutput-1 {
		inputY := input[y]
		for x := 0; x < len(input[0]); x++ {
			// Zero is white in the byte matrix
			if inputY[x] == 1 {
				output.Set(x+margin, yOutput)
			}
		}
	}
	return output, nil
}

// rotateArray takes and rotates the array 90 degrees
func rotateArray(bitarray [][]byte) [][]byte {
	temp := make([][]byte, len(bitarray[0]))
	for i := range temp {
		temp[i] = make([]byte, len(bitarray))
	}
	for ii := 0; ii < len(bitarray); ii++ {
		// This makes the direction consistent on screen when rotating the screen
		inverseii := len(bitarray) - ii - 1
		for jj := 0; jj < len(bitarray[0]); jj++ {
			temp[jj][inverseii] = bitarray[ii][jj]
		}
	}
	return temp
}

// getEncodingFromName gets an encoding from a charset name
func getEncodingFromName(name string) (encoding.Encoding, error) {
	switch name {
	case "ISO-8859-1", "ISO8859_1":
		return charmap.ISO8859_1, nil
	default:
		// Try to get from common character set ECI
		// For now, default to ISO-8859-1
		return charmap.ISO8859_1, nil
	}
}


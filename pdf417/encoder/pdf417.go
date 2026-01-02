/*
 * Copyright 2006 Jeremias Maerki in part, and ZXing Authors in part
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * This file has been modified from its original form in Barcode4J.
 */

package encoder

import (
	"math"

	"github.com/makiuchi-d/gozxing"
	"golang.org/x/text/encoding"
)

const (
	// The start pattern (17 bits)
	START_PATTERN = 0x1fea8
	// The stop pattern (18 bits)
	STOP_PATTERN = 0x3fa29

	PREFERRED_RATIO      = 3.0
	DEFAULT_MODULE_WIDTH = 0.357 // 1px in mm
	HEIGHT               = 2.0   // mm
)

// CODEWORD_TABLE is initialized in codeword_table.go

// PDF417 is the top-level class for the logic part of the PDF417 implementation.
type PDF417 struct {
	barcodeMatrix *BarcodeMatrix
	compact       bool
	compaction    Compaction
	encoding      encoding.Encoding
	minCols       int
	maxCols       int
	maxRows       int
	minRows       int
}

// NewPDF417 creates a new PDF417 encoder
func NewPDF417() *PDF417 {
	return NewPDF417WithCompact(false)
}

// NewPDF417WithCompact creates a new PDF417 encoder with compact mode
func NewPDF417WithCompact(compact bool) *PDF417 {
	return &PDF417{
		compact:    compact,
		compaction: Compaction_AUTO,
		encoding:   nil, // Use default
		minCols:    2,
		maxCols:    30,
		maxRows:    30,
		minRows:    2,
	}
}

// GetBarcodeMatrix returns the barcode matrix
func (p *PDF417) GetBarcodeMatrix() *BarcodeMatrix {
	return p.barcodeMatrix
}

// calculateNumberOfRows calculates the necessary number of rows as described in annex Q of ISO/IEC 15438:2001(E).
func calculateNumberOfRows(m, k, c int) int {
	r := ((m + 1 + k) / c) + 1
	if c*r >= (m + 1 + k + c) {
		r--
	}
	return r
}

// getNumberOfPadCodewords calculates the number of pad codewords as described in 4.9.2 of ISO/IEC 15438:2001(E).
func getNumberOfPadCodewords(m, k, c, r int) int {
	n := c*r - k
	if n > m+1 {
		return n - m - 1
	}
	return 0
}

// encodeChar encodes a character pattern
func encodeChar(pattern, len int, logic *BarcodeRow) {
	map_ := 1 << (len - 1)
	last := (pattern & map_) != 0 // Initialize to inverse of first bit
	width := 0
	for i := 0; i < len; i++ {
		black := (pattern & map_) != 0
		if last == black {
			width++
		} else {
			logic.AddBar(last, width)
			last = black
			width = 1
		}
		map_ >>= 1
	}
	logic.AddBar(last, width)
}

// encodeLowLevel performs low-level encoding
func (p *PDF417) encodeLowLevel(fullCodewords string, c, r, errorCorrectionLevel int, logic *BarcodeMatrix) {
	idx := 0
	for y := 0; y < r; y++ {
		cluster := y % 3
		logic.StartRow()
		encodeChar(START_PATTERN, 17, logic.GetCurrentRow())

		var left, right int
		if cluster == 0 {
			left = (30 * (y / 3)) + ((r - 1) / 3)
			right = (30 * (y / 3)) + (c - 1)
		} else if cluster == 1 {
			left = (30 * (y / 3)) + (errorCorrectionLevel * 3) + ((r - 1) % 3)
			right = (30 * (y / 3)) + ((r - 1) / 3)
		} else {
			left = (30 * (y / 3)) + (c - 1)
			right = (30 * (y / 3)) + (errorCorrectionLevel * 3) + ((r - 1) % 3)
		}

		pattern := CODEWORD_TABLE[cluster][left]
		encodeChar(pattern, 17, logic.GetCurrentRow())

		for x := 0; x < c; x++ {
			pattern = CODEWORD_TABLE[cluster][int(fullCodewords[idx])]
			encodeChar(pattern, 17, logic.GetCurrentRow())
			idx++
		}

		if p.compact {
			encodeChar(STOP_PATTERN, 1, logic.GetCurrentRow()) // encodes stop line for compact pdf417
		} else {
			pattern = CODEWORD_TABLE[cluster][right]
			encodeChar(pattern, 17, logic.GetCurrentRow())
			encodeChar(STOP_PATTERN, 18, logic.GetCurrentRow())
		}
	}
}

// GenerateBarcodeLogic generates the barcode logic
func (p *PDF417) GenerateBarcodeLogic(msg string, errorCorrectionLevel int) error {
	return p.GenerateBarcodeLogicWithAutoECI(msg, errorCorrectionLevel, false)
}

// GenerateBarcodeLogicWithAutoECI generates the barcode logic with autoECI support
func (p *PDF417) GenerateBarcodeLogicWithAutoECI(msg string, errorCorrectionLevel int, autoECI bool) error {
	// 1. step: High-level encoding
	errorCorrectionCodeWords, err := GetErrorCorrectionCodewordCount(errorCorrectionLevel)
	if err != nil {
		return err
	}
	highLevel, err := EncodeHighLevel(msg, p.compaction, p.encoding, autoECI)
	if err != nil {
		return err
	}
	sourceCodeWords := len(highLevel)

	dimension, err := determineDimensions(p.minCols, p.maxCols, p.minRows, p.maxRows, sourceCodeWords, errorCorrectionCodeWords)
	if err != nil {
		return err
	}

	cols := dimension[0]
	rows := dimension[1]

	pad := getNumberOfPadCodewords(sourceCodeWords, errorCorrectionCodeWords, cols, rows)

	// 2. step: construct data codewords
	if sourceCodeWords+errorCorrectionCodeWords+1 > 929 { // +1 for symbol length CW
		return gozxing.NewWriterException("Encoded message contains too many code words, message too big (%d bytes)", len(msg))
	}
	n := sourceCodeWords + pad + 1
	var dataCodewords []rune
	dataCodewords = append(dataCodewords, rune(n))
	dataCodewords = append(dataCodewords, []rune(highLevel)...)
	for i := 0; i < pad; i++ {
		dataCodewords = append(dataCodewords, 900) // PAD characters
	}
	dataCodewordsStr := string(dataCodewords)

	// 3. step: Error correction
	ec, err := GenerateErrorCorrection(dataCodewordsStr, errorCorrectionLevel)
	if err != nil {
		return err
	}

	// 4. step: low-level encoding
	p.barcodeMatrix = NewBarcodeMatrix(rows, cols)
	p.encodeLowLevel(dataCodewordsStr+ec, cols, rows, errorCorrectionLevel, p.barcodeMatrix)
	return nil
}

// determineDimensions determines optimal number of columns and rows for the specified number of codewords.
func determineDimensions(minCols, maxCols, minRows, maxRows, sourceCodeWords, errorCorrectionCodeWords int) ([2]int, error) {
	var ratio float64
	var dimension [2]int
	currentCol := minCols
	found := false

	for cols := minCols; cols <= maxCols; cols++ {
		currentCol = cols
		rows := calculateNumberOfRows(sourceCodeWords, errorCorrectionCodeWords, cols)

		if rows < minRows {
			break
		}

		if rows > maxRows {
			continue
		}

		newRatio := (float64(17*cols+69) * DEFAULT_MODULE_WIDTH) / (float64(rows) * HEIGHT)

		// ignore if previous ratio is closer to preferred ratio
		if found && math.Abs(newRatio-PREFERRED_RATIO) > math.Abs(ratio-PREFERRED_RATIO) {
			continue
		}

		ratio = newRatio
		dimension = [2]int{cols, rows}
		found = true
	}

	// Handle case when min values were larger than necessary
	if !found {
		rows := calculateNumberOfRows(sourceCodeWords, errorCorrectionCodeWords, currentCol)
		if rows < minRows {
			dimension = [2]int{minCols, minRows}
			found = true
		}
	}

	if !found {
		return [2]int{}, gozxing.NewWriterException("Unable to fit message in columns")
	}

	return dimension, nil
}

// SetDimensions sets max/min row/col values
func (p *PDF417) SetDimensions(maxCols, minCols, maxRows, minRows int) {
	p.maxCols = maxCols
	p.minCols = minCols
	p.maxRows = maxRows
	p.minRows = minRows
}

// SetCompaction sets the compaction mode to use
func (p *PDF417) SetCompaction(compaction Compaction) {
	p.compaction = compaction
}

// SetCompact sets if compact mode should be enabled
func (p *PDF417) SetCompact(compact bool) {
	p.compact = compact
}

// SetEncoding sets character encoding to use
func (p *PDF417) SetEncoding(enc encoding.Encoding) {
	p.encoding = enc
}

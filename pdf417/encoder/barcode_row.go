/*
 * Copyright 2011 ZXing authors
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

package encoder

// BarcodeRow represents a single row of a PDF417 barcode
type BarcodeRow struct {
	row             []byte
	currentLocation int
}

// NewBarcodeRow creates a BarcodeRow of the specified width
func NewBarcodeRow(width int) *BarcodeRow {
	return &BarcodeRow{
		row:             make([]byte, width),
		currentLocation: 0,
	}
}

// Set sets a specific location in the bar
func (br *BarcodeRow) Set(x int, value byte) {
	br.row[x] = value
}

// AddBar adds a bar to the row
// black: true if the bar is black, false if white
// width: how many spots wide the bar is
func (br *BarcodeRow) AddBar(black bool, width int) {
	value := byte(0)
	if black {
		value = 1
	}
	for i := 0; i < width; i++ {
		br.row[br.currentLocation] = value
		br.currentLocation++
	}
}

// GetScaledRow scales the row by the specified factor
// scale: how much you want the image to be scaled, must be >= 1
func (br *BarcodeRow) GetScaledRow(scale int) []byte {
	output := make([]byte, len(br.row)*scale)
	for i := 0; i < len(output); i++ {
		output[i] = br.row[i/scale]
	}
	return output
}

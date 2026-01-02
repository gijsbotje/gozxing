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

// BarcodeMatrix holds all of the information for a barcode in a format where it can be easily accessible
type BarcodeMatrix struct {
	matrix     []*BarcodeRow
	currentRow int
	height     int
	width      int
}

// NewBarcodeMatrix creates a new BarcodeMatrix
// height: the height of the matrix (Rows)
// width: the width of the matrix (Cols)
func NewBarcodeMatrix(height, width int) *BarcodeMatrix {
	matrix := make([]*BarcodeRow, height)
	// Initialize the array to the correct width
	for i := 0; i < height; i++ {
		matrix[i] = NewBarcodeRow((width+4)*17 + 1)
	}
	return &BarcodeMatrix{
		matrix:     matrix,
		currentRow: -1,
		height:     height,
		width:      width * 17,
	}
}

// Set sets a value at a specific position
func (bm *BarcodeMatrix) Set(x, y int, value byte) {
	bm.matrix[y].Set(x, value)
}

// StartRow starts a new row
func (bm *BarcodeMatrix) StartRow() {
	bm.currentRow++
}

// GetCurrentRow returns the current row
func (bm *BarcodeMatrix) GetCurrentRow() *BarcodeRow {
	return bm.matrix[bm.currentRow]
}

// GetMatrix returns the matrix as a 2D byte array
func (bm *BarcodeMatrix) GetMatrix() [][]byte {
	return bm.GetScaledMatrix(1, 1)
}

// GetScaledMatrix returns the matrix scaled by the specified factors
func (bm *BarcodeMatrix) GetScaledMatrix(xScale, yScale int) [][]byte {
	matrixOut := make([][]byte, bm.height*yScale)
	yMax := bm.height * yScale
	for i := 0; i < yMax; i++ {
		matrixOut[yMax-i-1] = bm.matrix[i/yScale].GetScaledRow(xScale)
	}
	return matrixOut
}

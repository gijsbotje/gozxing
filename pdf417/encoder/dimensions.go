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

package encoder

// Dimensions is a data object to specify the minimum and maximum number of rows and columns for a PDF417 barcode.
type Dimensions struct {
	minCols int
	maxCols int
	minRows int
	maxRows int
}

func NewDimensions(minCols, maxCols, minRows, maxRows int) *Dimensions {
	return &Dimensions{
		minCols: minCols,
		maxCols: maxCols,
		minRows: minRows,
		maxRows: maxRows,
	}
}

func (d *Dimensions) GetMinCols() int {
	return d.minCols
}

func (d *Dimensions) GetMaxCols() int {
	return d.maxCols
}

func (d *Dimensions) GetMinRows() int {
	return d.minRows
}

func (d *Dimensions) GetMaxRows() int {
	return d.maxRows
}


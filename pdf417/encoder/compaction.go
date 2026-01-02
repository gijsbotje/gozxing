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

import "github.com/makiuchi-d/gozxing"

// Compaction represents possible PDF417 barcode compaction types.
type Compaction int

const (
	Compaction_AUTO Compaction = iota
	Compaction_TEXT
	Compaction_BYTE
	Compaction_NUMERIC
)

// ValueOf returns the Compaction enum value for the given string
func Compaction_ValueOf(s string) (Compaction, error) {
	switch s {
	case "AUTO":
		return Compaction_AUTO, nil
	case "TEXT":
		return Compaction_TEXT, nil
	case "BYTE":
		return Compaction_BYTE, nil
	case "NUMERIC":
		return Compaction_NUMERIC, nil
	default:
		return Compaction_AUTO, gozxing.NewWriterException("Unknown compaction: %v", s)
	}
}

func (c Compaction) String() string {
	switch c {
	case Compaction_AUTO:
		return "AUTO"
	case Compaction_TEXT:
		return "TEXT"
	case Compaction_BYTE:
		return "BYTE"
	case Compaction_NUMERIC:
		return "NUMERIC"
	default:
		return "AUTO"
	}
}

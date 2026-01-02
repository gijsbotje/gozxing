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
	"math/big"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/common"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

const (
	// code for Text compaction
	TEXT_COMPACTION = 0
	// code for Byte compaction
	BYTE_COMPACTION = 1
	// code for Numeric compaction
	NUMERIC_COMPACTION = 2

	// Text compaction submode Alpha
	SUBMODE_ALPHA = 0
	// Text compaction submode Lower
	SUBMODE_LOWER = 1
	// Text compaction submode Mixed
	SUBMODE_MIXED = 2
	// Text compaction submode Punctuation
	SUBMODE_PUNCTUATION = 3

	// mode latch to Text Compaction mode
	LATCH_TO_TEXT = 900
	// mode latch to Byte Compaction mode (number of characters NOT a multiple of 6)
	LATCH_TO_BYTE_PADDED = 901
	// mode latch to Numeric Compaction mode
	LATCH_TO_NUMERIC = 902
	// mode shift to Byte Compaction mode
	SHIFT_TO_BYTE = 913
	// mode latch to Byte Compaction mode (number of characters a multiple of 6)
	LATCH_TO_BYTE = 924
	// identifier for a user defined Extended Channel Interpretation (ECI)
	ECI_USER_DEFINED = 925
	// identifier for a general purpose ECO format
	ECI_GENERAL_PURPOSE = 926
	// identifier for an ECI of a character set of code page
	ECI_CHARSET = 927
)

// Raw code table for text compaction Mixed sub-mode
var TEXT_MIXED_RAW = []byte{
	48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 38, 13, 9, 44, 58,
	35, 45, 46, 36, 47, 43, 37, 42, 61, 94, 0, 32, 0, 0, 0,
}

// Raw code table for text compaction: Punctuation sub-mode
var TEXT_PUNCTUATION_RAW = []byte{
	59, 60, 62, 64, 91, 92, 93, 95, 96, 126, 33, 13, 9, 44, 58,
	10, 45, 46, 36, 47, 34, 124, 42, 40, 41, 63, 123, 125, 39, 0,
}

var MIXED = make([]byte, 128)
var PUNCTUATION = make([]byte, 128)

var DEFAULT_ENCODING = charmap.ISO8859_1

func init() {
	// Construct inverse lookups
	for i := range MIXED {
		MIXED[i] = 255 // -1
	}
	for i, b := range TEXT_MIXED_RAW {
		if b > 0 {
			MIXED[b] = byte(i)
		}
	}
	for i := range PUNCTUATION {
		PUNCTUATION[i] = 255 // -1
	}
	for i, b := range TEXT_PUNCTUATION_RAW {
		if b > 0 {
			PUNCTUATION[b] = byte(i)
		}
	}
}

// EncodeHighLevel performs high-level encoding of a PDF417 message using the algorithm described in annex P
// of ISO/IEC 15438:2001(E). If byte compaction has been selected, then only byte compaction is used.
func EncodeHighLevel(msg string, compaction Compaction, charset encoding.Encoding, autoECI bool) (string, error) {
	if len(msg) == 0 {
		return "", gozxing.NewWriterException("Empty message not allowed")
	}

	if compaction == Compaction_TEXT {
		if err := checkCharset(msg, 127, "Consider specifying Compaction.AUTO instead of Compaction.TEXT"); err != nil {
			return "", err
		}
	}

	if charset == nil && !autoECI {
		if err := checkCharset(msg, 255, "Consider specifying EncodeHintType.PDF417_AUTO_ECI and/or EncodeTypeHint.CHARACTER_SET"); err != nil {
			return "", err
		}
	}

	// the codewords 0..928 are encoded as Unicode characters
	var sb []rune

	if charset == nil {
		charset = DEFAULT_ENCODING
	} else if charset != DEFAULT_ENCODING {
		eci, ok := common.GetCharacterSetECI(charset)
		if ok && eci != nil {
			encodingECI(eci.GetValue(), &sb)
		}
	}

	msgLen := len(msg)
	p := 0
	textSubMode := SUBMODE_ALPHA

	// User selected encoding mode
	switch compaction {
	case Compaction_TEXT:
		encodeText(msg, p, msgLen, &sb, textSubMode)
	case Compaction_BYTE:
		msgBytes, err := charset.NewEncoder().Bytes([]byte(msg))
		if err != nil {
			return "", gozxing.WrapWriterException(err)
		}
		encodeBinary(msgBytes, p, len(msgBytes), BYTE_COMPACTION, &sb)
	case Compaction_NUMERIC:
		sb = append(sb, LATCH_TO_NUMERIC)
		encodeNumeric(msg, p, msgLen, &sb)
	default: // AUTO
		encodingMode := TEXT_COMPACTION // Default mode, see 4.4.2.1
		for p < msgLen {
			n := determineConsecutiveDigitCount(msg, p)
			if n >= 13 {
				sb = append(sb, LATCH_TO_NUMERIC)
				encodingMode = NUMERIC_COMPACTION
				textSubMode = SUBMODE_ALPHA // Reset after latch
				encodeNumeric(msg, p, n, &sb)
				p += n
			} else {
				t := determineConsecutiveTextCount(msg, p)
				if t >= 5 || n == msgLen {
					if encodingMode != TEXT_COMPACTION {
						sb = append(sb, LATCH_TO_TEXT)
						encodingMode = TEXT_COMPACTION
						textSubMode = SUBMODE_ALPHA // start with submode alpha after latch
					}
					textSubMode = encodeText(msg, p, t, &sb, textSubMode)
					p += t
				} else {
					b := determineConsecutiveBinaryCount(msg, p, charset)
					if b == 0 {
						b = 1
					}
					bytes, err := charset.NewEncoder().Bytes([]byte(msg[p : p+b]))
					if err != nil {
						return "", gozxing.WrapWriterException(err)
					}
					if len(bytes) == 1 && encodingMode == TEXT_COMPACTION {
						// Switch for one byte (instead of latch)
						encodeBinary(bytes, 0, 1, TEXT_COMPACTION, &sb)
					} else {
						// Mode latch performed by encodeBinary()
						encodeBinary(bytes, 0, len(bytes), encodingMode, &sb)
						encodingMode = BYTE_COMPACTION
						textSubMode = SUBMODE_ALPHA // Reset after latch
					}
					p += b
				}
			}
		}
	}

	return string(sb), nil
}

func checkCharset(input string, max int, errorMessage string) error {
	for i, ch := range input {
		if ch > rune(max) {
			return gozxing.NewWriterException("Non-encodable character detected: %c (Unicode: %d) at position #%d - %s", ch, int(ch), i, errorMessage)
		}
	}
	return nil
}

// encodeText encodes parts of the message using Text Compaction as described in ISO/IEC 15438:2001(E), chapter 4.4.2.
func encodeText(input string, startpos, count int, sb *[]rune, initialSubmode int) int {
	var tmp []rune
	submode := initialSubmode
	idx := 0
	for idx < count {
		ch := rune(input[startpos+idx])
		switch submode {
		case SUBMODE_ALPHA:
			if isAlphaUpper(ch) {
				if ch == ' ' {
					tmp = append(tmp, 26) // space
				} else {
					tmp = append(tmp, ch-65)
				}
			} else {
				if isAlphaLower(ch) {
					submode = SUBMODE_LOWER
					tmp = append(tmp, 27) // ll
					continue
				} else if isMixed(ch) {
					submode = SUBMODE_MIXED
					tmp = append(tmp, 28) // ml
					continue
				} else {
					tmp = append(tmp, 29) // ps
					tmp = append(tmp, rune(PUNCTUATION[ch]))
				}
			}
		case SUBMODE_LOWER:
			if isAlphaLower(ch) {
				if ch == ' ' {
					tmp = append(tmp, 26) // space
				} else {
					tmp = append(tmp, ch-97)
				}
			} else {
				if isAlphaUpper(ch) {
					tmp = append(tmp, 27) // as
					tmp = append(tmp, ch-65)
					// space cannot happen here, it is also in "Lower"
				} else if isMixed(ch) {
					submode = SUBMODE_MIXED
					tmp = append(tmp, 28) // ml
					continue
				} else {
					tmp = append(tmp, 29) // ps
					tmp = append(tmp, rune(PUNCTUATION[ch]))
				}
			}
		case SUBMODE_MIXED:
			if isMixed(ch) {
				tmp = append(tmp, rune(MIXED[ch]))
			} else {
				if isAlphaUpper(ch) {
					submode = SUBMODE_ALPHA
					tmp = append(tmp, 28) // al
					continue
				} else if isAlphaLower(ch) {
					submode = SUBMODE_LOWER
					tmp = append(tmp, 27) // ll
					continue
				} else {
					if startpos+idx+1 < count && isPunctuation(rune(input[startpos+idx+1])) {
						submode = SUBMODE_PUNCTUATION
						tmp = append(tmp, 25) // pl
						continue
					}
					tmp = append(tmp, 29) // ps
					tmp = append(tmp, rune(PUNCTUATION[ch]))
				}
			}
		default: // SUBMODE_PUNCTUATION
			if isPunctuation(ch) {
				tmp = append(tmp, rune(PUNCTUATION[ch]))
			} else {
				submode = SUBMODE_ALPHA
				tmp = append(tmp, 29) // al
				continue
			}
		}
		idx++
	}
	var h rune
	tmpLen := len(tmp)
	for i := 0; i < tmpLen; i++ {
		odd := (i % 2) != 0
		if odd {
			h = (h * 30) + tmp[i]
			*sb = append(*sb, h)
		} else {
			h = tmp[i]
		}
	}
	if (tmpLen % 2) != 0 {
		*sb = append(*sb, (h*30)+29) // ps
	}
	return submode
}

// encodeBinary encodes parts of the message using Byte Compaction as described in ISO/IEC 15438:2001(E), chapter 4.4.3.
func encodeBinary(bytes []byte, startpos, count int, startmode int, sb *[]rune) {
	if count == 1 && startmode == TEXT_COMPACTION {
		*sb = append(*sb, SHIFT_TO_BYTE)
	} else {
		if (count % 6) == 0 {
			*sb = append(*sb, LATCH_TO_BYTE)
		} else {
			*sb = append(*sb, LATCH_TO_BYTE_PADDED)
		}
	}

	idx := startpos
	// Encode sixpacks
	if count >= 6 {
		chars := make([]rune, 5)
		for (startpos + count - idx) >= 6 {
			var t int64
			for i := 0; i < 6; i++ {
				t <<= 8
				t += int64(bytes[idx+i] & 0xff)
			}
			for i := 0; i < 5; i++ {
				chars[i] = rune(t % 900)
				t /= 900
			}
			for i := len(chars) - 1; i >= 0; i-- {
				*sb = append(*sb, chars[i])
			}
			idx += 6
		}
	}
	// Encode rest (remaining n<5 bytes if any)
	for i := idx; i < startpos+count; i++ {
		ch := bytes[i] & 0xff
		*sb = append(*sb, rune(ch))
	}
}

// encodeNumeric encodes numeric data
func encodeNumeric(input string, startpos, count int, sb *[]rune) {
	idx := 0
	num900 := big.NewInt(900)
	num0 := big.NewInt(0)
	for idx < count {
		segLen := count - idx
		if segLen > 44 {
			segLen = 44
		}
		part := "1" + input[startpos+idx:startpos+idx+segLen]
		bigint := new(big.Int)
		bigint, _ = bigint.SetString(part, 10)
		var tmp []rune
		for bigint.Cmp(num0) != 0 {
			mod := new(big.Int)
			mod.Mod(bigint, num900)
			tmp = append(tmp, rune(mod.Int64()))
			bigint.Div(bigint, num900)
		}
		// Reverse temporary string
		for i := len(tmp) - 1; i >= 0; i-- {
			*sb = append(*sb, tmp[i])
		}
		idx += segLen
	}
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isAlphaUpper(ch rune) bool {
	return ch == ' ' || (ch >= 'A' && ch <= 'Z')
}

func isAlphaLower(ch rune) bool {
	return ch == ' ' || (ch >= 'a' && ch <= 'z')
}

func isMixed(ch rune) bool {
	if ch < 0 || ch >= 128 {
		return false
	}
	return MIXED[ch] != 255
}

func isPunctuation(ch rune) bool {
	if ch < 0 || ch >= 128 {
		return false
	}
	return PUNCTUATION[ch] != 255
}

func isText(ch rune) bool {
	return ch == '\t' || ch == '\n' || ch == '\r' || (ch >= 32 && ch <= 126)
}

// determineConsecutiveDigitCount determines the number of consecutive characters that are encodable using numeric compaction.
func determineConsecutiveDigitCount(input string, startpos int) int {
	count := 0
	len := len(input)
	idx := startpos
	for idx < len && isDigit(rune(input[idx])) {
		count++
		idx++
	}
	return count
}

// determineConsecutiveTextCount determines the number of consecutive characters that are encodable using text compaction.
func determineConsecutiveTextCount(input string, startpos int) int {
	inputLen := len(input)
	idx := startpos
	for idx < inputLen {
		numericCount := 0
		for numericCount < 13 && idx < inputLen && isDigit(rune(input[idx])) {
			numericCount++
			idx++
		}
		if numericCount >= 13 {
			return idx - startpos - numericCount
		}
		if numericCount > 0 {
			// Heuristic: All text-encodable chars or digits are binary encodable
			continue
		}
		// Check if character is encodable
		if !isText(rune(input[idx])) {
			break
		}
		idx++
	}
	return idx - startpos
}

// determineConsecutiveBinaryCount determines the number of consecutive characters that are encodable using binary compaction.
func determineConsecutiveBinaryCount(input string, startpos int, enc encoding.Encoding) int {
	inputLen := len(input)
	idx := startpos
	encoder := enc.NewEncoder()
	for idx < inputLen {
		numericCount := 0
		i := idx
		for numericCount < 13 && i < inputLen && isDigit(rune(input[i])) {
			numericCount++
			i = idx + numericCount
			if i >= inputLen {
				break
			}
		}
		if numericCount >= 13 {
			return idx - startpos
		}
		// Check if character can be encoded - try to encode a single byte
		testBytes, err := encoder.Bytes([]byte{byte(input[idx])})
		if err != nil || len(testBytes) == 0 {
			ch := rune(input[idx])
			panic(gozxing.NewWriterException("Non-encodable character detected: %c (Unicode: %d)", ch, int(ch)))
		}
		idx++
	}
	return idx - startpos
}

func encodingECI(eci int, sb *[]rune) error {
	if eci >= 0 && eci < 900 {
		*sb = append(*sb, ECI_CHARSET)
		*sb = append(*sb, rune(eci))
	} else if eci < 810900 {
		*sb = append(*sb, ECI_GENERAL_PURPOSE)
		*sb = append(*sb, rune(eci/900-1))
		*sb = append(*sb, rune(eci%900))
	} else if eci < 811800 {
		*sb = append(*sb, ECI_USER_DEFINED)
		*sb = append(*sb, rune(810900-eci))
	} else {
		return gozxing.NewWriterException("ECI number not in valid range from 0..811799, but was %d", eci)
	}
	return nil
}

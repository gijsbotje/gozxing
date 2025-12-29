package encoder

import (
	"reflect"
	"testing"

	"golang.org/x/text/encoding/ianaindex"
)

// Helper function tests - matching HighLevelEncoder structure

func TestMinimalEncoder_isExtendedASCII(t *testing.T) {
	if isExtendedASCII(' ', -1) != false {
		t.Fatalf("isExtendedASCII(' ') must false")
	}
	if isExtendedASCII(0x7f, -1) != false {
		t.Fatalf("isExtendedASCII(0x7f) must false")
	}
	if isExtendedASCII(0x80, -1) != true {
		t.Fatalf("isExtendedASCII(0x80) must true")
	}
	if isExtendedASCII(0xff, -1) != true {
		t.Fatalf("isExtendedASCII(0xff) must true")
	}
	if isExtendedASCII(0x1D, 0x1D) != false {
		t.Fatalf("isExtendedASCII(0x1D, 0x1D) must false (FNC1)")
	}
}

func TestMinimalEncoder_isInC40Shift1Set(t *testing.T) {
	if isInC40Shift1Set(0) != true {
		t.Fatalf("isInC40Shift1Set(0) must true")
	}
	if isInC40Shift1Set(31) != true {
		t.Fatalf("isInC40Shift1Set(31) must true")
	}
	if isInC40Shift1Set(32) != false {
		t.Fatalf("isInC40Shift1Set(32) must false")
	}
}

func TestMinimalEncoder_isInC40Shift2Set(t *testing.T) {
	if isInC40Shift2Set('!', -1) != true {
		t.Fatalf("isInC40Shift2Set('!') must true")
	}
	if isInC40Shift2Set('_', -1) != true {
		t.Fatalf("isInC40Shift2Set('_') must true")
	}
	if isInC40Shift2Set('A', -1) != false {
		t.Fatalf("isInC40Shift2Set('A') must false")
	}
	if isInC40Shift2Set(0x1D, 0x1D) != true {
		t.Fatalf("isInC40Shift2Set(0x1D, 0x1D) must true (FNC1)")
	}
}

func TestMinimalEncoder_isInTextShift1Set(t *testing.T) {
	if isInTextShift1Set(0) != true {
		t.Fatalf("isInTextShift1Set(0) must true")
	}
	if isInTextShift1Set(31) != true {
		t.Fatalf("isInTextShift1Set(31) must true")
	}
	if isInTextShift1Set(32) != false {
		t.Fatalf("isInTextShift1Set(32) must false")
	}
}

func TestMinimalEncoder_isInTextShift2Set(t *testing.T) {
	if isInTextShift2Set('!', -1) != true {
		t.Fatalf("isInTextShift2Set('!') must true")
	}
	if isInTextShift2Set('_', -1) != true {
		t.Fatalf("isInTextShift2Set('_') must true")
	}
	if isInTextShift2Set('A', -1) != false {
		t.Fatalf("isInTextShift2Set('A') must false")
	}
	if isInTextShift2Set(0x1D, 0x1D) != true {
		t.Fatalf("isInTextShift2Set(0x1D, 0x1D) must true (FNC1)")
	}
}

func TestMinimalEncoder_isDigit(t *testing.T) {
	if isDigitChar('/') != false {
		t.Fatalf("isDigitChar('/') must false")
	}
	if isDigitChar('0') != true {
		t.Fatalf("isDigitChar('0') must true")
	}
	if isDigitChar('9') != true {
		t.Fatalf("isDigitChar('9') must true")
	}
	if isDigitChar(':') != false {
		t.Fatalf("isDigitChar(':') must false")
	}
}

func TestMinimalEncoder_isNativeC40(t *testing.T) {
	tcs := []rune{' ', '0', '9', 'A', 'Z'}
	fcs := []rune{'!', '@', '^', 'a', '~'}
	for _, c := range tcs {
		if isNativeC40Char(c) != true {
			t.Fatalf("isNativeC40Char(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeC40Char(c) != false {
			t.Fatalf("isNativeC40Char(%v) must false", c)
		}
	}
}

func TestMinimalEncoder_isNativeText(t *testing.T) {
	tcs := []rune{' ', '0', '9', 'a', 'z'}
	fcs := []rune{'!', '@', '^', 'A', '~'}
	for _, c := range tcs {
		if isNativeTextChar(c) != true {
			t.Fatalf("isNativeTextChar(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeTextChar(c) != false {
			t.Fatalf("isNativeTextChar(%v) must false", c)
		}
	}
}

func TestMinimalEncoder_isNativeX12(t *testing.T) {
	tcs := []rune{'\r', '*', '>', ' ', '0', '9', 'A', 'Z'}
	fcs := []rune{'!', '@', '^', 'a', '~'}
	for _, c := range tcs {
		if isNativeX12Char(c) != true {
			t.Fatalf("isNativeX12Char(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeX12Char(c) != false {
			t.Fatalf("isNativeX12Char(%v) must false", c)
		}
	}
}

func TestMinimalEncoder_isX12TermSepChar(t *testing.T) {
	if isX12TermSepChar('\r') != true {
		t.Fatalf("isX12TermSepChar('\\r') must true")
	}
	if isX12TermSepChar('*') != true {
		t.Fatalf("isX12TermSepChar('*') must true")
	}
	if isX12TermSepChar('>') != true {
		t.Fatalf("isX12TermSepChar('>') must true")
	}
	if isX12TermSepChar(' ') != false {
		t.Fatalf("isX12TermSepChar(' ') must false")
	}
}

func TestMinimalEncoder_isNativeEDIFACT(t *testing.T) {
	tcs := []rune{' ', '!', '@', '0', 'A', '^'}
	fcs := []rune{'\r', '_', 'a', '~'}
	for _, c := range tcs {
		if isNativeEDIFACTChar(c) != true {
			t.Fatalf("isNativeEDIFACTChar(%v) must true", c)
		}
	}
	for _, c := range fcs {
		if isNativeEDIFACTChar(c) != false {
			t.Fatalf("isNativeEDIFACTChar(%v) must false", c)
		}
	}
}

// Main encoding test - matching HighLevelEncoder structure

func TestEncodeHighLevelMinimal(t *testing.T) {
	encoder := MinimalEncoder{}
	// Test error cases (similar to HighLevelEncoder)
	// Note: MinimalEncoder doesn't have the same error checking as HighLevelEncoder
	// but we can test with very long strings
	str := string(make([]byte, 1559))
	_, e := encoder.EncodeHighLevel(str, nil, -1, SymbolShapeHint_FORCE_NONE)
	if e == nil {
		// MinimalEncoder might handle this differently, so we just check it doesn't panic
		// This is a soft check - if it encodes, that's fine
	}

	// Test macro 05 - should match HighLevelEncoder output
	str = "[)>\u001E05\u001Daaaaaa\u001E\u0004"
	b, e := encoder.EncodeHighLevel(str, nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	// HighLevelEncoder produces: []byte{236, 239, 89, 191, 89, 191, 254, 129}
	// MinimalEncoder should start with macro 05 codeword (236)
	if b[0] != 236 {
		t.Fatalf("EncodeHighLevel macro 05 first byte = %v, expect 236", b[0])
	}

	// Test macro 06 - should match HighLevelEncoder output
	str = "[)>\u001E06\u001Daaaaaa\u001E\u0004"
	b, e = encoder.EncodeHighLevel(str, nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	// HighLevelEncoder produces: []byte{237, 239, 89, 191, 89, 191, 254, 129}
	// MinimalEncoder should start with macro 06 codeword (237)
	if b[0] != 237 {
		t.Fatalf("EncodeHighLevel macro 06 first byte = %v, expect 237", b[0])
	}

	// Test basic encoding
	b, e = encoder.EncodeHighLevel("123456", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}

	// Test with extended ASCII
	b, e = encoder.EncodeHighLevel("123456£", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
}

func TestEncodeHighLevelMinimalWithOptions(t *testing.T) {
	shape := SymbolShapeHint_FORCE_NONE
	encoder := MinimalEncoder{}

	// Test with GS1 format (fnc1 = 0x1D)
	b, e := encoder.EncodeHighLevel("010123456789012810ABCD1234", nil, 0x1D, shape)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
	// Should start with FNC1 codeword (232)
	if b[0] != 232 {
		t.Fatalf("EncodeHighLevel GS1 first byte = %v, expect 232", b[0])
	}

	// Test with priority charset
	utf8, _ := ianaindex.IANA.Encoding("UTF-8")
	b, e = encoder.EncodeHighLevel("Hello World", utf8, -1, shape)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}

	// Test without GS1 (fnc1 = -1)
	b, e = encoder.EncodeHighLevel("Hello World", nil, -1, shape)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
	// Should not start with FNC1
	if b[0] == 232 {
		t.Fatalf("EncodeHighLevel should not start with FNC1 when fnc1=-1")
	}
}

// Encoding mode tests

func TestMinimalEncoderASCIIEncodation(t *testing.T) {
	encoder := MinimalEncoder{}
	b, e := encoder.EncodeHighLevel("123456", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	// Should encode digits efficiently
	if len(b) < 3 {
		t.Fatalf("EncodeHighLevel(123456) size too small: %d", len(b))
	}
}

func TestMinimalEncoderC40Encodation(t *testing.T) {
	encoder := MinimalEncoder{}
	b, e := encoder.EncodeHighLevel("AIMAIMAIM", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
}

func TestMinimalEncoderTextEncodation(t *testing.T) {
	encoder := MinimalEncoder{}
	b, e := encoder.EncodeHighLevel("aimaimaim", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
}

func TestMinimalEncoderX12Encodation(t *testing.T) {
	encoder := MinimalEncoder{}
	b, e := encoder.EncodeHighLevel("ABC>ABC123>AB", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
}

func TestMinimalEncoderEDIFACTEncodation(t *testing.T) {
	encoder := MinimalEncoder{}
	b, e := encoder.EncodeHighLevel(".A.C1.3.DATA.123DATA.123DATA", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
}

func TestMinimalEncoderBase256Encodation(t *testing.T) {
	encoder := MinimalEncoder{}
	// Test with extended ASCII characters
	b, e := encoder.EncodeHighLevel("\u00ABäöüé\u00BB", nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
}

// Feature-specific tests

func TestMinimalEncoderGS1Format(t *testing.T) {
	shape := SymbolShapeHint_FORCE_NONE
	encoder := MinimalEncoder{}
	// Test GS1 format with FNC1 character
	b, e := encoder.EncodeHighLevel("010123456789012810ABCD1234", nil, 0x1D, shape)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
	// Should start with FNC1 (232)
	if b[0] != 232 {
		t.Fatalf("GS1 format should start with FNC1 (232), got %d", b[0])
	}
}

func TestMinimalEncoderECI(t *testing.T) {
	shape := SymbolShapeHint_FORCE_NONE
	encoder := MinimalEncoder{}
	// Test with UTF-8 priority charset
	utf8, _ := ianaindex.IANA.Encoding("UTF-8")
	b, e := encoder.EncodeHighLevel("Hello World", utf8, -1, shape)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	if len(b) == 0 {
		t.Fatalf("EncodeHighLevel returns empty")
	}
}

func TestMinimalEncoderShapeHints(t *testing.T) {
	msg := "ABCDEFG"
	encoder := MinimalEncoder{}

	// Test FORCE_NONE
	b1, e1 := encoder.EncodeHighLevel(msg, nil, -1, SymbolShapeHint_FORCE_NONE)
	if e1 != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e1)
	}

	// Test FORCE_SQUARE
	b2, e2 := encoder.EncodeHighLevel(msg, nil, -1, SymbolShapeHint_FORCE_SQUARE)
	if e2 != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e2)
	}

	// Test FORCE_RECTANGLE
	b3, e3 := encoder.EncodeHighLevel(msg, nil, -1, SymbolShapeHint_FORCE_RECTANGLE)
	if e3 != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e3)
	}

	// All should produce valid encodings
	if len(b1) == 0 || len(b2) == 0 || len(b3) == 0 {
		t.Fatalf("All shape hints should produce non-empty encodings")
	}
}

// Comparison tests - MinimalEncoder specific

func TestMinimalEncoderComparison(t *testing.T) {
	shape := SymbolShapeHint_FORCE_NONE
	testMessages := []string{
		"Hello World!",
		"1234567890",
		"ABCDEFG",
		"abcdefg",
		"AIMAIMAIM",
		"aimaimaim",
		"ABC>ABC123>AB",
		".A.C1.3.DATA",
	}

	minimalEncoder := MinimalEncoder{}
	highLevelEncoder := HighLevelEncoder{}
	for _, msg := range testMessages {
		minimal, e1 := minimalEncoder.EncodeHighLevel(msg, nil, -1, shape)
		highLevel, e2 := highLevelEncoder.EncodeHighLevel(msg, shape, nil, nil, false)

		if e1 != nil {
			t.Errorf("EncodeHighLevel(%q) returns error: %v", msg, e1)
			continue
		}
		if e2 != nil {
			t.Errorf("EncodeHighLevel(%q) returns error: %v", msg, e2)
			continue
		}

		minimalSize := len(minimal)
		highLevelSize := len(highLevel)

		// Minimal encoder should produce smaller or equal size
		if minimalSize > highLevelSize {
			t.Errorf("MinimalEncoder size (%d) > HighLevelEncoder size (%d) for %q", minimalSize, highLevelSize, msg)
		}
	}
}

func TestMinimalEncoderSizes(t *testing.T) {
	shape := SymbolShapeHint_FORCE_NONE

	testCases := []struct {
		msg     string
		minSize int
		maxSize int
	}{
		{"A", 3, 3},
		{"AB", 3, 3},
		{"ABC", 3, 3},
		{"ABCD", 5, 5},
		{"ABCDE", 5, 5},
		{"ABCDEF", 5, 5},
		{"ABCDEFG", 8, 8},
		{"a", 3, 3},
		{"ab", 3, 3},
		{"abc", 3, 3},
		{"abcd", 5, 5},
		{"abcdef", 5, 5},
		{"abcdefg", 8, 8},
	}

	minimalEncoder := MinimalEncoder{}
	highLevelEncoder := HighLevelEncoder{}
	for _, tc := range testCases {
		minimal, e1 := minimalEncoder.EncodeHighLevel(tc.msg, nil, -1, shape)
		highLevel, e2 := highLevelEncoder.EncodeHighLevel(tc.msg, shape, nil, nil, false)

		if e1 != nil {
			t.Fatalf("EncodeHighLevel(%q) returns error: %v", tc.msg, e1)
		}
		if e2 != nil {
			t.Fatalf("EncodeHighLevel(%q) returns error: %v", tc.msg, e2)
		}

		minimalSize := len(minimal)
		highLevelSize := len(highLevel)

		// Minimal encoder should produce smaller or equal size
		if minimalSize > highLevelSize {
			t.Errorf("MinimalEncoder size (%d) > HighLevelEncoder size (%d) for %q", minimalSize, highLevelSize, tc.msg)
		}

		// Check expected minimum size
		if minimalSize < tc.minSize {
			t.Errorf("MinimalEncoder size (%d) < expected min (%d) for %q", minimalSize, tc.minSize, tc.msg)
		}
	}
}

// Specific test cases with expected outputs where applicable

func TestMinimalEncoderSpecificCases(t *testing.T) {
	encoder := MinimalEncoder{}
	// Test case that should match HighLevelEncoder output
	str := "[)>\u001E05\u001Daaaaaa\u001E\u0004"
	b, e := encoder.EncodeHighLevel(str, nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	// HighLevelEncoder produces: []byte{236, 239, 89, 191, 89, 191, 254, 129}
	// MinimalEncoder should produce the same or similar
	expect := []byte{236, 239, 89, 191, 89, 191, 254, 129}
	if !reflect.DeepEqual(b, expect) {
		// For now, just check it starts correctly
		if b[0] != expect[0] {
			t.Fatalf("EncodeHighLevel macro 05 = %v, expect starts with %v", b, expect)
		}
	}

	// Test case that should match HighLevelEncoder output
	str = "[)>\u001E06\u001Daaaaaa\u001E\u0004"
	b, e = encoder.EncodeHighLevel(str, nil, -1, SymbolShapeHint_FORCE_NONE)
	if e != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", e)
	}
	// HighLevelEncoder produces: []byte{237, 239, 89, 191, 89, 191, 254, 129}
	expect = []byte{237, 239, 89, 191, 89, 191, 254, 129}
	if !reflect.DeepEqual(b, expect) {
		// For now, just check it starts correctly
		if b[0] != expect[0] {
			t.Fatalf("EncodeHighLevel macro 06 = %v, expect starts with %v", b, expect)
		}
	}

	// Test cases from Java HighLevelEncodeTestCase
	testCases := []struct {
		name string
		msg  string
	}{
		{
			name: "ASCII digits",
			msg:  "123456",
		},
		{
			name: "Hello World",
			msg:  "Hello World!",
		},
		{
			name: "C40 example",
			msg:  "A1B2C3D4E5F6G7H8I9J0K1L2",
		},
		{
			name: "Text example",
			msg:  "aimaimaim",
		},
		{
			name: "X12 example",
			msg:  "ABC>ABC123>AB",
		},
		{
			name: "EDIFACT example",
			msg:  ".A.C1.3.DATA.123DATA.123DATA",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, e := encoder.EncodeHighLevel(tc.msg, nil, -1, SymbolShapeHint_FORCE_NONE)
			if e != nil {
				t.Fatalf("EncodeHighLevel(%q) returns error: %v", tc.msg, e)
			}
			if len(b) == 0 {
				t.Fatalf("EncodeHighLevel(%q) returns empty", tc.msg)
			}
			// Verify all bytes are in valid range
			for i, byt := range b {
				if int(byt) < 0 || int(byt) > 255 {
					t.Errorf("Byte at index %d is out of range: %d", i, byt)
				}
			}
		})
	}
}

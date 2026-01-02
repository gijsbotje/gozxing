package encoder

import (
	"testing"
)

func TestGetErrorCorrectionCodewordCount(t *testing.T) {
	for level := 0; level <= 8; level++ {
		count, err := GetErrorCorrectionCodewordCount(level)
		if err != nil {
			t.Fatalf("GetErrorCorrectionCodewordCount(%d) returns error: %v", level, err)
		}
		expected := 1 << (level + 1)
		if count != expected {
			t.Fatalf("GetErrorCorrectionCodewordCount(%d) = %d, expect %d", level, count, expected)
		}
	}

	_, err := GetErrorCorrectionCodewordCount(-1)
	if err == nil {
		t.Fatalf("GetErrorCorrectionCodewordCount(-1) must be error")
	}

	_, err = GetErrorCorrectionCodewordCount(9)
	if err == nil {
		t.Fatalf("GetErrorCorrectionCodewordCount(9) must be error")
	}
}

func TestGetRecommendedMinimumErrorCorrectionLevel(t *testing.T) {
	testCases := []struct {
		n      int
		expect int
	}{
		{1, 2},
		{40, 2},
		{41, 3},
		{160, 3},
		{161, 4},
		{320, 4},
		{321, 5},
		{863, 5},
	}

	for _, tc := range testCases {
		level, err := GetRecommendedMinimumErrorCorrectionLevel(tc.n)
		if err != nil {
			t.Fatalf("GetRecommendedMinimumErrorCorrectionLevel(%d) returns error: %v", tc.n, err)
		}
		if level != tc.expect {
			t.Fatalf("GetRecommendedMinimumErrorCorrectionLevel(%d) = %d, expect %d", tc.n, level, tc.expect)
		}
	}

	_, err := GetRecommendedMinimumErrorCorrectionLevel(0)
	if err == nil {
		t.Fatalf("GetRecommendedMinimumErrorCorrectionLevel(0) must be error")
	}

	_, err = GetRecommendedMinimumErrorCorrectionLevel(-1)
	if err == nil {
		t.Fatalf("GetRecommendedMinimumErrorCorrectionLevel(-1) must be error")
	}

	_, err = GetRecommendedMinimumErrorCorrectionLevel(864)
	if err == nil {
		t.Fatalf("GetRecommendedMinimumErrorCorrectionLevel(864) must be error")
	}
}

func TestGenerateErrorCorrection(t *testing.T) {
	// Test with a simple data codeword string
	// Each rune represents a codeword value (0-928)
	dataCodewords := string([]rune{5, 10, 15})
	ec, err := GenerateErrorCorrection(dataCodewords, 0)
	if err != nil {
		t.Fatalf("GenerateErrorCorrection returns error: %v", err)
	}
	if len([]rune(ec)) == 0 {
		t.Fatalf("GenerateErrorCorrection returns empty string")
	}

	// Test with different error correction levels
	for level := 0; level <= 8; level++ {
		ec, err := GenerateErrorCorrection(dataCodewords, level)
		if err != nil {
			t.Fatalf("GenerateErrorCorrection(level %d) returns error: %v", level, err)
		}
		expectedLen, _ := GetErrorCorrectionCodewordCount(level)
		// Count runes, not bytes
		actualLen := len([]rune(ec))
		if actualLen != expectedLen {
			t.Fatalf("GenerateErrorCorrection(level %d) length = %d, expect %d", level, actualLen, expectedLen)
		}
	}
}


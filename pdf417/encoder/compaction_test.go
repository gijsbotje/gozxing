package encoder

import (
	"testing"
)

func TestCompaction_ValueOf(t *testing.T) {
	testCases := []struct {
		input    string
		expected Compaction
		shouldErr bool
	}{
		{"AUTO", Compaction_AUTO, false},
		{"TEXT", Compaction_TEXT, false},
		{"BYTE", Compaction_BYTE, false},
		{"NUMERIC", Compaction_NUMERIC, false},
		{"INVALID", Compaction_AUTO, true},
	}

	for _, tc := range testCases {
		result, err := Compaction_ValueOf(tc.input)
		if tc.shouldErr {
			if err == nil {
				t.Fatalf("Compaction_ValueOf(%s) should return error", tc.input)
			}
		} else {
			if err != nil {
				t.Fatalf("Compaction_ValueOf(%s) returns error: %v", tc.input, err)
			}
			if result != tc.expected {
				t.Fatalf("Compaction_ValueOf(%s) = %v, expect %v", tc.input, result, tc.expected)
			}
		}
	}
}

func TestCompaction_String(t *testing.T) {
	testCases := []struct {
		compaction Compaction
		expected   string
	}{
		{Compaction_AUTO, "AUTO"},
		{Compaction_TEXT, "TEXT"},
		{Compaction_BYTE, "BYTE"},
		{Compaction_NUMERIC, "NUMERIC"},
	}

	for _, tc := range testCases {
		result := tc.compaction.String()
		if result != tc.expected {
			t.Fatalf("Compaction(%v).String() = %s, expect %s", tc.compaction, result, tc.expected)
		}
	}
}


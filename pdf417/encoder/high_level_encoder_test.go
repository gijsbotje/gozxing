package encoder

import (
	"testing"

	"golang.org/x/text/encoding/charmap"
)

func TestEncodeHighLevel_Text(t *testing.T) {
	result, err := EncodeHighLevel("Hello", Compaction_TEXT, charmap.ISO8859_1, false)
	if err != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", err)
	}
	if len(result) == 0 {
		t.Fatalf("EncodeHighLevel returns empty result")
	}
}

func TestEncodeHighLevel_Byte(t *testing.T) {
	result, err := EncodeHighLevel("test", Compaction_BYTE, charmap.ISO8859_1, false)
	if err != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", err)
	}
	if len(result) == 0 {
		t.Fatalf("EncodeHighLevel returns empty result")
	}
}

func TestEncodeHighLevel_Numeric(t *testing.T) {
	result, err := EncodeHighLevel("1234567890", Compaction_NUMERIC, charmap.ISO8859_1, false)
	if err != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", err)
	}
	if len(result) == 0 {
		t.Fatalf("EncodeHighLevel returns empty result")
	}
}

func TestEncodeHighLevel_Auto(t *testing.T) {
	result, err := EncodeHighLevel("Hello123", Compaction_AUTO, charmap.ISO8859_1, false)
	if err != nil {
		t.Fatalf("EncodeHighLevel returns error: %v", err)
	}
	if len(result) == 0 {
		t.Fatalf("EncodeHighLevel returns empty result")
	}
}

func TestEncodeHighLevel_EmptyMessage(t *testing.T) {
	_, err := EncodeHighLevel("", Compaction_AUTO, charmap.ISO8859_1, false)
	if err == nil {
		t.Fatalf("EncodeHighLevel with empty message must return error")
	}
}


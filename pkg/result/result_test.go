package result

import (
	"math"
	"testing"

	"github.com/patterninc/heimdall/pkg/result/column"
)

func TestDecodeDecimal(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		scale    int
		expected float64
	}{
		{
			name:     "zero value",
			bytes:    make([]byte, 16),
			scale:    2,
			expected: 0.0,
		},
		{
			name: "positive decimal 123.45 with scale 2",
			bytes: []byte{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0x30, 0x39, // 12345
			},
			scale:    2,
			expected: 123.45,
		},
		{
			name: "negative decimal -123.45 with scale 2",
			bytes: func() []byte {
				// Two's complement of 12345 in 16 bytes
				b := []byte{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xCF, 0xC7,
				}
				return b
			}(),
			scale:    2,
			expected: -123.45,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decodeDecimal(tt.bytes, tt.scale)
			f, ok := result.(float64)
			if !ok {
				t.Fatalf("expected float64, got %T", result)
			}
			if math.Abs(f-tt.expected) > 0.001 {
				t.Errorf("got %f, want %f", f, tt.expected)
			}
		})
	}
}

func TestDecodeDecimalNil(t *testing.T) {
	result := decodeDecimal(nil, 2)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestDecodeDecimalNonBytes(t *testing.T) {
	result := decodeDecimal("not bytes", 2)
	if result != "not bytes" {
		t.Errorf("expected passthrough for non-byte types, got %v", result)
	}
}

func TestExtractUnionValue(t *testing.T) {
	t.Run("primitive type lookup", func(t *testing.T) {
		c := &column.Column{Name: "name", Type: "string"}
		m := map[string]any{"string": "hello"}
		val := extractUnionValue(m, c)
		if val != "hello" {
			t.Errorf("got %v, want %q", val, "hello")
		}
	})

	t.Run("named fixed type lookup", func(t *testing.T) {
		c := &column.Column{Name: "sp_spend", Type: "fixed", AvroTypeName: "topLevelRecord.sp_spend.fixed"}
		raw := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x30, 0x39}
		m := map[string]any{"topLevelRecord.sp_spend.fixed": raw}
		val := extractUnionValue(m, c)
		if b, ok := val.([]byte); !ok {
			t.Errorf("expected []byte, got %T", val)
		} else if len(b) != 16 {
			t.Errorf("expected 16 bytes, got %d", len(b))
		}
	})

	t.Run("fallback to first value", func(t *testing.T) {
		c := &column.Column{Name: "data", Type: "record"}
		m := map[string]any{"some.record.name": "value"}
		val := extractUnionValue(m, c)
		if val != "value" {
			t.Errorf("got %v, want %q", val, "value")
		}
	})
}

func TestDecodeDecimalScaleZero(t *testing.T) {
	// Decimal with scale 0 is effectively an integer but still needs byte decoding
	raw := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 42}
	result := decodeDecimal(raw, 0)
	f, ok := result.(float64)
	if !ok {
		t.Fatalf("expected float64, got %T", result)
	}
	if f != 42.0 {
		t.Errorf("got %f, want 42.0", f)
	}
}

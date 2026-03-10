package column

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Type
		wantErr  bool
	}{
		{
			name:     "simple string type",
			input:    `"string"`,
			expected: Type("string"),
		},
		{
			name:     "simple int type",
			input:    `"int"`,
			expected: Type("int"),
		},
		{
			name:     "nullable string array",
			input:    `["null", "string"]`,
			expected: Type("string"),
		},
		{
			name:     "nullable long array",
			input:    `["null", "long"]`,
			expected: Type("long"),
		},
		{
			name:     "nullable string array reversed order",
			input:    `["string", "null"]`,
			expected: Type("string"),
		},
		{
			name:     "complex object type - decimal",
			input:    `{"type": "bytes", "logicalType": "decimal", "precision": 38, "scale": 2}`,
			expected: Type("bytes"),
		},
		{
			name:     "complex object type - record/struct",
			input:    `{"type": "record", "name": "my_struct", "fields": [{"name": "a", "type": "string"}]}`,
			expected: Type("record"),
		},
		{
			name:     "complex object type - array",
			input:    `{"type": "array", "items": "string"}`,
			expected: Type("array"),
		},
		{
			name:     "complex object type - map",
			input:    `{"type": "map", "values": "string"}`,
			expected: Type("map"),
		},
		{
			name:     "nullable complex type - decimal in array",
			input:    `["null", {"type": "bytes", "logicalType": "decimal", "precision": 38, "scale": 2}]`,
			expected: Type("bytes"),
		},
		{
			name:     "nullable complex type - record in array",
			input:    `["null", {"type": "record", "name": "my_struct", "fields": []}]`,
			expected: Type("record"),
		},
		{
			name:    "array of only nulls",
			input:   `["null"]`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Type
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestIsPrimitive(t *testing.T) {
	primitives := []Type{"boolean", "int", "long", "float", "double", "bytes", "string"}
	for _, p := range primitives {
		if !p.IsPrimitive() {
			t.Errorf("expected %q to be primitive", p)
		}
	}

	nonPrimitives := []Type{"record", "array", "map", "enum", "fixed"}
	for _, np := range nonPrimitives {
		if np.IsPrimitive() {
			t.Errorf("expected %q to NOT be primitive", np)
		}
	}
}

func TestColumnUnmarshalJSON(t *testing.T) {
	input := `{
		"fields": [
			{"name": "id", "type": "long"},
			{"name": "name", "type": ["null", "string"]},
			{"name": "amount", "type": {"type": "bytes", "logicalType": "decimal", "precision": 38, "scale": 2}},
			{"name": "tags", "type": ["null", {"type": "array", "items": "string"}]}
		]
	}`

	type schema struct {
		Fields []Column `json:"fields"`
	}

	var s schema
	if err := json.Unmarshal([]byte(input), &s); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	expected := []struct {
		name      string
		typeName  Type
		scale     int
		isDecimal bool
	}{
		{"id", Type("long"), 0, false},
		{"name", Type("string"), 0, false},
		{"amount", Type("bytes"), 2, true},
		{"tags", Type("array"), 0, false},
	}

	if len(s.Fields) != len(expected) {
		t.Fatalf("got %d fields, want %d", len(s.Fields), len(expected))
	}

	for i, e := range expected {
		if s.Fields[i].Name != e.name {
			t.Errorf("field %d: got name %q, want %q", i, s.Fields[i].Name, e.name)
		}
		if s.Fields[i].Type != e.typeName {
			t.Errorf("field %d: got type %q, want %q", i, s.Fields[i].Type, e.typeName)
		}
		if s.Fields[i].Scale != e.scale {
			t.Errorf("field %d: got scale %d, want %d", i, s.Fields[i].Scale, e.scale)
		}
		if s.Fields[i].IsDecimal() != e.isDecimal {
			t.Errorf("field %d: got IsDecimal() %v, want %v", i, s.Fields[i].IsDecimal(), e.isDecimal)
		}
	}
}

func TestSensorAggBrandAdExpansionTacosSchema(t *testing.T) {
	avroSchema := `{
		"type": "record",
		"name": "sensor_agg_brand_ad_expansion_tacos",
		"fields": [
			{"name": "brand_id", "type": ["null", "long"]},
			{"name": "brand", "type": ["null", "string"]},
			{"name": "ad_market_id", "type": ["null", "long"]},
			{"name": "ad_market", "type": ["null", "string"]},
			{"name": "target_currency_code", "type": ["null", "string"]},
			{"name": "sp_spend", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.sp_spend.fixed"}]},
			{"name": "sb_spend", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.sb_spend.fixed"}]},
			{"name": "sbv_spend", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.sbv_spend.fixed"}]},
			{"name": "sd_spend", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.sd_spend.fixed"}]},
			{"name": "dsp_spend", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.dsp_spend.fixed"}]},
			{"name": "total_ad_spend", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.total_ad_spend.fixed"}]},
			{"name": "total_sales", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.total_sales.fixed"}]},
			{"name": "tacos", "type": ["null", {"type": "fixed", "size": 16, "logicalType": "decimal", "precision": 38, "scale": 2, "name": "topLevelRecord.tacos.fixed"}]},
			{"name": "metadata", "type": ["null", "string"]},
			{"name": "_etl_when_created", "type": ["null", "long"]},
			{"name": "_etl_when_updated", "type": ["null", "long"]}
		]
	}`

	type avroFields struct {
		Fields []Column `json:"fields"`
	}

	var schema avroFields
	if err := json.Unmarshal([]byte(avroSchema), &schema); err != nil {
		t.Fatalf("failed to unmarshal schema: %v", err)
	}

	if len(schema.Fields) != 16 {
		t.Fatalf("expected 16 fields, got %d", len(schema.Fields))
	}

	decimalFields := []string{"sp_spend", "sb_spend", "sbv_spend", "sd_spend", "dsp_spend", "total_ad_spend", "total_sales", "tacos"}
	for _, name := range decimalFields {
		found := false
		for _, f := range schema.Fields {
			if f.Name == name {
				found = true
				if f.Type != Type("fixed") {
					t.Errorf("field %q: expected type %q, got %q", name, "fixed", f.Type)
				}
				if !f.IsDecimal() {
					t.Errorf("field %q: expected IsDecimal() to be true", name)
				}
				if f.Scale != 2 {
					t.Errorf("field %q: expected scale 2, got %d", name, f.Scale)
				}
				expectedAvroName := "topLevelRecord." + name + ".fixed"
				if f.AvroTypeName != expectedAvroName {
					t.Errorf("field %q: expected AvroTypeName %q, got %q", name, expectedAvroName, f.AvroTypeName)
				}
				break
			}
		}
		if !found {
			t.Errorf("field %q not found in schema", name)
		}
	}

	// Verify decimal fields with scale 0 are still detected as decimal
	t.Run("decimal with scale 0 is detected", func(t *testing.T) {
		scaleZeroSchema := `{"fields": [{"name": "count", "type": {"type": "bytes", "logicalType": "decimal", "precision": 38, "scale": 0}}]}`
		var s2 avroFields
		if err := json.Unmarshal([]byte(scaleZeroSchema), &s2); err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}
		if !s2.Fields[0].IsDecimal() {
			t.Errorf("expected IsDecimal() to be true for decimal with scale 0")
		}
		if s2.Fields[0].Scale != 0 {
			t.Errorf("expected scale 0, got %d", s2.Fields[0].Scale)
		}
	})

	nonDecimalFields := map[string]Type{
		"brand_id":             Type("long"),
		"brand":                Type("string"),
		"ad_market_id":         Type("long"),
		"ad_market":            Type("string"),
		"target_currency_code": Type("string"),
		"metadata":             Type("string"),
		"_etl_when_created":    Type("long"),
		"_etl_when_updated":    Type("long"),
	}
	for name, expectedType := range nonDecimalFields {
		for _, f := range schema.Fields {
			if f.Name == name {
				if f.Type != expectedType {
					t.Errorf("field %q: expected type %q, got %q", name, expectedType, f.Type)
				}
				if f.IsDecimal() {
					t.Errorf("field %q: expected IsDecimal() to be false", name)
				}
				break
			}
		}
	}
}

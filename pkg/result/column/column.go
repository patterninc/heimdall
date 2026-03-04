package column

import (
	"encoding/json"
	"strings"
)

type Column struct {
	Name         string `yaml:"name,omitempty" json:"name,omitempty"`
	Type         Type   `yaml:"type,omitempty" json:"type,omitempty"`
	Scale        int    `json:"-" yaml:"-"`
	AvroTypeName string `json:"-" yaml:"-"`
}

func (c *Column) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name string          `json:"name"`
		Type json.RawMessage `json:"type"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	c.Name = raw.Name

	if err := json.Unmarshal(raw.Type, &c.Type); err != nil {
		return err
	}

	c.extractDecimalMeta(raw.Type)
	return nil
}

func (c *Column) extractDecimalMeta(typeData json.RawMessage) {
	trimmed := strings.TrimSpace(string(typeData))

	if strings.HasPrefix(trimmed, "[") {
		var elems []json.RawMessage
		if err := json.Unmarshal(typeData, &elems); err != nil {
			return
		}
		for _, elem := range elems {
			elemStr := strings.TrimSpace(string(elem))
			if elemStr == `"null"` {
				continue
			}
			if strings.HasPrefix(elemStr, "{") {
				c.parseComplexMeta(elem)
			}
		}
	} else if strings.HasPrefix(trimmed, "{") {
		c.parseComplexMeta(typeData)
	}
}

func (c *Column) parseComplexMeta(data json.RawMessage) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return
	}

	if raw, ok := obj["name"]; ok {
		var name string
		if json.Unmarshal(raw, &name) == nil {
			c.AvroTypeName = name
		}
	}

	if raw, ok := obj["logicalType"]; ok {
		var logicalType string
		if json.Unmarshal(raw, &logicalType) == nil && logicalType == "decimal" {
			if s, ok := obj["scale"]; ok {
				json.Unmarshal(s, &c.Scale)
			}
		}
	}
}

func (c *Column) IsDecimal() bool {
	return c.Scale > 0
}

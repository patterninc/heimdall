package column

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/babourine/x/pkg/set"
)

var (
	primitiveTypes = set.New([]string{
		`boolean`,
		`int`,
		`long`,
		`float`,
		`double`,
		`bytes`,
		`string`,
	})
)

type Type string

func (t *Type) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))

	if strings.HasPrefix(trimmed, `[`) {
		var temp []json.RawMessage
		if err := json.Unmarshal(data, &temp); err != nil {
			return err
		}

		for _, elem := range temp {
			elemStr := strings.TrimSpace(string(elem))
			if elemStr == `"null"` {
				continue
			}
			if strings.HasPrefix(elemStr, `"`) {
				var s string
				if err := json.Unmarshal(elem, &s); err != nil {
					return err
				}
				*t = Type(s)
				return nil
			}
			// Complex type (object) — extract its "type" field name or serialize as string
			*t = extractComplexTypeName(elem)
			return nil
		}

		return fmt.Errorf("unexpected value: %s", trimmed)

	} else if strings.HasPrefix(trimmed, `{`) {
		*t = extractComplexTypeName(data)
		return nil

	} else {
		var temp string
		if err := json.Unmarshal(data, &temp); err != nil {
			return err
		}
		*t = Type(temp)
	}

	return nil
}

func extractComplexTypeName(data json.RawMessage) Type {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		return Type(string(data))
	}
	if raw, ok := obj["type"]; ok {
		var typeName string
		if err := json.Unmarshal(raw, &typeName); err == nil {
			return Type(typeName)
		}
	}
	return Type(string(data))
}

func (t Type) IsPrimitive() bool {
	return primitiveTypes.Has(string(t))
}

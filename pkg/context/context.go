package context

import (
	"encoding/json"
	"strings"
)

const redactedValue = `REDACTED`

// sensitiveKeys are context field names whose values must never leave the
// process via JSON APIs. Matching is case-insensitive. Values are retained
// in-memory and when persisting via String().
var sensitiveKeys = map[string]struct{}{
	`password`:              {},
	`private_key`:           {},
	`secret`:                {},
	`token`:                 {},
	`api_key`:               {},
	`access_key`:            {},
	`secret_key`:            {},
	`client_secret`:         {},
	`secret_access_key`:     {},
	`aws_secret_access_key`: {},
}

type Context map[string]any

func New(v any) *Context {

	// Avoid Context.MarshalJSON so secrets are not redacted when cloning.
	data, err := marshalRaw(v)
	if err != nil {
		panic(`cannot marshal json`)
	}

	value := make(map[string]any)

	if err := json.Unmarshal(data, &value); err != nil {
		panic(`cannot unmarshal json`)
	}

	return (*Context)(&value)

}

func marshalRaw(v any) ([]byte, error) {
	switch t := v.(type) {
	case Context:
		return json.Marshal(map[string]any(t))
	case *Context:
		if t == nil {
			return []byte(`null`), nil
		}
		return json.Marshal(map[string]any(*t))
	default:
		return json.Marshal(v)
	}
}

func (c *Context) UnmarshalYAML(unmarshal func(any) error) error {

	value := make(map[string]any)

	if err := unmarshal(&value); err != nil {
		return err
	}

	*c = value

	return nil

}

func (c *Context) UnmarshalJSON(data []byte) error {

	value := make(map[string]any)

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	*c = value

	return nil

}

// MarshalJSON redacts sensitive fields for API responses. Internal helpers
// (String, Unmarshal) marshal the underlying map so plugins and DB storage
// still see real secrets.
func (c Context) MarshalJSON() ([]byte, error) {
	if c == nil {
		return []byte(`null`), nil
	}
	return json.Marshal(redactMap(c))
}

func (c *Context) Unmarshal(v any) error {

	// Marshal the underlying map so sensitive values are not redacted.
	data, err := json.Marshal(map[string]any(*c))

	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)

}

func (c *Context) String() string {

	if c == nil {
		return ``
	}

	// Persist the real context; do not go through MarshalJSON redaction.
	data, _ := json.Marshal(map[string]any(*c))

	return string(data)

}

func isSensitiveKey(key string) bool {
	_, ok := sensitiveKeys[strings.ToLower(key)]
	return ok
}

func redactMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		if isSensitiveKey(k) {
			out[k] = redactedValue
			continue
		}
		out[k] = redactValue(v)
	}
	return out
}

func redactValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return redactMap(t)
	case []any:
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = redactValue(item)
		}
		return out
	default:
		return v
	}
}

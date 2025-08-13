package duration

import (
	"encoding/json"
	"time"

	"gopkg.in/yaml.v3"
)

type Duration time.Duration

// UnmarshalYAML allows YAML like "duration: 15s"
func (d *Duration) UnmarshalYAML(value *yaml.Node) error {

	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)

	return nil

}

// UnmarshalJSON allows JSON like "duration": "15s"
func (d *Duration) UnmarshalJSON(data []byte) error {

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(dur)

	return nil

}

// MarshalJSON marshals the duration as a string
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

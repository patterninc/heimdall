package status

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Status byte

const (
	Active   Status = 1
	Inactive Status = 2
	Deleted  Status = 3
)

const (
	formatErrUnknownStatus = "unknown status: %v"
)

var (
	statusMapping = map[string]Status{
		`active`:   Active,
		`inactive`: Inactive,
		`deleted`:  Deleted,
	}
)

func (s *Status) UnmarshalYAML(unmarshal func(interface{}) error) error {

	value, err := UnmarshalYaML(unmarshal, statusMapping)

	if err != nil {
		return err
	}

	*s = value

	return nil

}

func (s *Status) UnmarshalJSON(data []byte) error {

	unmarshal := func(v interface{}) error {
		return json.Unmarshal(data, v)
	}

	value, err := UnmarshalYaML(unmarshal, statusMapping)
	if err != nil {
		return err
	}

	*s = value

	return nil

}

func (s *Status) MarshalJSON() ([]byte, error) {

	for k, v := range statusMapping {
		if v == *s {
			return json.Marshal(strings.ToUpper(k))
		}
	}

	return nil, fmt.Errorf(formatErrUnknownStatus, *s)

}

func UnmarshalYaML(unmarshal func(interface{}) error, mapping map[string]Status) (Status, error) {

	var temp string

	if err := unmarshal(&temp); err != nil {
		return 0, err
	}

	value, found := mapping[strings.ToLower(temp)]
	if !found {
		return 0, fmt.Errorf(formatErrUnknownStatus, temp)
	}

	return value, nil

}

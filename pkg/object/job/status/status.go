package status

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/patterninc/heimdall/pkg/object/status"
)

type Status status.Status

const (
	New       Status = 1
	Accepted  Status = 2
	Running   Status = 3
	Failed    Status = 4
	Killed    Status = 5
	Succeeded Status = 6
)

const (
	formatErrUnknownStatus = "unknown status: %v"
)

var (
	statusMapping = map[string]status.Status{
		`new`:       status.Status(New),
		`accepted`:  status.Status(Accepted),
		`running`:   status.Status(Running),
		`failed`:    status.Status(Failed),
		`killed`:    status.Status(Killed),
		`succeeded`: status.Status(Succeeded),
	}
)

func (s *Status) UnmarshalYAML(unmarshal func(interface{}) error) error {

	value, err := status.UnmarshalYaML(unmarshal, statusMapping)

	if err != nil {
		return err
	}

	*s = Status(value)

	return nil

}

func (s *Status) MarshalJSON() ([]byte, error) {

	for k, v := range statusMapping {
		if v == status.Status(*s) {
			return json.Marshal(strings.ToUpper(k))
		}
	}

	return nil, fmt.Errorf(formatErrUnknownStatus, *s)

}

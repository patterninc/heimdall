package result

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

type columnType string

func (t *columnType) UnmarshalJSON(data []byte) error {

	if strings.HasPrefix(string(data), `[`) {

		temp := make([]string, 0, 2)

		if err := json.Unmarshal(data, &temp); err != nil {
			return err
		}

		for _, v := range temp {
			if v != `null` {
				*t = columnType(temp[0])
				return nil
			}
		}

		return fmt.Errorf("unexpected value: %s", string(data))

	} else {

		temp := ``

		if err := json.Unmarshal(data, &temp); err != nil {
			return err
		}

		*t = columnType(temp)

	}

	return nil

}

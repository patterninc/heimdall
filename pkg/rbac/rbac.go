package rbac

import (
	"errors"
	"fmt"

	"github.com/patterninc/heimdall/pkg/rbac/ranger"
	"gopkg.in/yaml.v3"
)

var (
	ErrRBACIDsAreNotUnique = errors.New("rbac IDs are not unique")
	supportedRBACs         = map[string]func() RBAC{
		`apache_ranger`: NewRanger,
	}
)

type RBAC interface {
	Init() error
	HasAccess(user string, query string) (bool, error)
	GetName() string
}

type RBACs map[string]RBAC

type configs struct {
	RBAC []RBAC
}

func (c *RBACs) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var temp configs

	if err := unmarshal(&temp); err != nil {
		return err
	}

	items := make(map[string]RBAC)

	for _, t := range temp.RBAC {
		items[t.GetName()] = t
	}

	if len(temp.RBAC) != len(items) {
		return ErrRBACIDsAreNotUnique
	}

	*c = items

	return nil

}

// Implements custom unmarshaling based on `type` field in YAML
func (c *configs) UnmarshalYAML(value *yaml.Node) error {
	for _, value := range value.Content {
		var probe struct {
			Type string `yaml:"type"`
		}
		if err := value.Decode(&probe); err != nil {
			return err
		}

		supportedRBAC, ok := supportedRBACs[probe.Type]
		if !ok {
			return fmt.Errorf("unsupported RBAC type: %s", probe.Type)
		}
		r := supportedRBAC()
		if err := value.Decode(r); err != nil {
			return err
		}
		c.RBAC = append(c.RBAC, r)
	}
	return nil
}

func NewRanger() RBAC {
	return &ranger.Ranger{}
}

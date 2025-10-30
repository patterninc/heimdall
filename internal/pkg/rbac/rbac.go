package rbac

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/patterninc/heimdall/internal/pkg/rbac/ranger"
	"github.com/patterninc/heimdall/pkg/rbac"
)

var (
	ErrRBACIDsAreNotUnique = errors.New("rbac IDs are not unique")
	supportedRBACs         = map[string]func() rbac.RBAC{
		`apache_ranger`: ranger.New,
	}
)

type R rbac.RBAC

type RBACs map[string]R

type configs struct {
	RBAC []R
}

func (c *RBACs) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var temp configs

	if err := unmarshal(&temp); err != nil {
		return err
	}

	items := make(map[string]R)

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

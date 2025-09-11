package rbac

import (
	"context"
	"errors"
	"fmt"

	"github.com/patterninc/heimdall/pkg/rbac/ranger"
	parserFactory "github.com/patterninc/heimdall/pkg/sql/parser/factory"
	"gopkg.in/yaml.v3"
)

var (
	ErrRBACIDsAreNotUnique = errors.New("rbac IDs are not unique")
)

type RBAC interface {
	Init(ctx context.Context) error
	HasAccess(user string, query string) (bool, error)
	GetName() string
}

type RBACs map[string]RBAC

type RBACConfigs struct {
	RBAC []RBAC
}

func (c *RBACs) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var temp RBACConfigs

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
func (c *RBACConfigs) UnmarshalYAML(value *yaml.Node) error {
	for _, value := range value.Content {
		var probe struct {
			Type string `yaml:"type"`
		}
		if err := value.Decode(&probe); err != nil {
			return err
		}

		switch probe.Type {
		case "apache_ranger":
			var r ranger.ApacheRanger
			if err := value.Decode(&r); err != nil {
				return err
			}
			c.RBAC = append(c.RBAC, &r)
			parser, err := parserFactory.CreateParserByType(r.Parser.Type, r.Parser.DefaultCatalog)
			if err != nil {
				return err
			}
			r.AccessReceiver = parser
		default:
			return fmt.Errorf("unknown RBAC type: %s", probe.Type)
		}
	}
	return nil
}

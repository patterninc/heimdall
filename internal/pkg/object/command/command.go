package command

import (
	"fmt"

	"github.com/babourine/x/pkg/set"

	"github.com/patterninc/heimdall/internal/pkg/object"
	"github.com/patterninc/heimdall/internal/pkg/object/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

var (
	ErrCommandIDsAreNotUnique = fmt.Errorf(`command IDs are not unique`)
)

type Command struct {
	object.Object `yaml:",inline" json:",inline"`
	Status        status.Status    `yaml:"status,omitempty" json:"status,omitempty"`
	Plugin        string           `yaml:"plugin,omitempty" json:"plugin,omitempty"`
	IsSync        bool             `yaml:"is_sync,omitempty" json:"is_sync,omitempty"`
	ClusterTags   *set.Set[string] `yaml:"cluster_tags,omitempty" json:"cluster_tags,omitempty"`
	Handler       plugin.Handler   `yaml:"-" json:"-"`
}

type Commands map[string]*Command

func (c *Commands) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var temp []*Command

	if err := unmarshal(&temp); err != nil {
		return err
	}

	items := make(map[string]*Command)

	for _, t := range temp {
		if t.ID == `` {
			t.ID = t.Name
		}
		items[t.ID] = t
	}

	if len(temp) != len(items) {
		return ErrCommandIDsAreNotUnique
	}

	*c = items

	return nil

}

func (c *Command) Init() error {

	if err := c.Object.Init(); err != nil {
		return err
	}

	if c.Status == 0 {
		c.Status = status.Active
	}

	return nil

}

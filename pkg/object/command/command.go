package command

import (
	"fmt"
	"regexp"

	"github.com/babourine/x/pkg/set"

	"github.com/patterninc/heimdall/pkg/object"
	"github.com/patterninc/heimdall/pkg/object/status"
	"github.com/patterninc/heimdall/pkg/plugin"
)

var (
	ErrCommandIDsAreNotUnique = fmt.Errorf(`command IDs are not unique`)
)

type Command struct {
	object.Object  `yaml:",inline" json:",inline"`
	Status         status.Status    `yaml:"status,omitempty" json:"status,omitempty"`
	Plugin         string           `yaml:"plugin,omitempty" json:"plugin,omitempty"`
	IsSync         bool             `yaml:"is_sync,omitempty" json:"is_sync,omitempty"`
	ClusterTags    *set.Set[string] `yaml:"cluster_tags,omitempty" json:"cluster_tags,omitempty"`
	AllowedCallers *set.Set[string] `yaml:"allowed_callers,omitempty" json:"allowed_callers,omitempty"`
	Handler        plugin.Handler   `yaml:"-" json:"-"`
	callerPatterns []*regexp.Regexp `yaml:"-" json:"-"`
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

	if err := c.setCallerPatterns(); err != nil {
		return err
	}

	return nil

}

func (c *Command) IsCallerAllowed(user string) bool {
	if len(c.callerPatterns) == 0 {
		return true
	}
	for _, re := range c.callerPatterns {
		if re.MatchString(user) {
			return true
		}
	}
	return false
}

func (c *Command) setCallerPatterns() error {
	if c.AllowedCallers == nil || c.AllowedCallers.Len() == 0 {
		return nil
	}
	patterns := make([]*regexp.Regexp, 0, c.AllowedCallers.Len())
	for _, raw := range c.AllowedCallers.Slice() {
		re, err := regexp.Compile(`^(?:` + raw + `)$`)
		if err != nil {
			return fmt.Errorf("command %s: invalid allowed_callers pattern %q: %w", c.ID, raw, err)
		}
		patterns = append(patterns, re)
	}
	c.callerPatterns = patterns
	return nil
}

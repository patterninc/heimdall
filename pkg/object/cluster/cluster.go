package cluster

import (
	"fmt"

	"github.com/patterninc/heimdall/internal/pkg/object"
	"github.com/patterninc/heimdall/pkg/object/status"
)

var (
	ErrClusterIDsAreNotUnique = fmt.Errorf(`cluster IDs are not unique`)
)

type Cluster struct {
	object.Object `yaml:",inline" json:",inline"`
	Status        status.Status `yaml:"status,omitempty" json:"status,omitempty"`
}

type Clusters map[string]*Cluster

func (c *Clusters) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var temp []*Cluster

	if err := unmarshal(&temp); err != nil {
		return err
	}

	items := make(map[string]*Cluster)

	for _, t := range temp {
		if t.ID == `` {
			t.ID = t.Name
		}
		items[t.ID] = t
	}

	if len(temp) != len(items) {
		return ErrClusterIDsAreNotUnique
	}

	*c = items

	return nil

}

func (c *Cluster) Init() error {

	if err := c.Object.Init(); err != nil {
		return err
	}

	if c.Status == 0 {
		c.Status = status.Active
	}

	return nil

}

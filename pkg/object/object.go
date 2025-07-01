package object

import (
	"fmt"

	"github.com/babourine/x/pkg/set"

	"github.com/patterninc/heimdall/pkg/context"
)

const (
	defaultUser      = `heimdall`
	formatIDTag      = "_id:%s"
	formatNameTag    = "_name:%s"
	formatVersionTag = "_version:%s"
)

type Object struct {
	SystemID    int64            `yaml:"-" json:"-"`
	ID          string           `yaml:"id,omitempty" json:"id,omitempty"`
	Name        string           `yaml:"name,omitempty" json:"name,omitempty"`
	Version     string           `yaml:"version,omitempty" json:"version,omitempty"`
	User        string           `yaml:"user,omitempty" json:"user,omitempty"`
	Description string           `yaml:"description,omitempty" json:"description,omitempty"`
	Tags        *set.Set[string] `yaml:"tags,omitempty" json:"tags,omitempty"`
	Context     *context.Context `yaml:"context,omitempty" json:"context,omitempty"`
	CreatedAt   int              `yaml:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   int              `yaml:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (o *Object) Init() error {

	if o.ID == `` {
		o.ID = o.Name
	}

	if o.User == `` {
		o.User = defaultUser
	}

	// adding system tags
	if o.Tags == nil {
		o.Tags = set.New([]string{})
	}

	o.Tags.Add(fmt.Sprintf(formatIDTag, o.ID))
	o.Tags.Add(fmt.Sprintf(formatNameTag, o.Name))
	o.Tags.Add(fmt.Sprintf(formatVersionTag, o.Version))

	return nil

}

func (o *Object) GetID() string {
	return o.ID
}

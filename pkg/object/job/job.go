package job

import (
	"github.com/babourine/x/pkg/set"
	"github.com/google/uuid"

	"github.com/patterninc/heimdall/pkg/object"
	"github.com/patterninc/heimdall/pkg/object/job/status"
	"github.com/patterninc/heimdall/pkg/result"
)

type Job struct {
	object.Object      `yaml:",inline" json:",inline"`
	Status             status.Status        `yaml:"status,omitempty" json:"status,omitempty"`
	IsSync             bool                 `yaml:"is_sync,omitempty" json:"is_sync,omitempty"`
	StoreResultSync    bool                 `yaml:"store_result_sync,omitempty" json:"store_result_sync,omitempty"`
	Error              string               `yaml:"error,omitempty" json:"error,omitempty"`
	CommandCriteria    *set.Set[string]     `yaml:"command_criteria,omitempty" json:"command_criteria,omitempty"`
	ClusterCriteria    *set.Set[string]     `yaml:"cluster_criteria,omitempty" json:"cluster_criteria,omitempty"`
	CommandID          string               `yaml:"command_id,omitempty" json:"command_id,omitempty"`
	CommandName        string               `yaml:"command_name,omitempty" json:"command_name,omitempty"`
	ClusterID          string               `yaml:"cluster_id,omitempty" json:"cluster_id,omitempty"`
	ClusterName        string               `yaml:"cluster_name,omitempty" json:"cluster_name,omitempty"`
	CanceledBy         string               `yaml:"canceled_by,omitempty" json:"canceled_by,omitempty"`
	ExtraJobAttributes map[string]Attribute `yaml:"extra_job_attributes,omitempty" json:"extra_job_attributes,omitempty"`
	Result             *result.Result       `yaml:"result,omitempty" json:"result,omitempty"`

	outputs map[string]string
}

// Attribute kinds for ExtraJobAttributes. Kind tells the UI how to render Value.
const (
	AttributeKindLink = "link"
	AttributeKindText = "text"
)

// Attribute is a generic, plugin-agnostic value surfaced on a job in the UI.
type Attribute struct {
	Kind  string `yaml:"kind,omitempty" json:"kind,omitempty"`
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
}

func (j *Job) Init() error {

	// we override job ID so clients dont submit theirs -- server must issue job ID
	j.ID = uuid.NewString()

	if err := j.Object.Init(); err != nil {
		return err
	}

	if j.Status == 0 {
		j.Status = status.New
	}

	return nil

}

func (j *Job) SetOutput(key, value string) {
	if j.outputs == nil {
		j.outputs = make(map[string]string)
	}
	j.outputs[key] = value
}

func (j *Job) Outputs() map[string]string {
	return j.outputs
}

package column

type Column struct {
	Name string `yaml:"name,omitempty" json:"name,omitempty"`
	Type Type   `yaml:"type,omitempty" json:"type,omitempty"`
}

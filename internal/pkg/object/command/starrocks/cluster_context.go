package starrocks

type clusterContext struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Database string `yaml:"database,omitempty" json:"database,omitempty"`
	UseTLS   bool   `yaml:"use_tls,omitempty" json:"use_tls,omitempty"`
}

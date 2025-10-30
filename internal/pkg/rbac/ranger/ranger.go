package ranger

import (
	"context"
	"errors"
	"log"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser"
	"github.com/patterninc/heimdall/internal/pkg/sql/parser/factory"
	"github.com/patterninc/heimdall/pkg/rbac"
)

var (
	ErrRangerClientConfigIsRequired         = errors.New("ranger client_config is required")
	ErrRangerParserConfigIsRequired         = errors.New("ranger parser_config is required")
	ErrRangerParserTypeIsRequired           = errors.New("ranger parser_config.type is required")
	ErrRangerParserDefaultCatalogIsRequired = errors.New("ranger parser_config.default_catalog is required")
	ErrRangerUnsupportedParserType          = errors.New("unsupported ranger parser_config.type. supported types: trino")
)

type Ranger struct {
	Name                  string `yaml:"name,omitempty" json:"name,omitempty"`
	ServiceName           string `yaml:"service_name,omitempty" json:"service_name,omitempty"`
	Client                Client
	SyncIntervalInMinutes int                   `yaml:"sync_interval_in_minutes,omitempty" json:"sync_interval_in_minutes,omitempty"`
	AccessReceiver        parser.AccessReceiver `yaml:"parser,omitempty" json:"parser,omitempty"`
	permissionsByUser     map[string]*userPermissions
}

type parserConfig struct {
	Type           string `yaml:"type,omitempty" json:"type,omitempty"`
	DefaultCatalog string `yaml:"default_catalog,omitempty" json:"default_catalog,omitempty"`
}

type clientConfig struct {
	Endpoint string `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

type userPermissions struct {
	AllowPolicies map[parser.Action][]*Policy
	DenyPolicies  map[parser.Action][]*Policy
}

func (r *Ranger) Init() error {
	// first time lets sync state explicitly
	if err := r.SyncState(); err != nil {
		return err
	}
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		ticker := time.NewTicker(time.Duration(r.SyncIntervalInMinutes) * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Apache Ranger sync goroutine")
				return
			case <-ticker.C:
				log.Println("Syncing policies from Apache Ranger for service:", r.ServiceName)
				if err := r.SyncState(); err != nil {
					log.Println("Error syncing users and groups from Apache Ranger", "error", err)
				}
			}
		}
	}()
	return nil
}

func (r *Ranger) HasAccess(user string, query string) (bool, error) {

	return true, nil
}

func (r *Ranger) GetName() string {
	return r.Name
}

func (r *Ranger) SyncState() error {

	return nil
}

func (r *Ranger) UnmarshalYAML(value *yaml.Node) error {
	type rawRanger struct {
		Name                  string        `yaml:"name,omitempty" json:"name,omitempty"`
		ServiceName           string        `yaml:"service_name,omitempty" json:"service_name,omitempty"`
		SyncIntervalInMinutes int           `yaml:"sync_interval_in_minutes,omitempty" json:"sync_interval_in_minutes,omitempty"`
		Client                *clientConfig `yaml:"client"`
		Parser                *parserConfig `yaml:"parser"`
	}

	var raw rawRanger
	if err := value.Decode(&raw); err != nil {
		return err
	}

	if raw.Client == nil {
		return ErrRangerClientConfigIsRequired
	}
	if raw.Parser == nil {
		return ErrRangerParserConfigIsRequired
	}
	if raw.Parser.Type == "" {
		return ErrRangerParserTypeIsRequired
	}
	if raw.Parser.DefaultCatalog == "" {
		return ErrRangerParserDefaultCatalogIsRequired
	}

	r.Name = raw.Name
	r.ServiceName = raw.ServiceName
	r.SyncIntervalInMinutes = raw.SyncIntervalInMinutes
	r.Client = NewClient(raw.Client.Endpoint, raw.Client.Username, raw.Client.Password)

	accessReceiver, err := factory.CreateParserByType(raw.Parser.Type, raw.Parser.DefaultCatalog)
	if err != nil {
		return ErrRangerUnsupportedParserType
	}

	r.AccessReceiver = accessReceiver
	if r.SyncIntervalInMinutes == 0 {
		r.SyncIntervalInMinutes = 5
	}
	return nil
}

func New() rbac.RBAC {
	return &Ranger{}
}

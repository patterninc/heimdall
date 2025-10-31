package ranger

import (
	"context"
	"errors"
	"log"
	"strings"
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
				if err := r.SyncState(); err != nil {
					log.Println("Error syncing users and groups from Apache Ranger", "error", err)
				}
			}
		}
	}()
	return nil
}

func (r *Ranger) HasAccess(user string, query string) (bool, error) {
	user = strings.ToLower(user)
	if _, ok := r.permissionsByUser[user]; !ok {
		log.Println("User not found in ranger policies. User: ", user)
		return false, nil
	}
	accessList, err := r.AccessReceiver.ParseAccess(query)
	if err != nil {
		return false, err
	}

	permissions := r.permissionsByUser[user]

	for _, access := range accessList {
		for _, permition := range permissions.DenyPolicies[access.Action()] {
			if permition.doesControlAnAccess(access) {
				log.Println("Access denied by ranger policy", "user", user, "policy", permition.Name, "action", access.Action(), "resource", access.QualifiedName())
				return false, nil
			}
		}
		foundAllowPolicy := false
		for _, permition := range permissions.AllowPolicies[access.Action()] {
			if permition.doesControlAnAccess(access) {
				log.Println("Access allowed by ranger policy", "user", user, "policy", permition.Name, "action", access.Action(), "resource", access.QualifiedName())
				foundAllowPolicy = true
				break
			}
		}
		if !foundAllowPolicy {
			log.Println("Access denied by ranger policy", "user", user, "action", access.Action(), "resource", access.QualifiedName())
			return false, nil
		}
	}
	return true, nil
}

func (r *Ranger) GetName() string {
	return r.Name
}

func (r *Ranger) SyncState() error {
	policies, err := r.Client.GetPolicies(r.ServiceName)
	if err != nil {
		return err
	}
	users, err := r.Client.GetUsers()
	if err != nil {
		return err
	}

	usersByGroup := map[string][]string{}
	for _, user := range users {
		for _, gName := range user.GroupNameList {
			usersByGroup[gName] = append(usersByGroup[gName], user.Name)
		}
	}

	newPermissionsByUser := map[string]*userPermissions{}
	for _, policy := range policies {
		if !policy.IsEnabled {
			continue
		}
		if policy.Resources == nil || policy.Resources.Catalog == nil || policy.Resources.Schema == nil || policy.Resources.Table == nil {
			// Skip policies that do not have catalog, schema, or table defined
			continue
		}

		if err := policy.init(); err != nil {
			log.Println("Error initializing policy:", err)
			return err
		}

		controlledActions := policy.getControlledActions(usersByGroup)
		for userName, actions := range controlledActions.allowedActionsByUser {
			if _, ok := newPermissionsByUser[userName]; !ok {
				newPermissionsByUser[userName] = &userPermissions{
					AllowPolicies: map[parser.Action][]*Policy{},
					DenyPolicies:  map[parser.Action][]*Policy{},
				}
			}
			for _, action := range actions {
				newPermissionsByUser[userName].AllowPolicies[action] = append(newPermissionsByUser[userName].AllowPolicies[action], policy)
			}
		}
		for userName, actions := range controlledActions.deniedActionsByUser {
			if _, ok := newPermissionsByUser[userName]; !ok {
				newPermissionsByUser[userName] = &userPermissions{
					AllowPolicies: map[parser.Action][]*Policy{},
					DenyPolicies:  map[parser.Action][]*Policy{},
				}
			}
			for _, action := range actions {
				newPermissionsByUser[userName].DenyPolicies[action] = append(newPermissionsByUser[userName].DenyPolicies[action], policy)
			}
		}
	}

	r.permissionsByUser = newPermissionsByUser
	log.Println("Syncing users and groups from Apache Ranger for service:", r.ServiceName)
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

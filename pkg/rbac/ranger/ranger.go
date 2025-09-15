package ranger

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/patterninc/heimdall/pkg/sql/parser"
)

type ApacheRanger struct {
	Name                  string        `yaml:"name,omitempty" json:"name,omitempty"`
	ServiceName           string        `yaml:"service_name,omitempty" json:"service_name,omitempty"`
	Client                ClientWrapper `yaml:"client,omitempty" json:"client,omitempty"`
	SyncIntervalInMinutes int           `yaml:"sync_interval_in_minutes,omitempty" json:"sync_interval_in_minutes,omitempty"`
	AccessReceiver        parser.AccessReceiver
	permissionsByUser      map[string]*UserPermissions
	Parser                ParserConfig `yaml:"parser,omitempty" json:"parser,omitempty"`
}

type ParserConfig struct {
	Type           string `yaml:"type,omitempty" json:"type,omitempty"`
	DefaultCatalog string `yaml:"default_catalog,omitempty" json:"default_catalog,omitempty"`
}

type PermissionStatus int

const (
	PermissionStatusAllow PermissionStatus = iota
	PermissionStatusDeny
	PermissionStatusUnknown
)

type UserPermissions struct {
	AllowPolicys map[parser.Action][]*Policy
	DenyPolicys  map[parser.Action][]*Policy
}

func (ar *ApacheRanger) Init(ctx context.Context) error {
	// first time lets sync state explicitly
	if err := ar.SyncState(); err != nil {
		return err
	}
	ar.startSyncPolicies(ctx)
	return nil
}

func (ar *ApacheRanger) HasAccess(user string, query string) (bool, error) {
	user = strings.ToLower(user)
	if _, ok := ar.permissionsByUser[user]; !ok {
		log.Println("User not found in ranger policies. User: ", user)
		return false, nil
	}
	accessList, err := ar.AccessReceiver.ParseAccess(query)
	if err != nil {
		return false, err
	}

	permissions := ar.permissionsByUser[user]

	for _, access := range accessList {
		for _, permition := range permissions.DenyPolicys[access.Action()] {
			if permition.doesControlAnAccess(access) {
				log.Println("Access denied by ranger policy", "user", user, "query", query, "policy", permition.Name, "action", access.Action(), "resource", access.QualifiedName())
				return false, nil
			}
		}
		foundAllowPolicy := false
		for _, permition := range permissions.AllowPolicys[access.Action()] {
			if permition.doesControlAnAccess(access) {
				log.Println("Access allowed by ranger policy", "user", user, "query", query, "policy", permition.Name, "action", access.Action(), "resource", access.QualifiedName())
				foundAllowPolicy = true
				break
			}
		}
		if !foundAllowPolicy {
			log.Println("Access denied by ranger policy", "user", user, "query", query, "action", access.Action(), "resource", access.QualifiedName())
			return false, nil
		}
	}
	return true, nil
}

func (r *ApacheRanger) GetName() string {
	return r.Name
}

func (r *ApacheRanger) SyncState() error {
	policies, err := r.Client.GetPolicies(r.ServiceName)
	if err != nil {
		return err
	}
	users, err := r.Client.GetUsers()
	if err != nil {
		return err
	}
	groups, err := r.Client.GetGroups()
	if err != nil {
		return err
	}

	log.Println("Users:", len(users), "Groups:", len(groups), "Policies:", len(policies))

	groupByID := map[int64]*Group{}
	usersByGroup := map[string][]string{}
	for _, group := range groups {
		groupByID[group.ID] = group
	}

	for _, user := range users {
		for _, gid := range user.GroupIdList {
			if group, ok := groupByID[gid]; ok {
				usersByGroup[group.Name] = append(usersByGroup[group.Name], user.Name)
			}
		}
	}

	newPermissionsByUser := map[string]*UserPermissions{}
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
				newPermissionsByUser[userName] = &UserPermissions{
					AllowPolicys: map[parser.Action][]*Policy{},
					DenyPolicys:  map[parser.Action][]*Policy{},
				}
			}
			for _, action := range actions {
				newPermissionsByUser[userName].AllowPolicys[action] = append(newPermissionsByUser[userName].AllowPolicys[action], policy)
			}
		}
		for userName, actions := range controlledActions.deniedActionsByUser {
			if _, ok := newPermissionsByUser[userName]; !ok {
				newPermissionsByUser[userName] = &UserPermissions{
					AllowPolicys: map[parser.Action][]*Policy{},
					DenyPolicys:  map[parser.Action][]*Policy{},
				}
			}
			for _, action := range actions {
				newPermissionsByUser[userName].DenyPolicys[action] = append(newPermissionsByUser[userName].DenyPolicys[action], policy)
			}
		}
	}

	r.permissionsByUser = newPermissionsByUser
	log.Println("Syncing users and groups from Apache Ranger for service:", r.ServiceName)
	return nil
}

func (ar *ApacheRanger) startSyncPolicies(ctx context.Context) {
	go func() {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		ticker := time.NewTicker(time.Duration(ar.SyncIntervalInMinutes) * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping Apache Ranger sync goroutine")
				return
			case <-ticker.C:
				log.Println("Syncing policies from Apache Ranger for service:", ar.ServiceName)
				if err := ar.SyncState(); err != nil {
					log.Println("Error syncing users and groups from Apache Ranger", "error", err)
				}
			}
		}
	}()
}

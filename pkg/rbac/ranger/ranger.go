package ranger

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/patterninc/heimdall/pkg/sql/parser"
)

type ApacheRanger struct {
	Name                  string `yaml:"name,omitempty" json:"name,omitempty"`
	ServiceName           string `yaml:"service_name,omitempty" json:"service_name,omitempty"`
	Client                Client
	SyncIntervalInMinutes int `yaml:"sync_interval_in_minutes,omitempty" json:"sync_interval_in_minutes,omitempty"`
	AccessReceiver        parser.AccessReceiver
	permitionsByUser      map[string]*UserPermitions
}

type PermitionStatus int

const (
	PermitionStatusAllow PermitionStatus = iota
	PermitionStatusDeny
	PermitionStatusUnknown
)

type UserPermitions struct {
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
	if _, ok := ar.permitionsByUser[user]; !ok {
		log.Println("User not found in ranger policies", "user", user)
		return false, nil
	}
	accessList, err := ar.AccessReceiver.ParseAccess(query)
	if err != nil {
		return false, err
	}

	permitions := ar.permitionsByUser[user]
	for _, access := range accessList {
		for _, permition := range permitions.DenyPolicys[access.Action()] {
			if permition.controlAnAccess(access) {
				log.Println("Access denied by ranger policy", "user", user, "query", query, "policy", permition.Name, "action", access.Action(), "resource", access.QualifiedName())
				return false, nil
			}
		}
		for _, permition := range permitions.AllowPolicys[access.Action()] {
			if permition.controlAnAccess(access) {
				log.Println("Access allowed by ranger policy", "user", user, "query", query, "policy", permition.Name, "action", access.Action(), "resource", access.QualifiedName())
				return true, nil
			}
		}
	}
	return false, nil
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

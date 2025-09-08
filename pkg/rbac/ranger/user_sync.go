package ranger

import (
	"log"

	"github.com/patterninc/heimdall/pkg/sql/parser"
)

func (r *ApacheRanger) SyncState() error {
	users, err := r.Client.GetUsers()
	if err != nil {
		return err
	}
	groups, err := r.Client.GetGroups()
	if err != nil {
		return err
	}
	policies, err := r.Client.GetPolicies(r.ServiceName)
	if err != nil {
		return err
	}
	println("Users:", len(users), "Groups:", len(groups), "Policies:", len(policies))

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
	
	newPermitionsByUser := map[string]*UserPermitions{}
	for _, policy := range policies {
		if !policy.IsEnabled {
			continue
		}
		if len(policy.Resources.Catalog.Values) == 0 || len(policy.Resources.Schema.Values) == 0 || len(policy.Resources.Table.Values) == 0 {
			// Skip policies that do not have catalog, schema, and table defined
			continue
		}

		if err := policy.init(); err != nil {
			log.Println("Error initializing policy:", err)
			return err
		}
		controlledActions := policy.getControlledActions(usersByGroup)
		for userName, actions := range controlledActions.allowedActionsByUser {
			if _, ok := newPermitionsByUser[userName]; !ok {
				newPermitionsByUser[userName] = &UserPermitions{
					AllowPolicys: map[parser.Action][]*Policy{},
					DenyPolicys:  map[parser.Action][]*Policy{},
				}
			}
			for _, action := range actions {
				newPermitionsByUser[userName].AllowPolicys[action] = append(newPermitionsByUser[userName].AllowPolicys[action], policy)
			}
		}
		for userName, actions := range controlledActions.deniedActionsByUser {
			if _, ok := newPermitionsByUser[userName]; !ok {
				newPermitionsByUser[userName] = &UserPermitions{
					AllowPolicys: map[parser.Action][]*Policy{},
					DenyPolicys:  map[parser.Action][]*Policy{},
				}
			}
			for _, action := range actions {
				newPermitionsByUser[userName].DenyPolicys[action] = append(newPermitionsByUser[userName].DenyPolicys[action], policy)
			}
		}
	}

	r.permitionsByUser = newPermitionsByUser
	log.Println("Syncing users and groups from Apache Ranger for service:", r.ServiceName)
	return nil
}

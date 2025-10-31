package ranger

import (
	"fmt"
	"log"
	"strings"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser"
)

const (
	allActionAccessType = "all"
)

var (
	actionByName = map[string]parser.Action{
		"select":                   parser.SELECT,
		"insert":                   parser.INSERT,
		"update":                   parser.UPDATE,
		"delete":                   parser.DELETE,
		"create":                   parser.CREATE,
		"drop":                     parser.DROP,
		"use":                      parser.USE,
		"alter":                    parser.ALTER,
		"grant":                    parser.GRANT,
		"revoke":                   parser.REVOKE,
		"show":                     parser.SHOW,
		"impersonate":              parser.IMPERSONATE,
		"execute":                  parser.EXECUTE,
		"read_system_information":  parser.READ_SYSTEM_INFORMATION,
		"write_system_information": parser.WRITE_SYSTEM_INFORMATION,
	}

	allActions = []parser.Action{
		parser.SELECT,
		parser.INSERT,
		parser.UPDATE,
		parser.DELETE,
		parser.CREATE,
		parser.DROP,
		parser.USE,
		parser.ALTER,
		parser.GRANT,
		parser.REVOKE,
		parser.SHOW,
		parser.IMPERSONATE,
		parser.EXECUTE,
		parser.READ_SYSTEM_INFORMATION,
		parser.WRITE_SYSTEM_INFORMATION,
	}
)

type Policy struct {
	ID                  int         `json:"id"`
	GUID                string      `json:"guid"`
	IsEnabled           bool        `json:"isEnabled"`
	Version             int         `json:"version"`
	Service             string      `json:"service"`
	Name                string      `json:"name"`
	PolicyType          int         `json:"policyType"`
	PolicyPriority      int         `json:"policyPriority"`
	Description         string      `json:"description"`
	IsAuditEnabled      bool        `json:"isAuditEnabled"`
	Resources           *Resource   `json:"resources"`
	AdditionalResources []*Resource `json:"additionalResources"`
	AllResources        []*Resource
	PolicyItems         []*PolicyItem `json:"policyItems"`
	DenyPolicyItems     []*PolicyItem `json:"denyPolicyItems"`
	AllowExceptions     []*PolicyItem `json:"allowExceptions"`
	DenyExceptions      []*PolicyItem `json:"denyExceptions"`
	ServiceType         string       `json:"serviceType"`
}

type Resource struct {
	Schema  *ResourceField `json:"schema,omitempty"`
	Catalog *ResourceField `json:"catalog,omitempty"`
	Table   *ResourceField `json:"table,omitempty"`
	Column  *ResourceField `json:"column,omitempty"`
}

type ResourceField struct {
	RawValues  []string `json:"values"`
	IsExcludes bool     `json:"isExcludes"`
	values     []*value
}

type value struct {
	value    string
	isRegexp bool
}

func (rf *ResourceField) IsPrefixFor(s string) bool {
	for _, v := range rf.values {
		if v.isRegexp && strings.HasPrefix(s, v.value) || s == v.value {
			return true
		}
	}
	return false
}

type Access struct {
	Type string `json:"type"`
}

type PolicyItem struct {
	Accesses []*Access `json:"accesses"`
	Users    []string `json:"users,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Actions  []parser.Action
}

type ControlledActions struct {
	allowedActionsByUser map[string][]parser.Action
	deniedActionsByUser  map[string][]parser.Action
}

func (p *Policy) init() (err error) {
	p.AllResources = append([]*Resource{p.Resources}, p.AdditionalResources...)
	for _, v := range p.AllResources {
		v.Catalog.values, err = preprocessValues(v.Catalog.RawValues)
		if err != nil {
			return err
		}
		v.Schema.values, err = preprocessValues(v.Schema.RawValues)
		if err != nil {
			return err
		}
		v.Table.values, err = preprocessValues(v.Table.RawValues)
		if err != nil {
			return err
		}
	}
	return nil
}

func preprocessValues(rawValues []string) ([]*value, error) {
	result := make([]*value, len(rawValues))
	for i, v := range rawValues {
		isRegexp := strings.HasSuffix(v, "*")
		v = strings.TrimRight(v, "*")
		if strings.Count(v, "*") > 0 {
			return nil, fmt.Errorf("Invalid value value %s* is allowed only at the end of the string. ")
		}
		result[i] = &value{
			isRegexp: isRegexp,
			value:    v,
		}
	}
	return result, nil
}

func (p *Policy) doesControlAnAccess(access parser.Access) bool {
	switch a := access.(type) {
	case *parser.TableAccess:
		return p.doesControlTableAccess(a)
	}
	return false
}

func (p *Policy) doesControlTableAccess(a *parser.TableAccess) bool {
	for _, v := range p.AllResources {

		if v.Catalog.IsPrefixFor(a.Catalog) == v.Catalog.IsExcludes {
			continue
		}

		if v.Schema.IsPrefixFor(a.Schema) == v.Schema.IsExcludes {
			continue
		}
		if v.Table.IsPrefixFor(a.Table) == v.Table.IsExcludes {
			continue
		}

		return true
	}
	return false
}

func (p *Policy) getControlledActions(usersByGroup map[string][]string) ControlledActions {
	return ControlledActions{
		allowedActionsByUser: p.getAllPolicyByUser(p.PolicyItems, p.AllowExceptions, usersByGroup),
		deniedActionsByUser:  p.getAllPolicyByUser(p.DenyPolicyItems, p.DenyExceptions, usersByGroup),
	}
}

func (p *Policy) getAllPolicyByUser(
	items []*PolicyItem,
	exceptions []*PolicyItem,
	usersByGroup map[string][]string,
) map[string][]parser.Action {
	policiesItem := policyItemsToActionsByUser(items, usersByGroup)
	exceptionsItem := policyItemsToActionsByUser(exceptions, usersByGroup)

	for user, actions := range exceptionsItem {
		if _, ok := policiesItem[user]; !ok {
			continue
		}
		for action := range actions {
			delete(policiesItem[user], action)
		}
		if len(policiesItem[user]) == 0 {
			delete(policiesItem, user)
		}
	}

	result := map[string][]parser.Action{}
	for user, actionsMap := range policiesItem {
		actions := make([]parser.Action, 0, len(actionsMap))
		for action := range actionsMap {
			actions = append(actions, action)
		}
		result[user] = actions
	}
	return result
}

func (p *PolicyItem) getPermissions() []parser.Action {
	if p.Actions != nil {
		return p.Actions
	}
	p.Actions = make([]parser.Action, 0)
	for _, access := range p.Accesses {
		accessType := strings.ToLower(access.Type)

		if accessType == allActionAccessType {
			return allActions
		}
		if action, ok := actionByName[accessType]; ok {
			p.Actions = append(p.Actions, action)
			continue
		} else {
			log.Println("Unknown action type in ranger policy:", accessType)
		}
	}
	return p.Actions

}

func policyItemsToActionsByUser(items []*PolicyItem, usersByGroup map[string][]string) map[string]map[parser.Action]struct{} {
	permissions := make(map[string]map[parser.Action]struct{})

	for _, item := range items {
		actions := item.getPermissions()
		for _, user := range item.Users {
			addActionsToPermissions(permissions, user, actions)
		}
		for _, group := range item.Groups {
			for _, user := range usersByGroup[group] {
				addActionsToPermissions(permissions, user, actions)
			}
		}
	}

	return permissions
}

func addActionsToPermissions(permissions map[string]map[parser.Action]struct{}, user string, actions []parser.Action) {
	user = strings.ToLower(user)
	if _, ok := permissions[user]; !ok {
		permissions[user] = make(map[parser.Action]struct{})
	}
	for _, action := range actions {
		permissions[user][action] = struct{}{}
	}
}

package ranger

import (
	"regexp"
	"strings"

	"github.com/patterninc/heimdall/pkg/sql/parser"
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
	ID                  int          `json:"id"`
	GUID                string       `json:"guid"`
	IsEnabled           bool         `json:"isEnabled"`
	Version             int          `json:"version"`
	Service             string       `json:"service"`
	Name                string       `json:"name"`
	PolicyType          int          `json:"policyType"`
	PolicyPriority      int          `json:"policyPriority"`
	Description         string       `json:"description"`
	IsAuditEnabled      bool         `json:"isAuditEnabled"`
	Resources           *Resource    `json:"resources"`
	AdditionalResources []*Resource  `json:"additionalResources"`
	PolicyItems         []PolicyItem `json:"policyItems"`
	DenyPolicyItems     []PolicyItem `json:"denyPolicyItems"`
	AllowExceptions     []PolicyItem `json:"allowExceptions"`
	DenyExceptions      []PolicyItem `json:"denyExceptions"`
	ServiceType         string       `json:"serviceType"`
}

type ResourceField struct {
	Values     []string `json:"values"`
	IsExcludes bool     `json:"isExcludes"`
	regexp     *regexp.Regexp
}

type Resource struct {
	Schema  *ResourceField `json:"schema,omitempty"`
	Catalog *ResourceField `json:"catalog,omitempty"`
	Table   *ResourceField `json:"table,omitempty"`
	Column  *ResourceField `json:"column,omitempty"`
}

type Access struct {
	Type string `json:"type"`
}

type PolicyItem struct {
	Accesses []Access `json:"accesses"`
	Users    []string `json:"users,omitempty"`
	Groups   []string `json:"groups,omitempty"`
	Actions  []parser.Action
}

type ControlledActions struct {
	allowedActionsByUser map[string][]parser.Action
	deniedActionsByUser  map[string][]parser.Action
}

func (p *Policy) init() error {
	for _, v := range append([]*Resource{p.Resources}, p.AdditionalResources...) {
		if len(v.Catalog.Values) != 0 {
			v.Catalog.regexp = regexp.MustCompile("^(" + patternsToRegex(v.Catalog.Values) + ")$")
		}
		if len(v.Schema.Values) != 0 {
			v.Schema.regexp = regexp.MustCompile("^(" + patternsToRegex(v.Schema.Values) + ")$")
		}
		if len(v.Table.Values) != 0 {
			v.Table.regexp = regexp.MustCompile("^(" + patternsToRegex(v.Table.Values) + ")$")
		}
	}
	return nil
}

func (p *Policy) controlAnAccess(access parser.Access) bool {
	switch a := access.(type) {
	case *parser.TableAccess:
		return p.controlTableAccess(a)
	}
	return false
}

func (p *Policy) controlTableAccess(a *parser.TableAccess) bool {
	for _, v := range append([]*Resource{p.Resources}, p.AdditionalResources...) {
		if len(v.Catalog.Values) != 0 {
			matchCatalog := v.Catalog.regexp.MatchString(a.Catalog)
			if matchCatalog && v.Catalog.IsExcludes {
				continue
			}
			if !matchCatalog && !v.Catalog.IsExcludes {
				continue
			}
		}
		if len(v.Schema.Values) != 0 {
			matchSchema := v.Schema.regexp.MatchString(a.Schema)
			if matchSchema && v.Schema.IsExcludes {
				continue
			}
			if !matchSchema && !v.Schema.IsExcludes {
				continue
			}
		}
		if len(v.Table.Values) != 0 {
			matchTable := v.Table.regexp.MatchString(a.Table)
			if matchTable && v.Table.IsExcludes {
				continue
			}
			if !matchTable && !v.Table.IsExcludes {
				continue
			}
		}

		return true
	}
	return false
}

func (p *Policy) getControlledActions(usersByGroup map[string][]string) ControlledActions {
	return ControlledActions{
		allowedActionsByUser: p.getAllAllowPolicyByUser(usersByGroup),
		deniedActionsByUser:  p.getAllDenyPoliciesByUser(usersByGroup),
	}
}

func (p *Policy) getAllAllowPolicyByUser(usersByGroup map[string][]string) map[string][]parser.Action {
	allowPoliciesItem := policyItemsToActionsByUser(p.PolicyItems, usersByGroup)
	excludeAllowPolicyItems := policyItemsToActionsByUser(p.AllowExceptions, usersByGroup)

	for user, actions := range excludeAllowPolicyItems {
		if _, ok := allowPoliciesItem[user]; !ok {
			continue
		}
		for action := range actions {
			delete(allowPoliciesItem[user], action)
		}
		if len(allowPoliciesItem[user]) == 0 {
			delete(allowPoliciesItem, user)
		}
	}

	result := map[string][]parser.Action{}
	for user, actionsMap := range allowPoliciesItem {
		actions := make([]parser.Action, 0, len(actionsMap))
		for action := range actionsMap {
			actions = append(actions, action)
		}
		result[user] = actions
	}
	return result
}

func (p *Policy) getAllDenyPoliciesByUser(usersByGroup map[string][]string) map[string][]parser.Action {
	denyPoliciesItem := policyItemsToActionsByUser(p.DenyPolicyItems, usersByGroup)
	excludeDenyPolicyItems := policyItemsToActionsByUser(p.DenyExceptions, usersByGroup)

	for user, actions := range excludeDenyPolicyItems {
		if _, ok := denyPoliciesItem[user]; !ok {
			continue
		}
		for action := range actions {
			delete(denyPoliciesItem[user], action)
		}
		if len(denyPoliciesItem[user]) == 0 {
			delete(denyPoliciesItem, user)
		}
	}

	result := map[string][]parser.Action{}
	for user, actionsMap := range denyPoliciesItem {
		actions := make([]parser.Action, 0, len(actionsMap))
		for action := range actionsMap {
			actions = append(actions, action)
		}
		result[user] = actions
	}
	return result
}

func globToRegex(pattern string) string {
	escaped := regexp.QuoteMeta(pattern)
	escaped = strings.ReplaceAll(escaped, "\\*", ".*")
	escaped = strings.ReplaceAll(escaped, "\\?", ".")
	return escaped
}

func patternsToRegex(patterns []string) string {
	var regexes []string
	for _, pat := range patterns {
		regexes = append(regexes, globToRegex(pat))
	}
	return strings.Join(regexes, "|")
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
		}
	}
	return p.Actions

}

func policyItemsToActionsByUser(items []PolicyItem, usersByGroup map[string][]string) map[string]map[parser.Action]struct{} {
	permissions := make(map[string]map[parser.Action]struct{})

	for _, item := range items {
		for _, user := range item.Users {
			user = strings.ToLower(user)
			if _, ok := permissions[user]; !ok {
				permissions[user] = make(map[parser.Action]struct{})
			}
			for _, action := range item.getPermissions() {
				permissions[user][action] = struct{}{}
			}
		}
		for _, group := range item.Groups {
			users, ok := usersByGroup[group]
			if !ok {
				continue
			}
			for _, user := range users {
				if _, ok := permissions[user]; !ok {
					permissions[user] = make(map[parser.Action]struct{})
				}
				for _, action := range item.getPermissions() {
					permissions[user][action] = struct{}{}
				}
			}
		}
	}

	return permissions
}

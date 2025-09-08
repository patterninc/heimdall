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
	Resources           Resource     `json:"resources"`
	AdditionalResources []Resource   `json:"additionalResources"`
	PolicyItems         []PolicyItem `json:"policyItems"`
	DenyPolicyItems     []PolicyItem `json:"denyPolicyItems"`
	AllowExceptions     []PolicyItem `json:"allowExceptions"`
	DenyExceptions      []PolicyItem `json:"denyExceptions"`
	ServiceType         string       `json:"serviceType"`

	supportedTables *regexp.Regexp // todo init that regexp
}

type ResourceField struct {
	Values     []string `json:"values"`
	IsExcludes bool     `json:"isExcludes"`
}

type Resource struct {
	Schema  ResourceField `json:"schema,omitempty"`
	Catalog ResourceField `json:"catalog,omitempty"`
	Table   ResourceField `json:"table,omitempty"`
	Column  ResourceField `json:"column,omitempty"`
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
	resourceRegexpParts := []string{}
	for _, resource := range append([]Resource{p.Resources}, p.AdditionalResources...) {
		resourceRegexpParts = append(resourceRegexpParts, resource.getMatchRegexp())
	}
	p.supportedTables = regexp.MustCompile("^(" + strings.Join(resourceRegexpParts, "|") + ")$")
	return nil
}

func (p *Policy) controlAnAccess(access parser.Access) bool {
	switch a := access.(type) {
	case *parser.TableAccess:
		return p.supportedTables.Match([]byte(a.QualifiedName()))
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

func (r *Resource) getMatchRegexp() string {
	catalogPart := ".*"
	if len(r.Catalog.Values) > 0 {
		catalogRegexp := patternsToRegex(r.Catalog.Values)
		if r.Catalog.IsExcludes {
			catalogPart = "(?!" + catalogRegexp + ").*"
		} else {
			catalogPart = "(" + catalogRegexp + ")"
		}
	}

	schemaPart := ".*"
	if len(r.Schema.Values) > 0 {
		schemaRegexp := patternsToRegex(r.Schema.Values)
		if r.Schema.IsExcludes {
			schemaPart = "(?!" + schemaRegexp + ").*"
		} else {
			schemaPart = "(" + schemaRegexp + ")"
		}
	}

	tablePart := ".*"
	if len(r.Table.Values) > 0 {
		tableRegexp := patternsToRegex(r.Table.Values)
		if r.Table.IsExcludes {
			tablePart = "(?!" + tableRegexp + ").*"
		} else {
			tablePart = "(" + tableRegexp + ")"
		}
	}

	return catalogPart + `\.` + schemaPart + `\.` + tablePart
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

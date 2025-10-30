package ranger

import (
	"regexp"

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
	PolicyItems         []PolicyItem `json:"policyItems"`
	DenyPolicyItems     []PolicyItem `json:"denyPolicyItems"`
	AllowExceptions     []PolicyItem `json:"allowExceptions"`
	DenyExceptions      []PolicyItem `json:"denyExceptions"`
	ServiceType         string       `json:"serviceType"`
}

type Resource struct {
	Schema  *ResourceField `json:"schema,omitempty"`
	Catalog *ResourceField `json:"catalog,omitempty"`
	Table   *ResourceField `json:"table,omitempty"`
	Column  *ResourceField `json:"column,omitempty"`
}

type ResourceField struct {
	Values     []string `json:"values"`
	IsExcludes bool     `json:"isExcludes"`
	regexp     *regexp.Regexp
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

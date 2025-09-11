package tests

import (
	"testing"

	"github.com/patterninc/heimdall/pkg/rbac/ranger"
)

func TestAllowPermissionsForGroups(t *testing.T) {
	tests := []testCase{
		{
			name:           "Policy allows all actions for group",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"all"})},
		},
		{
			name:           "Policy allows select action for group",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"select"})},
		},
		{
			name:           "Policy allows insert action for group, but query is select",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"insert"})},
		},
		{
			name:           "Policy allows multiple actions including select for group",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"insert", "select", "update"})},
		},
		{
			name:           "Policy allows multiple actions excluding select for group",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"insert", "update", "delete"})},
		},
		{
			name:           "No policy for group",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{},
		},
		{
			name:           "Policy allows select but query requires also insert",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"select"})},
		},
		{
			name:           "Policy allows all actions",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"all"})},
		},
		{
			name:           "Policy many actions and many are required",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultGroupAllowPolicy([]string{"delete", "insert", "select", "update"})},
		},
		{
			name:           "Policy exclude user from the select action",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultAllActionsGroupPolicyWithExcludeForDefaultGroup([]string{"select"})},
		},
		{
			name:           "Policy exclude user from the insert action, but query is select and insert",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultAllActionsGroupPolicyWithExcludeForDefaultGroup([]string{"insert"})},
		},
		{
			name:           "Policy exclude user from the insert action, but query is select ",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       []*ranger.Policy{getDefaultAllActionsGroupPolicyWithExcludeForDefaultGroup([]string{"insert"})},
		},
	}

	runTests(t, tests)
}

func TestDenyPermissionsForGroups(t *testing.T) {
	tests := []testCase{}

	runTests(t, tests)
}

func getDefaultGroupAllowPolicy(accessType []string) *ranger.Policy {
	return &ranger.Policy{
		ID:             1,
		GUID:           "policy-1",
		IsEnabled:      true,
		Name:           "Allow select for alice",
		PolicyType:     0,
		PolicyPriority: 1,
		Resources: &ranger.Resource{
			Catalog: &ranger.ResourceField{
				Values:     []string{"default_catalog"},
				IsExcludes: false,
			},
			Schema: &ranger.ResourceField{
				Values:     []string{"public"},
				IsExcludes: false,
			},
			Table: &ranger.ResourceField{
				Values:     []string{"table1"},
				IsExcludes: false,
			},
		},
		PolicyItems: []ranger.PolicyItem{
			{
				Groups: []string{testGroupName},
				Accesses: func() []ranger.Access {
					var accesses []ranger.Access
					for _, at := range accessType {
						accesses = append(accesses, ranger.Access{Type: at})
					}
					return accesses
				}(),
			},
		},
	}
}

func getDefaultAllActionsGroupPolicyWithExcludeForDefaultGroup(excludeAccess []string) *ranger.Policy {
	return &ranger.Policy{
		ID:             1,
		GUID:           "policy-1",
		IsEnabled:      true,
		Name:           "Allow select for alice",
		PolicyType:     0,
		PolicyPriority: 1,
		Resources: &ranger.Resource{
			Catalog: &ranger.ResourceField{
				Values:     []string{"default_catalog"},
				IsExcludes: false,
			},
			Schema: &ranger.ResourceField{
				Values:     []string{"public"},
				IsExcludes: false,
			},
			Table: &ranger.ResourceField{
				Values:     []string{"table1"},
				IsExcludes: false,
			},
		},
		PolicyItems: []ranger.PolicyItem{
			{
				Groups: []string{testGroupName},
				Accesses: []ranger.Access{
					{Type: "all"},
				},
			},
		},
		AllowExceptions: []ranger.PolicyItem{
			{
				Groups: []string{testGroupName},
				Accesses: func() []ranger.Access {
					var accesses []ranger.Access
					for _, ex := range excludeAccess {
						accesses = append(accesses, ranger.Access{Type: ex})
					}
					return accesses
				}(),
			},
		},
	}
}

package tests

import (
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/rbac/ranger"
)

func TestDenyPermissionsForUser(t *testing.T) {
	tests := []testCase{
		{
			name:           "Policy denies select action for user",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,

			policies: getAllowAllPolicyWithDenyForUser([]string{"select"}),
		},
		{
			name:           "Policy denies insert action for user",
			query:          "INSERT INTO default_catalog.public.table1 VALUES (1, 'data')",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,

			policies: getAllowAllPolicyWithDenyForUser([]string{"insert"}),
		},
		{
			name:           "Policy denies update action for user but query is select",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,

			policies: getAllowAllPolicyWithDenyForUser([]string{"update"}),
		},
		{
			name:           "Policy denies select and insert actions for user",
			query:          "INSERT INTO default_catalog.public.table1 VALUES (1, 'data')",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,

			policies: getAllowAllPolicyWithDenyAndExceptionForUser([]string{"select", "insert"}, []string{"all"}),
		},
		{
			name:           "Policy denies all actions for user but exception allows select",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,

			policies: getAllowAllPolicyWithDenyAndExceptionForUser([]string{"all"}, []string{"select"}),
		},
	}

	runTests(t, tests)
}

func TestAllowPermissionsForUser(t *testing.T) {
	tests := []testCase{
		{
			name:           "Policy allows all actions for user",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultUserAllowPolicy([]string{"all"})},
		},
		{
			name:           "Policy allows select action for user",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultUserAllowPolicy([]string{"select"})},
		},
		{
			name:           "Policy allows insert action for user, but query is select",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultUserAllowPolicy([]string{"insert"})},
		},
		{
			name:           "Policy allows multiple actions including select for user",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultUserAllowPolicy([]string{"insert", "select", "update"})},
		},
		{
			name:           "Policy allows multiple actions excluding select for user",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultUserAllowPolicy([]string{"insert", "update", "delete"})},
		},
		{
			name:           "No policy for user",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{},
		},
		{
			name:           "Policy allows select but query requires also insert",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultUserAllowPolicy([]string{"select"})},
		},
		{
			name:           "Policy allows all actions",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,

			policies: []*ranger.Policy{getDefaultUserAllowPolicy([]string{"all"})},
		},
		{
			name:           "Policy many actions and many are required",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultUserAllowPolicy([]string{"delete", "insert", "select", "update"})},
		},
		{
			name:           "Policy exclude user from the select action",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultAllActionsUserPolicyWithExcludeForDefaultUser([]string{"select"})},
		},
		{
			name:           "Policy exclude user from the insert action, but query is select and insert",
			query:          "INSERT INTO default_catalog.public.table1 as SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultAllActionsUserPolicyWithExcludeForDefaultUser([]string{"insert"})},
		},
		{
			name:           "Policy exclude user from the insert action, but query is select ",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getDefaultAllActionsUserPolicyWithExcludeForDefaultUser([]string{"insert"})},
		},
	}

	runTests(t, tests)

}

func getAllowAllPolicyWithDenyForUser(denyAccess []string) []*ranger.Policy {
	return []*ranger.Policy{
		{
			ID:             1,
			GUID:           "policy-1",
			IsEnabled:      true,
			Name:           "Allow select for alice",
			PolicyType:     0,
			PolicyPriority: 1,
			Resources: &ranger.Resource{
				Catalog: &ranger.ResourceField{
					RawValues:  []string{"default_catalog"},
					IsExcludes: false,
				},
				Schema: &ranger.ResourceField{
					RawValues:  []string{"public"},
					IsExcludes: false,
				},
				Table: &ranger.ResourceField{
					RawValues:  []string{"table1"},
					IsExcludes: false,
				},
			},
			PolicyItems: []*ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: []*ranger.Access{
						{Type: "all"},
					},
				},
			},
			DenyPolicyItems: []*ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: func() []*ranger.Access {
						var accesses []*ranger.Access
						for _, a := range denyAccess {
							accesses = append(accesses, &ranger.Access{Type: a})
						}
						return accesses
					}(),
				},
			},
		},
	}
}

func getAllowAllPolicyWithDenyAndExceptionForUser(denyAccess, exceptionAccess []string) []*ranger.Policy {
	return []*ranger.Policy{
		{
			ID:             1,
			GUID:           "policy-1",
			IsEnabled:      true,
			Name:           "Allow select for alice",
			PolicyType:     0,
			PolicyPriority: 1,
			Resources: &ranger.Resource{
				Catalog: &ranger.ResourceField{
					RawValues:  []string{"default_catalog"},
					IsExcludes: false,
				},
				Schema: &ranger.ResourceField{
					RawValues:  []string{"public"},
					IsExcludes: false,
				},
				Table: &ranger.ResourceField{
					RawValues:  []string{"table1"},
					IsExcludes: false,
				},
			},
			PolicyItems: []*ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: []*ranger.Access{
						{Type: "all"},
					},
				},
			},
			DenyPolicyItems: []*ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: func() []*ranger.Access {
						var accesses []*ranger.Access
						for _, a := range denyAccess {
							accesses = append(accesses, &ranger.Access{Type: a})
						}
						return accesses
					}(),
				},
			},
			DenyExceptions: []*ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: func() []*ranger.Access {
						var accesses []*ranger.Access
						for _, a := range exceptionAccess {
							accesses = append(accesses, &ranger.Access{Type: a})
						}
						return accesses
					}(),
				},
			},
		},
	}
}

func getDefaultUserAllowPolicy(accessType []string) *ranger.Policy {
	return &ranger.Policy{
		ID:             1,
		GUID:           "policy-1",
		IsEnabled:      true,
		Name:           "Allow select for alice",
		PolicyType:     0,
		PolicyPriority: 1,
		Resources: &ranger.Resource{
			Catalog: &ranger.ResourceField{
				RawValues:  []string{"default_catalog"},
				IsExcludes: false,
			},
			Schema: &ranger.ResourceField{
				RawValues:  []string{"public"},
				IsExcludes: false,
			},
			Table: &ranger.ResourceField{
				RawValues:  []string{"table1"},
				IsExcludes: false,
			},
		},
		PolicyItems: []*ranger.PolicyItem{
			{
				Users: []string{testUserName},
				Accesses: func() []*ranger.Access {
					var accesses []*ranger.Access
					for _, at := range accessType {
						accesses = append(accesses, &ranger.Access{Type: at})
					}
					return accesses
				}(),
			},
		},
	}
}

func getDefaultAllActionsUserPolicyWithExcludeForDefaultUser(excludeAccess []string) *ranger.Policy {
	return &ranger.Policy{
		ID:             1,
		GUID:           "policy-1",
		IsEnabled:      true,
		Name:           "Allow select for alice",
		PolicyType:     0,
		PolicyPriority: 1,
		Resources: &ranger.Resource{
			Catalog: &ranger.ResourceField{
				RawValues:  []string{"default_catalog"},
				IsExcludes: false,
			},
			Schema: &ranger.ResourceField{
				RawValues:  []string{"public"},
				IsExcludes: false,
			},
			Table: &ranger.ResourceField{
				RawValues:  []string{"table1"},
				IsExcludes: false,
			},
		},
		PolicyItems: []*ranger.PolicyItem{
			{
				Users: []string{testUserName},
				Accesses: []*ranger.Access{
					{Type: "all"},
				},
			},
		},
		AllowExceptions: []*ranger.PolicyItem{
			{
				Users: []string{testUserName},
				Accesses: func() []*ranger.Access {
					var accesses []*ranger.Access
					for _, ex := range excludeAccess {
						accesses = append(accesses, &ranger.Access{Type: ex})
					}
					return accesses
				}(),
			},
		},
	}
}

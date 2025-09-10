package tests

import (
	"testing"

	"github.com/patterninc/heimdall/pkg/rbac/ranger"
)

func TestDenyPermissionsForUser(t *testing.T) {
	tests := []testCase{
		{
			name:           "Policy denies select action for user",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       getAllowAllPolicyWithDenyForUser([]string{"select"}),
		},
		{
			name:           "Policy denies insert action for user",
			query:          "INSERT INTO default_catalog.public.table1 VALUES (1, 'data')",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       getAllowAllPolicyWithDenyForUser([]string{"insert"}),
		},
		{
			name:           "Policy denies update action for user but query is select",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       getAllowAllPolicyWithDenyForUser([]string{"update"}),
		},
		{
			name:           "Policy denies select and insert actions for user",
			query:          "INSERT INTO default_catalog.public.table1 VALUES (1, 'data')",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       getAllowAllPolicyWithDenyAndExceptionForUser([]string{"select", "insert"}, []string{"all"}),
		},
		{
			name:           "Policy denies all actions for user but exception allows select",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			groups:         testDefaultGroups,
			policies:       getAllowAllPolicyWithDenyAndExceptionForUser([]string{"all"}, []string{"select"}),
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
					Users: []string{testUserName},
					Accesses: []ranger.Access{
						{Type: "all"},
					},
				},
			},
			DenyPolicyItems: []ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: func() []ranger.Access {
						var accesses []ranger.Access
						for _, a := range denyAccess {
							accesses = append(accesses, ranger.Access{Type: a})
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
					Users: []string{testUserName},
					Accesses: []ranger.Access{
						{Type: "all"},
					},
				},
			},
			DenyPolicyItems: []ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: func() []ranger.Access {
						var accesses []ranger.Access
						for _, a := range denyAccess {
							accesses = append(accesses, ranger.Access{Type: a})
						}
						return accesses
					}(),
				},
			},
			DenyExceptions: []ranger.PolicyItem{
				{
					Users: []string{testUserName},
					Accesses: func() []ranger.Access {
						var accesses []ranger.Access
						for _, a := range exceptionAccess {
							accesses = append(accesses, ranger.Access{Type: a})
						}
						return accesses
					}(),
				},
			},
		},
	}
}

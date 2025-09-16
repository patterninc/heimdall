package tests

import (
	"testing"

	"github.com/patterninc/heimdall/pkg/rbac/ranger"
	"github.com/patterninc/heimdall/pkg/rbac/ranger/mocks"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

const (
	serviceName   = "test_service"
	testGroupName = "test_group"
	testUserName  = "test_user"
)

var (
	testDefaultUsers = map[string]*ranger.User{
		testUserName: {ID: 11, Name: testUserName, GroupNameList: []string{testGroupName}},
	}
)

type testCase struct {
	name           string
	query          string
	username       string
	expectedResult bool
	users          map[string]*ranger.User
	policies       []*ranger.Policy
}

func TestRangerPolicyCheck(t *testing.T) {
	tests := []testCase{
		{
			name:           "User with direct allow policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "alice",
			expectedResult: true,
			users: map[string]*ranger.User{
				"alice": {ID: 1, Name: "alice", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
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
							Users: []string{"alice"},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User with group allow policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "bob",
			expectedResult: true,
			users: map[string]*ranger.User{
				"bob": {ID: 2, Name: "bob", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             2,
					GUID:           "policy-2",
					IsEnabled:      true,
					Name:           "Allow select for testGroupName",
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
								{Type: "select"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User with deny policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "charlie",
			expectedResult: false,
			users: map[string]*ranger.User{
				"charlie": {ID: 3, Name: "charlie", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             3,
					GUID:           "policy-3",
					IsEnabled:      true,
					Name:           "Deny select for charlie",
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
							Users: []string{"charlie"},
							Accesses: []ranger.Access{
								{Type: "all"},
							},
						},
					},
					DenyPolicyItems: []ranger.PolicyItem{
						{
							Users: []string{"charlie"},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User without any policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "dave",
			expectedResult: false,
			users: map[string]*ranger.User{
				"dave": {ID: 4, Name: "dave", GroupNameList: []string{}},
			},
			policies: []*ranger.Policy{},
		},
		{
			name:           "User with conflicting allow and deny policies (deny should take precedence)",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "eve",
			expectedResult: false,
			users: map[string]*ranger.User{
				"eve": {ID: 5, Name: "eve", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             4,
					GUID:           "policy-4",
					IsEnabled:      true,
					Name:           "Allow select for eve",
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
							Users: []string{"eve"},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
				{
					ID:             5,
					GUID:           "policy-5",
					IsEnabled:      true,
					Name:           "Deny select for group3",
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
					DenyPolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{testGroupName},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User has different type of an access",
			query:          "SHOW TABLES FROM default_catalog.public",
			username:       "frank",
			expectedResult: false,
			users: map[string]*ranger.User{
				"frank": {ID: 6, Name: "frank", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             6,
					GUID:           "policy-6",
					IsEnabled:      true,
					Name:           "Allow show for group4",
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
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{testGroupName},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User allowed via regexp in table name",
			query:          "SELECT * FROM default_catalog.public.table_xyz",
			username:       "grace",
			expectedResult: true,
			users: map[string]*ranger.User{
				"grace": {ID: 7, Name: "grace", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             7,
					GUID:           "policy-7",
					IsEnabled:      true,
					Name:           "Allow select for testGroupName on tables matching regex",
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
							Values:     []string{"*"},
							IsExcludes: false,
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{testGroupName},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User denied via regexp in table name",
			query:          "SELECT * FROM default_catalog.public.table_abc",
			username:       "heidi",
			expectedResult: false,
			users: map[string]*ranger.User{
				"heidi": {ID: 8, Name: "heidi", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             8,
					GUID:           "policy-8",
					IsEnabled:      true,
					Name:           "Deny select for testGroupName on tables matching regex",
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
							Values:     []string{"table_*"},
							IsExcludes: false,
						},
					},
					DenyPolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{testGroupName},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{testGroupName},
							Accesses: []ranger.Access{
								{Type: "*"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User is exclude from allow policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "ivan",
			expectedResult: false,
			users: map[string]*ranger.User{
				"ivan": {ID: 9, Name: "ivan", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             9,
					GUID:           "policy-9",
					IsEnabled:      true,
					Name:           "Allow select for testGroupName excluding user ivan",
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
								{Type: "select"},
							},
						},
					},
					AllowExceptions: []ranger.PolicyItem{
						{
							Users: []string{"ivan"},
							Accesses: []ranger.Access{
								{Type: "all"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User is denied via group policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "judy",
			expectedResult: false,
			users: map[string]*ranger.User{
				"judy": {ID: 10, Name: "judy", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             10,
					GUID:           "policy-10",
					IsEnabled:      true,
					Name:           "Deny select for testGroupName excluding user judy",
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
							Users: []string{"judy"},
							Accesses: []ranger.Access{
								{Type: "all"},
							},
						},
					},
					DenyPolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{testGroupName},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User is exclude from deny policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "judy",
			expectedResult: true,
			users: map[string]*ranger.User{
				"judy": {ID: 10, Name: "judy", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             10,
					GUID:           "policy-10",
					IsEnabled:      true,
					Name:           "Deny select for testGroupName excluding user judy",
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
							Users: []string{"judy"},
							Accesses: []ranger.Access{
								{Type: "all"},
							},
						},
					},
					DenyPolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group8"},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
					DenyExceptions: []ranger.PolicyItem{
						{
							Users: []string{"judy"},
							Accesses: []ranger.Access{
								{Type: "all"},
							},
						},
					},
				},
			},
		},
		{
			name:           "User is allowed when resource has excludes",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "kate",
			expectedResult: true,
			users: map[string]*ranger.User{
				"kate": {ID: 11, Name: "kate", GroupNameList: []string{testGroupName}},
			},
			policies: []*ranger.Policy{
				{
					ID:             11,
					GUID:           "policy-11",
					IsEnabled:      true,
					Name:           "Allow select for testGroupName excluding schema 'internal'",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: &ranger.Resource{
						Catalog: &ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: &ranger.ResourceField{
							Values:     []string{"internal"},
							IsExcludes: true,
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
								{Type: "select"},
							},
						},
					},
				},
			},
		},
	}

	runTests(t, tests)
}

// TestResourcesSelection tests the resource selection logic in Ranger policies.
// In this tests users always have all permissions
func TestResourcesSelection(t *testing.T) {
	tests := []testCase{
		{
			name:           "Policy doesn't control the resource, different table name",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("default_catalog", "public", "table2"), nil)},
		},
		{
			name:           "Policy doesn't control the resource, different schema name",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("default_catalog", "private", "table1"), nil)},
		},
		{
			name:           "Policy doesn't control the resource, different catalog name",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("not_default_catalog", "public", "table1"), nil)},
		},
		{
			name:           "Policy and subpolicy doesn't control the resource, different catalog/table name",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("not_default_catalog", "public", "table1"), createResource("default_catalog", "public", "table2"))},
		},
		{
			name:           "Policy controls the resource, catalog is regexp",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("default_*", "public", "table1"), nil)},
		},
		{
			name:           "Policy controls the resource, schema is regexp",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("default_catalog", "p*c", "table1"), nil)},
		},
		{
			name:           "Policy controls the resource, table is regexp",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("default_catalog", "public", "t*l*"), nil)},
		},
		{
			name:           "Policy controls the resource, table is regexp",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("default_catalog", "public", "t*l*"), nil)},
		},
		{
			name:           "Policy controls the resource, exact match",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResource("default_catalog", "public", "table1"), nil)},
		},
		{
			name:           "Policy controls the resource, catalog is exclude",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForCatalog("catalog", "public", "table1", true), nil)},
		},
		{
			name:           "Policy controls the resource, catalog is exclude regexp",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForCatalog("catalo*", "public", "table1", true), nil)},
		},
		{
			name:           "Policy doesn't control the resource, catalog is exclude regexp but match",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForCatalog("defa*", "public", "table1", true), nil)},
		},
		{
			name:           "Policy and subpolicy control the resource, catalog is exclude but subpolicy match",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForCatalog("defa*", "public", "table1", true), createResource("default_catalog", "public", "table1"))},
		},
		{
			name:           "Policy doesn't control the resource, schema is exclude and match",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForSchema("defa*", "public", "table1", true), nil)},
		},
		{
			name:           "Policy controls the resource, schema is exclude but not match",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForSchema("defa*", "privat*", "table1", true), nil)},
		},
		{
			name:           "Policy and subpolicy control the resource, schema is exclude but subpolicy match",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForSchema("defa*", "public", "table1", true), createResource("default_catalog", "public", "table1"))},
		},
		{
			name:           "Policy doesn't control the resources table is excluded",
			query:          "SELECT * from default_catalog.public.table1",
			username:       testUserName,
			expectedResult: false,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForTable("defau*", "public", "table1", true), nil)},
		},
		{
			name:           "Policy does control the resources, table is excluded but doesn't match",
			query:          "SELECT * from default_catalog.public.table1",
			username:       testUserName,
			expectedResult: true,
			users:          testDefaultUsers,
			policies:       []*ranger.Policy{getAllowAllPolicy(createResourceWithExcludeOptionForTable("defau*", "public", "table2", true), nil)},
		},
	}
	runTests(t, tests)

}

func runTests(t *testing.T, tests []testCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rbac := &ranger.Ranger{
				AccessReceiver: trino.NewTrinoAccessReceiver("default_catalog"),
				Client:         getMockRangerClient(tt.users, tt.policies),
				ServiceName:    serviceName,
			}
			rbac.SyncState()
			actualResult, err := rbac.HasAccess(tt.username, tt.query)
			if err != nil {
				t.Errorf("error checking access for %q: %v", tt.name, err)
				return
			}
			if actualResult != tt.expectedResult {
				t.Errorf("unexpected result for %q: got %v, want %v", tt.name, actualResult, tt.expectedResult)
			}
		})
	}
}

func createResourceWithExcludeOptionForTable(catalogs, schemas, table string, excludeTable bool) *ranger.Resource {
	return &ranger.Resource{
		Catalog: &ranger.ResourceField{
			Values:     []string{catalogs},
			IsExcludes: false,
		},
		Schema: &ranger.ResourceField{
			Values:     []string{schemas},
			IsExcludes: false,
		},
		Table: &ranger.ResourceField{
			Values:     []string{table},
			IsExcludes: excludeTable,
		},
	}
}

func createResourceWithExcludeOptionForSchema(catalog, schema, table string, excludeSchema bool) *ranger.Resource {
	return &ranger.Resource{
		Catalog: &ranger.ResourceField{
			Values:     []string{catalog},
			IsExcludes: false,
		},
		Schema: &ranger.ResourceField{
			Values:     []string{schema},
			IsExcludes: excludeSchema,
		},
		Table: &ranger.ResourceField{
			Values:     []string{table},
			IsExcludes: false,
		},
	}
}

func createResourceWithExcludeOptionForCatalog(catalogs, schemas, tables string, excludeCatalog bool) *ranger.Resource {
	return &ranger.Resource{
		Catalog: &ranger.ResourceField{
			Values:     []string{catalogs},
			IsExcludes: excludeCatalog,
		},
		Schema: &ranger.ResourceField{
			Values:     []string{schemas},
			IsExcludes: false,
		},
		Table: &ranger.ResourceField{
			Values:     []string{tables},
			IsExcludes: false,
		},
	}
}

func createResource(catalogs, schemas, tables string) *ranger.Resource {
	return &ranger.Resource{
		Catalog: &ranger.ResourceField{
			Values:     []string{catalogs},
			IsExcludes: false,
		},
		Schema: &ranger.ResourceField{
			Values:     []string{schemas},
			IsExcludes: false,
		},
		Table: &ranger.ResourceField{
			Values:     []string{tables},
			IsExcludes: false,
		},
	}
}

func getAllowAllPolicy(resource *ranger.Resource, additionalResource *ranger.Resource) *ranger.Policy {
	var additionalResources []*ranger.Resource
	if additionalResource != nil {
		additionalResources = []*ranger.Resource{additionalResource}
	}
	return &ranger.Policy{
		ID:                  1,
		GUID:                "policy-1",
		IsEnabled:           true,
		Name:                "Allow select for alice",
		PolicyType:          0,
		PolicyPriority:      1,
		Resources:           resource,
		AdditionalResources: additionalResources,
		PolicyItems: []ranger.PolicyItem{
			{
				Users: []string{testUserName},
				Accesses: []ranger.Access{
					{Type: "all"},
				},
			},
		},
	}
}

func getMockRangerClient(users map[string]*ranger.User, policies []*ranger.Policy) *ranger.ClientWrapper {
	m := new(mocks.Client)
	m.On("GetUsers").Return(users, nil)
	m.On("GetPolicies", serviceName).Return(policies, nil)
	return &ranger.ClientWrapper{Client: m}
}

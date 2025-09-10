package tests

import (
	"testing"

	"github.com/patterninc/heimdall/pkg/rbac/ranger"
	"github.com/patterninc/heimdall/pkg/rbac/ranger/mocks"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

const (
	serviceName = "test_service"
)

func TestRangerPolicyCheck(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		username       string
		expectedResult bool
		users          map[string]*ranger.User
		groups         map[string]*ranger.Group
		policies       []*ranger.Policy
	}{
		{
			name:           "User with direct allow policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "alice",
			expectedResult: true,
			users: map[string]*ranger.User{
				"alice": {ID: 1, Name: "alice", GroupIdList: []int64{1}},
			},
			groups: map[string]*ranger.Group{
				"group1": {ID: 1, Name: "group1"},
			},
			policies: []*ranger.Policy{
				{
					ID:             1,
					GUID:           "policy-1",
					IsEnabled:      true,
					Name:           "Allow select for alice",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
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
				"bob": {ID: 2, Name: "bob", GroupIdList: []int64{1}},
			},
			groups: map[string]*ranger.Group{
				"group1": {ID: 1, Name: "group1"},
			},
			policies: []*ranger.Policy{
				{
					ID:             2,
					GUID:           "policy-2",
					IsEnabled:      true,
					Name:           "Allow select for group1",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
							Values:     []string{"table1"},
							IsExcludes: false,
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group1"},
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
				"charlie": {ID: 3, Name: "charlie", GroupIdList: []int64{2}},
			},
			groups: map[string]*ranger.Group{
				"group2": {ID: 2, Name: "group2"},
			},
			policies: []*ranger.Policy{
				{
					ID:             3,
					GUID:           "policy-3",
					IsEnabled:      true,
					Name:           "Deny select for charlie",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
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
				"dave": {ID: 4, Name: "dave", GroupIdList: []int64{}},
			},
			groups:   map[string]*ranger.Group{},
			policies: []*ranger.Policy{},
		},
		{
			name:           "User with conflicting allow and deny policies (deny should take precedence)",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "eve",
			expectedResult: false,
			users: map[string]*ranger.User{
				"eve": {ID: 5, Name: "eve", GroupIdList: []int64{3}},
			},
			groups: map[string]*ranger.Group{
				"group3": {ID: 3, Name: "group3"},
			},
			policies: []*ranger.Policy{
				{
					ID:             4,
					GUID:           "policy-4",
					IsEnabled:      true,
					Name:           "Allow select for eve",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
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
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
							Values:     []string{"table1"},
							IsExcludes: false,
						},
					},
					DenyPolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group3"},
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
				"frank": {ID: 6, Name: "frank", GroupIdList: []int64{4}},
			},
			groups: map[string]*ranger.Group{
				"group4": {ID: 4, Name: "group4"},
			},
			policies: []*ranger.Policy{
				{
					ID:             6,
					GUID:           "policy-6",
					IsEnabled:      true,
					Name:           "Allow show for group4",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group4"},
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
				"grace": {ID: 7, Name: "grace", GroupIdList: []int64{5}},
			},
			groups: map[string]*ranger.Group{
				"group5": {ID: 5, Name: "group5"},
			},
			policies: []*ranger.Policy{
				{
					ID:             7,
					GUID:           "policy-7",
					IsEnabled:      true,
					Name:           "Allow select for group5 on tables matching regex",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
							Values:     []string{"*"},
							IsExcludes: false,
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group5"},
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
				"heidi": {ID: 8, Name: "heidi", GroupIdList: []int64{6}},
			},
			groups: map[string]*ranger.Group{
				"group6": {ID: 6, Name: "group6"},
			},
			policies: []*ranger.Policy{
				{
					ID:             8,
					GUID:           "policy-8",
					IsEnabled:      true,
					Name:           "Deny select for group6 on tables matching regex",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
							Values:     []string{"table_*"},
							IsExcludes: false,
						},
					},
					DenyPolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group6"},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group6"},
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
				"ivan": {ID: 9, Name: "ivan", GroupIdList: []int64{7}},
			},
			groups: map[string]*ranger.Group{
				"group7": {ID: 7, Name: "group7"},
			},
			policies: []*ranger.Policy{
				{
					ID:             9,
					GUID:           "policy-9",
					IsEnabled:      true,
					Name:           "Allow select for group7 excluding user ivan",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
							Values:     []string{"table1"},
							IsExcludes: false,
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group7"},
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
				"judy": {ID: 10, Name: "judy", GroupIdList: []int64{8}},
			},
			groups: map[string]*ranger.Group{
				"group8": {ID: 8, Name: "group8"},
			},
			policies: []*ranger.Policy{
				{
					ID:             10,
					GUID:           "policy-10",
					IsEnabled:      true,
					Name:           "Deny select for group8 excluding user judy",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
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
				},
			},
		},
		{
			name:           "User is exclude from deny policy",
			query:          "SELECT * FROM default_catalog.public.table1",
			username:       "judy",
			expectedResult: true,
			users: map[string]*ranger.User{
				"judy": {ID: 10, Name: "judy", GroupIdList: []int64{8}},
			},
			groups: map[string]*ranger.Group{
				"group8": {ID: 8, Name: "group8"},
			},
			policies: []*ranger.Policy{
				{
					ID:             10,
					GUID:           "policy-10",
					IsEnabled:      true,
					Name:           "Deny select for group8 excluding user judy",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"public"},
							IsExcludes: false,
						},
						Table: ranger.ResourceField{
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
				"kate": {ID: 11, Name: "kate", GroupIdList: []int64{9}},
			},
			groups: map[string]*ranger.Group{
				"group9": {ID: 9, Name: "group9"},
			},
			policies: []*ranger.Policy{
				{
					ID:             11,
					GUID:           "policy-11",
					IsEnabled:      true,
					Name:           "Allow select for group9 excluding schema 'internal'",
					PolicyType:     0,
					PolicyPriority: 1,
					Resources: ranger.Resource{
						Catalog: ranger.ResourceField{
							Values:     []string{"default_catalog"},
							IsExcludes: false,
						},
						Schema: ranger.ResourceField{
							Values:     []string{"internal"},
							IsExcludes: true,
						},
						Table: ranger.ResourceField{
							Values:     []string{"table1"},
							IsExcludes: false,
						},
					},
					PolicyItems: []ranger.PolicyItem{
						{
							Groups: []string{"group9"},
							Accesses: []ranger.Access{
								{Type: "select"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			rbac := &ranger.ApacheRanger{
				AccessReceiver: trino.NewTrinoAccessReceiver("default_catalog"),
				Client:         getMockRangerClient(tt.users, tt.groups, tt.policies),
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

func getMockRangerClient(users map[string]*ranger.User, groups map[string]*ranger.Group, policies []*ranger.Policy) *mocks.Client {
	m := new(mocks.Client)
	m.On("GetUsers").Return(users, nil)
	m.On("GetGroups").Return(groups, nil)
	m.On("GetPolicies", serviceName).Return(policies, nil)
	return m
}

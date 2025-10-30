package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser"
	"github.com/patterninc/heimdall/internal/pkg/sql/parser/trino"
)

func TestParseSQLDelete(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []parser.Access
	}{
		{
			name:  "simple DELETE from single table without catalog",
			query: "DELETE FROM public.sales",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.DELETE,
				},
			},
		},
		{
			name:  "simple DELETE from single table with catalog",
			query: "DELETE FROM catalog_name.schema_name.sales",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "schema_name",
					Catalog: "catalog_name",
					Act:     parser.DELETE,
				},
			},
		},
		{
			name:  "simple DELETE from single table with catalog",
			query: "DELETE FROM catalog_name.schema_name.sales",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "schema_name",
					Catalog: "catalog_name",
					Act:     parser.DELETE,
				},
			},
		},
		{
			name:  "DELETE with WHERE clause",
			query: "DELETE FROM public.sales WHERE id = 1",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.DELETE,
				},
			},
		},
		{
			name:  "DELETE with subquery in where clause",
			query: "DELETE FROM public.sales WHERE id IN (SELECT id FROM catalog_name.public.second_table WHERE amount > 100)",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.DELETE,
				},
				&parser.TableAccess{
					Table:   "second_table",
					Schema:  "public",
					Catalog: "catalog_name",
					Act:     parser.SELECT,
				},
			},
		},
		{
			name: "DELETE with join in select clause",
			query: `DELETE FROM analytics.public.orders WHERE customer_id IN (
					SELECT c.customer_id
					FROM public.customers c
					JOIN public.blacklist b ON c.email = b.email
					WHERE b.active = true
				);`,
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "orders",
					Schema:  "public",
					Catalog: "analytics",
					Act:     parser.DELETE,
				},
				&parser.TableAccess{
					Table:   "customers",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.SELECT,
				},
				&parser.TableAccess{
					Table:   "blacklist",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.SELECT,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := trino.NewTrinoAccessReceiver(defaultCatalog)
			result, err := receiver.ParseAccess(tt.query)
			if err != nil {
				t.Errorf("Unexpected error in test %s: %v", tt.name, err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				expectedBytes, _ := json.Marshal(tt.expected)
				resultBytes, _ := json.Marshal(result)
				t.Errorf("test %s: expected %+v, got %+v", tt.name, string(expectedBytes), string(resultBytes))
			}
		})
	}

}

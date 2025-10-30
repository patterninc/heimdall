package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser/trino"
	"github.com/patterninc/heimdall/pkg/sql/parser"
)

func TestParseSQLUpdate(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []parser.Access
	}{
		{
			name:  "simple UPDATE on single table without catalog",
			query: "UPDATE public.sales SET amount = 200 WHERE id = 1",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.UPDATE,
				},
			},
		},
		{
			name:  "UPDATE with catalog and schema",
			query: "UPDATE my_catalog.public.sales SET amount = 200 WHERE id = 1",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "public",
					Catalog: "my_catalog",
					Act:     parser.UPDATE,
				},
			},
		},
		{
			name: "UPDATE with subquery",
			query: `UPDATE public.employees e
					SET rank = 'senior'
					WHERE EXISTS (
						SELECT 1 FROM catalog_name.public.promotions p WHERE p.employee_id = e.employee_id
					)`,
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "employees",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.UPDATE,
				},
				&parser.TableAccess{
					Table:   "promotions",
					Schema:  "public",
					Catalog: "catalog_name",
					Act:     parser.SELECT,
				},
			},
		},
		{
			name: "Update with join in select statement",
			query: `UPDATE analytics.public.payments
					SET flagged = true
					WHERE customer_id IN (
						SELECT c.customer_id
						FROM public.customers c
						JOIN public.review_list r ON c.email = r.email
						WHERE r.needs_review = true
					)`,
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "payments",
					Schema:  "public",
					Catalog: "analytics",
					Act:     parser.UPDATE,
				},
				&parser.TableAccess{
					Table:   "customers",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.SELECT,
				},
				&parser.TableAccess{
					Table:   "review_list",
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

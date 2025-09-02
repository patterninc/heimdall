package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

func TestParseSQLUpdate(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*parser.TableAccess
	}{
		{
			name:  "simple UPDATE on single table without catalog",
			query: "UPDATE public.sales SET amount = 200 WHERE id = 1",
			expected: []*parser.TableAccess{
				{
					Name:    "sales",
					Schema:  "public",
					Catalog: defaultCatalog,
					Access:  parser.UPDATE,
				},
			},
		},
		{
			name:  "UPDATE with catalog and schema",
			query: "UPDATE my_catalog.public.sales SET amount = 200 WHERE id = 1",
			expected: []*parser.TableAccess{
				{
					Name:    "sales",
					Schema:  "public",
					Catalog: "my_catalog",
					Access:  parser.UPDATE,
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
			expected: []*parser.TableAccess{
				{
					Name:    "employees",
					Schema:  "public",
					Catalog: defaultCatalog,
					Access:  parser.UPDATE,
				},
				{
					Name:    "promotions",
					Schema:  "public",
					Catalog: "catalog_name",
					Access:  parser.SELECT,
				},
			},
		},
		{
			name:"Update with join in select statement",
			query:`UPDATE analytics.public.payments
					SET flagged = true
					WHERE customer_id IN (
						SELECT c.customer_id
						FROM public.customers c
						JOIN public.review_list r ON c.email = r.email
						WHERE r.needs_review = true
					)`,
			expected: []*parser.TableAccess{
				{
					Name:    "payments",
					Schema:  "public",
					Catalog: "analytics",
					Access:  parser.UPDATE,
				},
				{
					Name:    "customers",
					Schema:  "public",
					Catalog: defaultCatalog,
					Access:  parser.SELECT,
				},
				{
					Name:    "review_list",
					Schema:  "public",
					Catalog: defaultCatalog,
					Access:  parser.SELECT,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := trino.NewTrinoAccessReceiver(defaultCatalog)
			result, err := receiver.ParseTableAccess(tt.query)
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

package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

func TestParseSQLCreateTable(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*parser.TableAccess
	}{
		{
			name:  "simple CREATE TABLE without catalog",
			query: "CREATE TABLE public.sales (id INT, amount DECIMAL)",
			expected: []*parser.TableAccess{{
				Name:    "sales",
				Schema:  "public",
				Catalog: defaultCatalog,
				Access:  parser.CREATE,
			}},
		},
		{
			name: "Create table with partitions",
			query: `CREATE TABLE default.logs (
				log_id BIGINT,
				log_message VARCHAR,
				log_date DATE
			)
			WITH (
				partitioned_by = ARRAY['log_date']
			);`,
			expected: []*parser.TableAccess{{
				Name:    "logs",
				Schema:  "default",
				Catalog: defaultCatalog,
				Access:  parser.CREATE,
			}},
		},
		{
			name: "Create table as select",
			query: `
			CREATE TABLE test_catalog.default.top_customers AS
				SELECT customer_id, SUM(amount) AS total_spent
				FROM catalog.public.orders
				GROUP BY customer_id
				HAVING SUM(amount) > 1000;`,
			expected: []*parser.TableAccess{
				{
					Name:    "top_customers",
					Schema:  "default",
					Catalog: "test_catalog",
					Access:  parser.CREATE,
				},
				{
					Name:    "orders",
					Schema:  "public",
					Catalog: "catalog",
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

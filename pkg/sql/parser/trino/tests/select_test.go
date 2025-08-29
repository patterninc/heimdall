package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

const (
	defaultCatalog = "default_catalog"
)

func TestParseSQLSelect(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*parser.TableAccess
	}{
		{
			name:  "Single table select",
			query: "SELECT * FROM schema1.table1",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Single table select with alias",
			query: "SELECT * FROM schema1.table1 as t1",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Multiple tables",
			query: "SELECT * FROM schema1.table1, testCatalog.schema2.table2",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: "testCatalog",
				},
			},
		},
		{
			name:  "Join tables",
			query: "SELECT * FROM ( SELECT * FROM schema1.table1, schema2.table2 )as b",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Join tables",
			query: "SELECT * FROM ( SELECT * FROM test_catalog.schema1.table1, schema2.table2 )as b",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: "test_catalog",
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Join with ON clause",
			query: "SELECT a.col1, b.col2 FROM schema1.table1 a JOIN schema2.table2 b ON a.id = b.id",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Subquery in FROM clause",
			query: "SELECT * FROM schema1.table1 WHERE id IN (SELECT id FROM schema2.table2)",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Subquery with alias",
			query: "SELECT * FROM (SELECT * FROM schema.table WHERE id IN (SELECT id FROM schema2.table2)) as sub",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema",
					Name:    "table",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Complex query with multiple joins and subqueries",
			query: "SELECT * FROM schema.table LEFT JOIN schema2.table2 ON schema.table.id = schema2.table2.id",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema",
					Name:    "table",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "EXISTS subquery",
			query: "SELECT * FROM schema1.table1 WHERE EXISTS (SELECT 1 FROM schema2.table2 WHERE schema2.table2.id = schema1.table1.id)",
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name: "Select with WHERE, AND, OR and IS clauses",
			query: `
				SELECT * FROM schema1.table1 as t1
				WHERE col1 = (SELECT col1 FROM schema2.table2 WHERE col2 = 'value2')
				AND (
					col2 = (SELECT col2 FROM schema3.table3 WHERE col2 = 'value2') OR
					col3 IS NOT NULL
				)
				AND col4 NOT IN (SELECT col4 FROM schema4.table4 WHERE col4 = 'value4')
			`,
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema2",
					Name:    "table2",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema3",
					Name:    "table3",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "schema4",
					Name:    "table4",
					Catalog: defaultCatalog,
				},
			},
		},

		{
			name: "Select WITH clause",
			query: `
			WITH recent_orders AS (
				SELECT
					order_id,
					customer_id,
					order_date
				FROM
					sales.orders
				WHERE
					order_date >= DATE '2024-01-01'
			)
			SELECT
				recent_orders.order_id,
				customers.customer_name
			FROM
				recent_orders
				JOIN sales.customers ON recent_orders.customer_id = customers.customer_id
			WHERE
				customers.is_active = TRUE;`,
			expected: []*parser.TableAccess{
				{
					Access:  parser.SELECT,
					Schema:  "sales",
					Name:    "orders",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "sales",
					Name:    "customers",
					Catalog: defaultCatalog,
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

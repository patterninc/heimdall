package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

func TestParseSQLInsert(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*parser.TableAccess
	}{
		{
			name:  "Insert into single table",
			query: "INSERT INTO schema1.table1 (col1, col2) VALUES ('value1', 'value2')",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Insert with select subquery",
			query: "INSERT INTO schema1.table1 (col1) SELECT col1 FROM schema2.table2",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
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
			name:  "Insert into table with schema and columns",
			query: "INSERT INTO schema1.table1 (col1, col2, col3) VALUES (1, 2, 3)",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Insert into table with schema and columns",
			query: "INSERT INTO schema1.table1 (col1, col2, col3) VALUES (1, 2, 3)",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Insert ignore",
			query: "INSERT IGNORE INTO schema1.table1 (col1) VALUES ('a')",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Insert with on duplicate key update",
			query: "INSERT INTO schema1.table1(col1) VALUES ('a') ON DUPLICATE KEY UPDATE col1 = 'b'",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Insert multiple rows",
			query: "INSERT INTO schema1.table1(col1, col2) VALUES ('a', 'b'), ('c', 'd')",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "schema1",
					Name:    "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Insert select with join",
			query: "INSERT INTO schema1.table1(col1) SELECT t1.col1 FROM schema2.table2 t1 JOIN schema3.table3 t2 ON t1.id = t2.id",
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
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
			},
		},
		{
			name: "Insert from SELECT with WHERE, AND, OR and IS clauses",
			query: `
				INSERT INTO schema1.table1(col1)
				SELECT col1 FROM schema2.table2
				WHERE col2 = (SELECT col2 FROM schema3.table3 WHERE col2 = 'value2')
				AND (
					col2 = (SELECT col2 FROM schema3.table3 WHERE col2 = 'value2') OR
					col3 IS NOT NULL
				)
				AND col4 NOT IN (SELECT col4 FROM schema4.table4 WHERE col4 = 'value4')
			`,
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
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

package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/sql"
)

func TestParseSQLInsert(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*sql.ResourceAccess
	}{
		{
			name:  "Insert into single table",
			query: "INSERT INTO myschema.mytable (col1, col2) VALUES ('value1', 'value2')",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Insert with select subquery",
			query: "INSERT INTO myschema.mytable (col1) SELECT col1 FROM myschema2.mytable2",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Insert into table with schema and columns",
			query: "INSERT INTO myschema.mytable (col1, col2, col3) VALUES (1, 2, 3)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Insert ignore",
			query: "INSERT IGNORE INTO myschema.mytable (col1) VALUES ('a')",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Insert with on duplicate key update",
			query: "INSERT INTO myschema.mytable (col1) VALUES ('a') ON DUPLICATE KEY UPDATE col1 = 'b'",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:     "Insert into table without schema",
			query:    "INSERT INTO mytable (col1) VALUES ('a')",
			expected: nil,
		},
		{
			name:  "Insert multiple rows",
			query: "INSERT INTO myschema.mytable (col1, col2) VALUES ('a', 'b'), ('c', 'd')",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Insert select with join",
			query: "INSERT INTO myschema.mytable (col1) SELECT t1.col1 FROM myschema2.mytable2 t1 JOIN myschema3.mytable3 t2 ON t1.id = t2.id",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
				{
					AccessType: "select",
					Schema:     "myschema3",
					Table:      "mytable3",
				},
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name: "Insert from SELECT with WHERE, AND, OR and IS clauses",
			query: `
				INSERT INTO myschema.mytable (col1)
				SELECT col1 FROM myschema2.mytable2
				WHERE col2 = (SELECT col2 FROM myschema3.mytable3 WHERE col2 = 'value2')
				AND (
					col2 = (SELECT col2 FROM myschema3.mytable3 WHERE col2 = 'value2') OR 
					col3 IS NOT NULL
				)
				AND col4 NOT IN (SELECT col4 FROM myschema4.mytable4 WHERE col4 = 'value4')
			`,
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
				{
					AccessType: "select",
					Schema:     "myschema3",
					Table:      "mytable3",
				},
				{
					AccessType: "select",
					Schema:     "myschema3",
					Table:      "mytable3",
				},
				{
					AccessType: "select",
					Schema:     "myschema4",
					Table:      "mytable4",
				},
				{
					AccessType: "insert",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sql.ParseSQLToResourceAccess(tt.query)
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

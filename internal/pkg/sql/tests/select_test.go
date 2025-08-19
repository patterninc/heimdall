package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/sql"
)

func TestParseSQLSelect(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*sql.ResourceAccess
	}{
		{
			name:  "Single table",
			query: "SELECT * FROM myschema.mytable",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Multiple tables",
			query: "SELECT * FROM myschema.mytable, myschema2.mytable2",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
			},
		},
		{
			name:  "Join tables",
			query: "SELECT * FROM ( SELECT * FROM myschema.mytable, myschema2.mytable2 )as b",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
			},
		},
		{
			name:  "Join with ON clause",
			query: "SELECT a.col1, b.col2 FROM myschema.mytable a JOIN myschema2.mytable2 b ON a.id = b.id",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
			},
		},
		{
			name:  "Subquery in FROM clause",
			query: "SELECT * FROM myschema.mytable WHERE id IN (SELECT id FROM myschema2.mytable2)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
			},
		},
		{
			name:  "Subquery with alias",
			query: "SELECT * FROM (SELECT * FROM myschema.mytable WHERE id IN (SELECT id FROM myschema2.mytable2)) as sub",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
			},
		},
		{
			name:  "Complex query with multiple joins and subqueries",
			query: "SELECT * FROM myschema.mytable LEFT JOIN myschema2.mytable2 ON myschema.mytable.id = myschema2.mytable2.id",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
			},
		},
		{
			name:  "EXISTS subquery",
			query: "SELECT * FROM myschema.mytable WHERE EXISTS (SELECT 1 FROM myschema2.mytable2 WHERE myschema2.mytable2.id = myschema.mytable.id)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "select",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
			},
		},
		{
			name: "Select with WHERE, AND, OR and IS clauses",
			query: `
				SELECT * FROM myschema.mytable  
				WHERE col1 = (SELECT col1 FROM myschema2.mytable2 WHERE col2 = 'value2')
				AND (
					col2 = (SELECT col2 FROM myschema3.mytable3 WHERE col2 = 'value2') OR 
					col3 IS NOT NULL
				)
				AND col4 NOT IN (SELECT col4 FROM myschema4.mytable4 WHERE col4 = 'value4')
			`,
			expected: []*sql.ResourceAccess{
				{
					AccessType: "select",
					Schema:     "myschema",
					Table:      "mytable",
				},
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
					Schema:     "myschema4",
					Table:      "mytable4",
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

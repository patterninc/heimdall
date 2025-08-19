package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/sql"
)

func TestParseSQLDelete(t *testing.T) {
	type testCase struct {
		name     string
		query    string
		expected []*sql.ResourceAccess
	}
	tests := []testCase{
		{
			name:  "Delete from single table",
			query: "DELETE FROM myschema.mytable WHERE col1 = 'value1'",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Delete with join",
			query: "DELETE a FROM myschema.mytable a JOIN myschema2.mytable2 b ON a.id = b.id WHERE b.col1 = 'value1'",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
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
			name:  "Delete with subquery",
			query: "DELETE FROM myschema.mytable WHERE id IN (SELECT id FROM myschema2.mytable2)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
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
			name:  "Delete with multiple conditions",
			query: "DELETE FROM myschema.mytable WHERE col1 = 'value1' AND col2 = 'value2'",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Delete with complex condition",
			query: "DELETE FROM myschema.mytable WHERE col1 = 'value1' OR (col2 = 'value2' AND col3 = 'value3')",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Delete without where clause",
			query: "DELETE FROM myschema.mytable",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Delete with alias and subquery",
			query: "DELETE a FROM myschema.mytable a WHERE a.id IN (SELECT b.id FROM myschema2.mytable2 b WHERE b.flag = 1)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
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
			name:  "Delete with exists subquery",
			query: "DELETE FROM myschema.mytable WHERE EXISTS (SELECT 1 FROM myschema2.mytable2 WHERE mytable.id = mytable2.id)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
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
			name:  "Delete with using clause",
			query: "DELETE FROM myschema.mytable USING myschema2.mytable2 WHERE myschema.mytable.id = myschema2.mytable2.id",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
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
			name:  "Delete from table with quoted identifiers",
			query: "DELETE FROM `myschema`.`mytable` WHERE `col1` = `value1`",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Delete with limit",
			query: "DELETE FROM myschema.mytable WHERE col1 = 'value1' LIMIT 10",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Delete with order by",
			query: "DELETE FROM myschema.mytable WHERE col1 = 'value1' ORDER BY col2 DESC",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name: "Delete with multiple tables",
			query: `DELETE myschema.mytable, myschema2.mytable2 
			FROM myschema.mytable JOIN myschema2.mytable2 ON myschema.mytable.id = myschema2.mytable2.id`,
			expected: []*sql.ResourceAccess{
				{
					AccessType: "delete",
					Schema:     "myschema",
					Table:      "mytable",
				},
				{
					AccessType: "delete",
					Schema:     "myschema2",
					Table:      "mytable2",
				},
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

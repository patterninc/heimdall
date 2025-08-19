package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/sql"
)

func TestParseSQLCreate(t *testing.T) {
	type testCase struct {
		name     string
		query    string
		expected []*sql.ResourceAccess
	}
	tests := []testCase{
		{
			name:  "Create table with schema",
			query: "CREATE TABLE myschema.mytable (id INT, name VARCHAR(100))",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "create",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Create table with quoted schema and table",
			query: "CREATE TABLE `myschema`.`mytable` (id INT, name VARCHAR(100))",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "create",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Create table with IF NOT EXISTS",
			query: "CREATE TABLE IF NOT EXISTS myschema.mytable (id INT)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "create",
					Schema:     "myschema",
					Table:      "mytable",
				},
			},
		},
		{
			name:  "Create table with complex column types",
			query: "CREATE TABLE public.mytable (id SERIAL PRIMARY KEY, data JSONB)",
			expected: []*sql.ResourceAccess{
				{
					AccessType: "create",
					Schema:     "public",
					Table:      "mytable",
				},
			},
		},
		{
			name: "Create view from csv files",
			query: `CREATE OR REPLACE VIEW myschema.myview
					USING CSV
					OPTIONS (
  						path 's3://heimdall/*.csv',
  						header 'true'
					);`,
			expected: []*sql.ResourceAccess{
				{
					AccessType: "create",
					Schema:     "myschema",
					Table:      "myview",
				},
			},
		},
		{
			name:  "Create table using iceberg",
			query: `CREATE TABLE IF NOT EXISTS myschema.mytable (id int) using iceberg`,
			expected: []*sql.ResourceAccess{
				{
					AccessType: "create",
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

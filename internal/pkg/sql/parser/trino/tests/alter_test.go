package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser"
	"github.com/patterninc/heimdall/internal/pkg/sql/parser/trino"
)

func TestParseSQLAlter(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []parser.Access
	}{
		{
			name:  "Alter table add column",
			query: "ALTER TABLE public.users ADD COLUMN email VARCHAR;",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "users",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.ALTER,
				},
			},
		},
		{
			name:  "Alter table with catalog",
			query: "ALTER TABLE test_catalog.public.inventory DROP COLUMN obsolete_flag;",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "inventory",
					Schema:  "public",
					Catalog: "test_catalog",
					Act:     parser.ALTER,
				},
			},
		},
		{
			name:  "Alter table rename table",
			query: "ALTER TABLE schema_name.old_table_name RENAME TO schema_name.new_table_name;",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "old_table_name",
					Schema:  "schema_name",
					Catalog: defaultCatalog,
					Act:     parser.ALTER,
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

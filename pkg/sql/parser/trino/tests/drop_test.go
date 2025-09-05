package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

func TestParseSQLDrop(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []parser.Access
	}{
		{
			name:  "simple Drop TABLE without catalog",
			query: "DROP TABLE public.sales",
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "public",
					Catalog: defaultCatalog,
					Act:     parser.DROP,
				},
			},
		},
		{
			name:  "Drop table with catalog",
			query: `DROP TABLE  IF EXISTS test_catalog.public.sales;`,
			expected: []parser.Access{
				&parser.TableAccess{
					Table:   "sales",
					Schema:  "public",
					Catalog: "test_catalog",
					Act:     parser.DROP,
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

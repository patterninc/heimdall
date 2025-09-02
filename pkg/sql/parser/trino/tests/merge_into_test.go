package sql

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

func TestParseSQLMergeInto(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*parser.TableAccess
	}{
		{
			name: "Simple MERGE INTO with catalog and schema",
			query: `MERGE INTO analytics.public.customers AS t
					USING analytics.staging.new_customers AS s
					ON t.customer_id = s.customer_id
					WHEN MATCHED THEN
						UPDATE SET
							name = s.name,
							email = s.email,
							updated_at = CURRENT_TIMESTAMP
					WHEN NOT MATCHED THEN
						INSERT (customer_id, name, email, updated_at)
						VALUES (s.customer_id, s.name, s.email, s.updated_at);`,
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "public",
					Name:    "customers",
					Catalog: "analytics",
				},
				{
					Access:  parser.SELECT,
					Schema:  "staging",
					Name:    "new_customers",
					Catalog: "analytics",
				},
			},
		},
		{
			name: "MERGE INTO with default catalog",
			query: `MERGE INTO public.products AS t
					USING public.incoming_products AS s
					ON t.sku = s.sku
					WHEN NOT MATCHED THEN
						INSERT (sku, name, category)
						VALUES (s.sku, s.name, s.category);`,
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Schema:  "public",
					Name:    "products",
					Catalog: defaultCatalog,
				},
				{
					Access:  parser.SELECT,
					Schema:  "public",
					Name:    "incoming_products",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name: "MERGE INTO with CTE",
			query: `
			MERGE INTO catalog_name.public.table_name target
				USING (
				WITH table_name AS (
					SELECT
						id,
						ROW_NUMBER() OVER (PARTITION BY id ORDER BY updated_at DESC) AS rn
					FROM catalog_name2.public2.table_name2
				)
					SELECT
						id,
					FROM table_name
					WHERE rn = 1
				) source
				ON target.id = source.id
				WHEN MATCHED THEN
				UPDATE SET
					name = source.name,
					email = source.email,
					is_active = source.is_active,
					updated_at = source.updated_at
				WHEN NOT MATCHED THEN
				INSERT (id )
				VALUES (source.id);`,
			expected: []*parser.TableAccess{
				{
					Access:  parser.INSERT,
					Catalog: "catalog_name",
					Schema:  "public",
					Name:    "table_name",
				},
				{
					Access:  parser.SELECT,
					Catalog: "catalog_name2",
					Schema:  "public2",
					Name:    "table_name2",
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

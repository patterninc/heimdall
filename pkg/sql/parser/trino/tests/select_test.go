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
		expected []parser.Access
	}{
		{
			name:  "Single table select",
			query: "SELECT * FROM schema1.table1",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Single table select with alias",
			query: "SELECT * FROM schema1.table1 as t1",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Multiple tables",
			query: "SELECT * FROM schema1.table1, testCatalog.schema2.table2",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: "testCatalog",
				},
			},
		},
		{
			name:  "Join tables",
			query: "SELECT * FROM ( SELECT * FROM schema1.table1, schema2.table2 )as b",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Join tables",
			query: "SELECT * FROM ( SELECT * FROM test_catalog.schema1.table1, schema2.table2 )as b",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: "test_catalog",
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Join with ON clause",
			query: "SELECT a.col1, b.col2 FROM schema1.table1 a JOIN schema2.table2 b ON a.id = b.id",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Subquery in FROM clause",
			query: "SELECT * FROM schema1.table1 WHERE id IN (SELECT id FROM schema2.table2)",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Subquery with alias",
			query: "SELECT * FROM (SELECT * FROM schema.table WHERE id IN (SELECT id FROM schema2.table2)) as sub",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema",
					Table:   "table",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "Complex query with multiple joins and subqueries",
			query: "SELECT * FROM schema.table LEFT JOIN schema2.table2 ON schema.table.id = schema2.table2.id",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema",
					Table:   "table",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name:  "EXISTS subquery",
			query: "SELECT * FROM schema1.table1 WHERE EXISTS (SELECT 1 FROM schema2.table2 WHERE schema2.table2.id = schema1.table1.id)",
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
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
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema1",
					Table:   "table1",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema2",
					Table:   "table2",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema3",
					Table:   "table3",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "schema4",
					Table:   "table4",
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
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "sales",
					Table:   "orders",
					Catalog: defaultCatalog,
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "sales",
					Table:   "customers",
					Catalog: defaultCatalog,
				},
			},
		},
		{
			name: "SELECT with cross join in subquery",
			query: `SELECT
					dag_id,
					COUNT(*) AS task_count,
					COUNT_IF(CAST(json_extract(task, '$.t') AS varchar) = 'EQUAL') AS task_equal_count
				FROM (
					SELECT
						dag_id,
						task
					FROM
						airflow.public.serialized_dag
						CROSS JOIN UNNEST(
							MAP_VALUES(
								CAST(json_extract(data, '$.tasks') AS map(varchar, json))
							)
						) AS t(task)
				) AS task_unnested
				GROUP BY dag_id`,
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "public",
					Table:   "serialized_dag",
					Catalog: "airflow",
				},
			},
		},
		{
			name: "Complex query with multiples joins #1",
			query: `SELECT username AS username 
					FROM (SELECT j.job_id,
						j.job_name,
						j.job_version,
						j.job_error,
						j.username,
						from_unixtime(j.created_at) AS created_at,
						js.job_status_name AS job_status,
						if(j.is_sync, 'sync', 'async') AS job_type,
						cl.cluster_id,
						cm.command_id
					FROM heimdall.public.jobs j
					JOIN heimdall.public.job_statuses js ON js.job_status_id = j.job_status_id
					JOIN heimdall.public.clusters cl ON cl.system_cluster_id = j.job_cluster_id
					JOIN heimdall.public.commands cm ON cm.system_command_id = j.job_command_id
					WHERE from_unixtime(j.created_at) >= CURRENT_TIMESTAMP - INTERVAL '45' DAY
					) AS virtual_table GROUP BY username ORDER BY username ASC
					LIMIT 1000`,
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Catalog: "heimdall",
					Schema:  "public",
					Table:   "jobs",
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Catalog: "heimdall",
					Schema:  "public",
					Table:   "job_statuses",
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Catalog: "heimdall",
					Schema:  "public",
					Table:   "clusters",
				},
				&parser.TableAccess{
					Act:     parser.SELECT,
					Catalog: "heimdall",
					Schema:  "public",
					Table:   "commands",
				},
			},
		},
		{
			name: "Complex query with unnest",
			query: `SELECT
						*
						FROM
						airflow.public.serialized_dag,
						UNNEST(
							CAST(
							json_extract(CAST(data AS JSON), '$.dag.tasks') AS ARRAY(JSON)
							)
						) AS t (task)
						where dag_id='dag_id_from_query'  
						LIMIT 1001`,
			expected: []parser.Access{
				&parser.TableAccess{
					Act:     parser.SELECT,
					Schema:  "public",
					Table:   "serialized_dag",
					Catalog: "airflow",
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

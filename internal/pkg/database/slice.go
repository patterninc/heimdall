package database

import (
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
)

var (
	sliceMetrics = telemetry.NewMethod("db_connection", "database_slice")
)

func GetSlice(db *Database, query string) (any, error) {

	// Track DB connection for slice query operation
	defer sliceMetrics.RecordLatency(time.Now(), "operation", "get_slice")
	sliceMetrics.CountRequest("operation", "get_slice")

	// open connection
	sess, err := db.NewSession(false)
	if err != nil {
		sliceMetrics.LogAndCountError(err, "operation", "get_slice")
		return nil, err
	}
	defer sess.Close()

	rows, err := sess.Query(query)
	if err != nil {
		sliceMetrics.LogAndCountError(err, "operation", "get_slice")
		return nil, err
	}
	defer rows.Close()

	result := make([]any, 0, 10)

	for rows.Next() {

		var item any
		if err := rows.Scan(&item); err != nil {
			sliceMetrics.LogAndCountError(err, "operation", "get_slice")
			return nil, err
		}

		result = append(result, item)

	}

	sliceMetrics.CountSuccess("operation", "get_slice")
	return result, nil

}

package database

import (
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
)

var (
	getSliceMethod = telemetry.NewMethod("db_connection", "get_slice")
)

func GetSlice(db *Database, query string) (any, error) {

	// Track DB connection for slice query operation
	defer getSliceMethod.RecordLatency(time.Now())
	getSliceMethod.CountRequest()

	// open connection
	sess, err := db.NewSession(false)
	if err != nil {
		getSliceMethod.LogAndCountError(err, "new_session")
		return nil, err
	}
	defer sess.Close()

	rows, err := sess.Query(query)
	if err != nil {
		getSliceMethod.LogAndCountError(err, "query")
		return nil, err
	}
	defer rows.Close()

	result := make([]any, 0, 10)

	for rows.Next() {

		var item any
		if err := rows.Scan(&item); err != nil {
			getSliceMethod.LogAndCountError(err, "scan")
			return nil, err
		}

		result = append(result, item)

	}

	getSliceMethod.CountSuccess()
	return result, nil

}

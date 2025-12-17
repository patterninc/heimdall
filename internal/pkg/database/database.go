package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/babourine/x/pkg/set"
	"github.com/hladush/go-telemetry/pkg/telemetry"
	_ "github.com/lib/pq"
)

const (
	dbDriverName = `postgres`
)

var (
	ctx              = context.Background()
	newSessionMethod = telemetry.NewMethod("db_connection", "new_session")
)

type Database struct {
	ConnectionString string `yaml:"connection_string,omitempty" json:"connection_string,omitempty"`
}

type Session struct {
	db        *sql.DB
	trx       *sql.Tx
	committed bool
}

func (d *Database) NewSession(withTransaction bool) (*Session, error) {

	// Track session creation metrics
	transactionLabel := "false"
	if withTransaction {
		transactionLabel = "true"
	}

	defer newSessionMethod.RecordLatency(time.Now(), "with_transaction", transactionLabel)
	newSessionMethod.CountRequest("with_transaction", transactionLabel)

	var err error

	s := &Session{}

	// Track database connection creation
	connectionStart := time.Now()
	newSessionMethod.CountRequest("with_transaction", transactionLabel)

	// open connection
	if s.db, err = sql.Open(dbDriverName, d.ConnectionString); err != nil {
		newSessionMethod.LogAndCountError(err, "with_transaction", transactionLabel)
		return nil, err
	}

	newSessionMethod.RecordLatency(connectionStart, "with_transaction", transactionLabel)
	newSessionMethod.CountSuccess("with_transaction", transactionLabel)

	// start transaction
	if withTransaction {
		if s.trx, err = s.db.BeginTx(ctx, nil); err != nil {
			s.db.Close() // Close the connection before returning error
			newSessionMethod.LogAndCountError(err, "with_transaction", transactionLabel)
			return nil, err
		}
	}

	newSessionMethod.CountSuccess("with_transaction", transactionLabel)
	return s, nil

}

func (s *Session) Close() error {

	// do we have an uncommitted transaction going? rollback!
	if s.trx != nil && !s.committed {
		s.trx.Rollback()
	}

	if s.db != nil {
		return s.db.Close()
	}

	return nil

}

func (s *Session) Commit() error {

	if s.trx != nil {
		if err := s.trx.Commit(); err != nil {
			return err
		}
		s.committed = true
	}

	return nil

}

func (s *Session) Rollback() error {

	if s.trx != nil {
		return s.trx.Rollback()
	}

	return nil

}

func (s *Session) InsertRow(query string, args ...any) (int64, error) {

	rowID := int64(0)

	// are we within transaction?
	if s.trx != nil {

		if err := s.trx.QueryRowContext(ctx, query, args...).Scan(&rowID); err != nil {
			return 0, err
		}

	} else {

		if err := s.db.QueryRow(query, args...).Scan(&rowID); err != nil {
			return 0, err
		}

	}

	return rowID, nil

}

func (s *Session) Exec(query string, args ...any) (int64, error) {

	// are we within transaction?
	if s.trx != nil {

		if _, err := s.trx.ExecContext(ctx, query, args...); err != nil {
			return 0, err
		}

	} else {

		result, err := s.db.Exec(query, args...)
		if err != nil {
			return 0, err
		}
		return result.RowsAffected()

	}

	return 0, nil

}

func (s *Session) QueryRow(query string, args ...any) (*sql.Row, error) {

	row := s.db.QueryRow(query, args...)

	if err := row.Err(); err != nil {
		return nil, err
	}

	return row, nil

}

func (s *Session) Query(query string, args ...any) (*sql.Rows, error) {

	return s.db.Query(query, args...)

}

func (s *Session) SelectSet(query string, args ...any) (*set.Set[string], error) {

	result := set.New([]string{})

	// let's add tags
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item string

	for rows.Next() {
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		result.Add(item)
	}

	return result, nil

}

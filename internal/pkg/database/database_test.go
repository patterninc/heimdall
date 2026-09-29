package database

import (
	"strings"
	"testing"
)

func TestSessionsSharePool(t *testing.T) {

	d := &Database{
		ConnectionString: "postgres://heimdall:heimdall@127.0.0.1:1/heimdall_pool_test?sslmode=disable",
		MaxOpenConns:     4,
		MaxIdleConns:     2,
	}

	s1, err := d.NewSession(false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s1.Close() })

	s2, err := d.NewSession(false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s2.Close() })

	if s1.db != s2.db {
		t.Fatal("sessions opened separate connection pools")
	}
	if s1.db.Stats().MaxOpenConnections != 4 {
		t.Fatalf("max open connections = %d, want 4", s1.db.Stats().MaxOpenConnections)
	}

	if err := s1.Close(); err != nil {
		t.Fatal(err)
	}

	// Closing a session must not close the shared pool.
	if err := s2.db.Ping(); err != nil && strings.Contains(err.Error(), "database is closed") {
		t.Fatal("closing a session closed the shared pool")
	}

	// A failed transaction must not close the pool either.
	if _, err := d.NewSession(true); err == nil {
		t.Fatal("expected begin transaction to fail without a database")
	}
	if err := s2.db.Ping(); err != nil && strings.Contains(err.Error(), "database is closed") {
		t.Fatal("failed NewSession closed the shared pool")
	}

}

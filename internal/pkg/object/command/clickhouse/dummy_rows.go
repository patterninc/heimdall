package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func dummyRows() driver.Rows {
	return &dummyRowsStruct{}
}

type dummyRowsStruct struct {
	// driver.Rows
}

// ColumnTypes implements driver.Rows.
func (d *dummyRowsStruct) ColumnTypes() []driver.ColumnType {
	return nil
}

// Scan implements driver.Rows.
func (d *dummyRowsStruct) Scan(dest ...any) error {
	return nil
}

// ScanStruct implements driver.Rows.
func (d *dummyRowsStruct) ScanStruct(dest any) error {
	return nil
}

// Totals implements driver.Rows.
func (d *dummyRowsStruct) Totals(dest ...any) error {
	return nil
}

func (d *dummyRowsStruct) Next() bool {
	return false
}

func (d *dummyRowsStruct) Columns() []string {
	return nil
}
func (d *dummyRowsStruct) Close() error {
	return nil
}
func (d *dummyRowsStruct) Err() error {
	return nil
}

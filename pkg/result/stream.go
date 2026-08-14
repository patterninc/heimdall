package result

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/patterninc/heimdall/pkg/result/column"
)

// RowWriter incrementally writes a Result-shaped {"columns":...,"data":...} object: WriteColumns once (optional), WriteRow per row, then Close.
type RowWriter struct {
	w            *bufio.Writer
	wroteColumns bool
	rowCount     int
	err          error
}

func NewRowWriter(w io.Writer) *RowWriter {
	rw := &RowWriter{w: bufio.NewWriter(w)}
	rw.err = rw.token(`{`)
	return rw
}

// WriteColumns must be called at most once, before any WriteRow; an empty/nil slice is a no-op, matching Result.Columns' omitempty.
func (rw *RowWriter) WriteColumns(columns []*column.Column) error {
	if rw.err != nil {
		return rw.err
	}
	if len(columns) == 0 {
		return nil
	}

	columnsJSON, err := json.Marshal(columns)
	if err != nil {
		return rw.fail(err)
	}
	if err := rw.token(`"columns":`); err != nil {
		return err
	}
	if _, err := rw.w.Write(columnsJSON); err != nil {
		return rw.fail(err)
	}
	rw.wroteColumns = true
	return nil
}

func (rw *RowWriter) WriteRow(row []any) error {
	if rw.err != nil {
		return rw.err
	}

	rowJSON, err := json.Marshal(row)
	if err != nil {
		return rw.fail(err)
	}

	var prefix string
	switch {
	case rw.rowCount == 0 && rw.wroteColumns:
		prefix = `,"data":[`
	case rw.rowCount == 0:
		prefix = `"data":[`
	default:
		prefix = `,`
	}
	if err := rw.token(prefix); err != nil {
		return err
	}
	if _, err := rw.w.Write(rowJSON); err != nil {
		return rw.fail(err)
	}
	rw.rowCount++
	return nil
}

// Close writes the closing bracket/brace and flushes; call exactly once, after all rows are written.
func (rw *RowWriter) Close() error {
	if rw.err != nil {
		return rw.err
	}
	if rw.rowCount > 0 {
		if err := rw.token(`]`); err != nil {
			return err
		}
	}
	if err := rw.token(`}`); err != nil {
		return err
	}
	return rw.w.Flush()
}

func (rw *RowWriter) fail(err error) error {
	rw.err = err
	return err
}

func (rw *RowWriter) token(s string) error {
	if _, err := rw.w.WriteString(s); err != nil {
		return rw.fail(err)
	}
	return nil
}

package trino

import "errors"

var (
	ErrEmptyClusterContext = errors.New("empty cluster context")
	ErrEmptyJobContext     = errors.New("empty job context")
)

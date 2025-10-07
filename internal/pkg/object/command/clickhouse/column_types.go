package clickhouse

import (
	"strings"
	"time"
)

// Map of base type -> handler
var chTypeHandlers = map[string]chScanHandler{
	"UInt8":       handleUInt8,
	"UInt16":      handleUInt16,
	"UInt32":      handleUInt32,
	"UInt64":      handleUInt64,
	"Int8":        handleInt8,
	"Int16":       handleInt16,
	"Int32":       handleInt32,
	"Int64":       handleInt64,
	"Float32":     handleFloat32,
	"Float64":     handleFloat64,
	"String":      handleString,
	"FixedString": handleString,
	"Date":        handleTime,
	"Date32":      handleTime,
	"DateTime":    handleTime,
	"DateTime64":  handleTime,
}

// Type handler signature
type chScanHandler func(nullable bool) (scanTarget any, reader func() any)

// Individual handlers

func handleUInt8(nullable bool) (any, func() any) {
	if nullable {
		var p *uint8
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v uint8
	return &v, func() any { return v }
}

func handleUInt16(nullable bool) (any, func() any) {
	if nullable {
		var p *uint16
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v uint16
	return &v, func() any { return v }
}

func handleUInt32(nullable bool) (any, func() any) {
	if nullable {
		var p *uint32
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v uint32
	return &v, func() any { return v }
}

func handleUInt64(nullable bool) (any, func() any) {
	if nullable {
		var p *uint64
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v uint64
	return &v, func() any { return v }
}

func handleInt8(nullable bool) (any, func() any) {
	if nullable {
		var p *int8
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v int8
	return &v, func() any { return v }
}

func handleInt16(nullable bool) (any, func() any) {
	if nullable {
		var p *int16
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v int16
	return &v, func() any { return v }
}

func handleInt32(nullable bool) (any, func() any) {
	if nullable {
		var p *int32
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v int32
	return &v, func() any { return v }
}

func handleInt64(nullable bool) (any, func() any) {
	if nullable {
		var p *int64
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v int64
	return &v, func() any { return v }
}

func handleFloat32(nullable bool) (any, func() any) {
	if nullable {
		var p *float32
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v float32
	return &v, func() any { return v }
}

func handleFloat64(nullable bool) (any, func() any) {
	if nullable {
		var p *float64
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v float64
	return &v, func() any { return v }
}

func handleString(nullable bool) (any, func() any) {
	if nullable {
		var p *string
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v string
	return &v, func() any { return v }
}

func handleTime(nullable bool) (any, func() any) {
	if nullable {
		var p *time.Time
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v time.Time
	return &v, func() any { return v }
}

func handleDefault(base string, nullable bool) (any, func() any) {
	// Treat Decimal, UUID, IPv4, IPv6 as string; fallback also string
	if nullable {
		var p *string
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v string
	return &v, func() any { return v }
}

func unwrapCHType(t string) (base string, nullable bool) {
	// Unwrap Nullable and LowCardinality; keep Array out for brevity
	s := t
	for {
		if strings.HasPrefix(s, "Nullable(") && strings.HasSuffix(s, ")") {
			nullable = true
			s = s[len("Nullable(") : len(s)-1]
			continue
		}
		if strings.HasPrefix(s, "LowCardinality(") && strings.HasSuffix(s, ")") {
			s = s[len("LowCardinality(") : len(s)-1]
			continue
		}
		break
	}
	// Decimal(N,S) normalize to "Decimal"
	if isDecimal(s) {
		return "Decimal", nullable
	}
	return s, nullable
}

func isDecimal(s string) bool {
	return strings.HasPrefix(s, "Decimal(") || strings.HasPrefix(s, "Decimal32(") ||
		strings.HasPrefix(s, "Decimal64(") || strings.HasPrefix(s, "Decimal128(") ||
		strings.HasPrefix(s, "Decimal256(")
}

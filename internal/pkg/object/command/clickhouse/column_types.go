package clickhouse

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Decimal struct {
	d *decimal.Decimal
}

func (d *Decimal) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case decimal.Decimal:
		if d.d == nil {
			d.d = new(decimal.Decimal)
		}
		*d.d = v
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Decimal", value)
	}
}

var chTypeToResultTypeName = map[string]string{
	"UInt8":       "int",
	"UInt16":      "int",
	"UInt32":      "int",
	"UInt64":      "long",
	"Int8":        "int",
	"Int16":       "int",
	"Int32":       "int",
	"Int64":       "long",
	"Float32":     "float",
	"Float64":     "double",
	"Decimal":     "double",
	"String":      "string",
	"FixedString": "string",
	"Date":        "int32",
	"Date32":      "int32",
	"DateTime":    "long",
	"DateTime64":  "long",
	"Bool":        "boolean",
}

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
	"Date":        handleDate,
	"Date32":      handleDate,
	"DateTime":    handleDateTime,
	"DateTime64":  handleDateTime,
	"Decimal":     handleDecimal,
	"Bool":        handleBool,
	"Array":       handleArray,
	"Tuple":       handleTuple,
}

// Type handler signature
type chScanHandler func(nullable bool) (scanTarget any, reader func() any)

// unified nullable helper
func makeScanTarget[T any](nullable bool) (any, func() any) {
	if nullable {
		var p *T
		return &p, func() any {
			if p == nil {
				return nil
			}
			return *p
		}
	}
	var v T
	return &v, func() any { return v }
}

// Individual handlers (all via generic helper)
func handleUInt8(nullable bool) (any, func() any)   { return makeScanTarget[uint8](nullable) }
func handleUInt16(nullable bool) (any, func() any)  { return makeScanTarget[uint16](nullable) }
func handleUInt32(nullable bool) (any, func() any)  { return makeScanTarget[uint32](nullable) }
func handleUInt64(nullable bool) (any, func() any)  { return makeScanTarget[uint64](nullable) }
func handleInt8(nullable bool) (any, func() any)    { return makeScanTarget[int8](nullable) }
func handleInt16(nullable bool) (any, func() any)   { return makeScanTarget[int16](nullable) }
func handleInt32(nullable bool) (any, func() any)   { return makeScanTarget[int32](nullable) }
func handleInt64(nullable bool) (any, func() any)   { return makeScanTarget[int64](nullable) }
func handleFloat32(nullable bool) (any, func() any) { return makeScanTarget[float32](nullable) }
func handleFloat64(nullable bool) (any, func() any) { return makeScanTarget[float64](nullable) }
func handleString(nullable bool) (any, func() any)  { return makeScanTarget[string](nullable) }
func handleBool(nullable bool) (any, func() any)    { return makeScanTarget[bool](nullable) }

func handleDate(nullable bool) (any, func() any) {
	target, reader := makeScanTarget[time.Time](nullable)
	return target, func() any {
		v := reader()
		if v == nil {
			return nil
		}
		t := v.(time.Time)
		// Return date as "YYYY-MM-DD"
		return t.Format("2006-01-02")
	}
}

func handleDateTime(nullable bool) (any, func() any) {
	target, reader := makeScanTarget[time.Time](nullable)
	return target, func() any {
		v := reader()
		if v == nil {
			return nil
		}
		t := v.(time.Time)
		// Check if there are any subsecond components
		if t.Nanosecond() > 0 {
			// Format with microseconds
			return t.Format("2006-01-02T15:04:05.999999999")
		}
		// DateTime (second precision)
		return t.Format("2006-01-02T15:04:05")
	}
}

func handleDecimal(nullable bool) (any, func() any) {
	var v Decimal
	return &v, func() any {
		if v.d == nil {
			return nil
		}
		val, _ := v.d.Float64()
		return val
	}
}
func handleTuple(nullable bool) (any, func() any) {
	if nullable {
		var p *any
		return &p, func() any {
			if p == nil || *p == nil {
				return nil
			}
			return *p
		}
	}
	var v any
	return &v, func() any { return v }
}

func handleArray(nullable bool) (any, func() any) {
	if nullable {
		var p *any
		return &p, func() any {
			if p == nil || *p == nil {
				return nil
			}
			return *p
		}
	}
	var v any
	return &v, func() any { return v }
}

func handleDefault(nullable bool) (any, func() any) {
	// Treat Decimal, UUID, IPv4, IPv6 as string; fallback also string
	return makeScanTarget[string](nullable)
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

	if strings.HasPrefix(s, "DateTime64(") {
		return "DateTime64", nullable
	}
	if strings.HasPrefix(s, "DateTime(") {
		return "DateTime", nullable
	}
	if strings.HasPrefix(s, "Array(") {
		return "Array", nullable
	}
	if strings.HasPrefix(s, "Tuple(") {
		return "Tuple", nullable
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

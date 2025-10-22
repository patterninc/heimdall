package clickhouse

import (
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

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
	"Date":        handleTime,
	"Date32":      handleTime,
	"DateTime":    handleTime,
	"DateTime64":  handleTime,
	"Decimal":     handleDecimal,
	"Bool":        handleBool,
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
func handleTime(nullable bool) (any, func() any)    { return makeScanTarget[time.Time](nullable) }
func handleBool(nullable bool) (any, func() any)    { return makeScanTarget[bool](nullable) }
func handleDecimal(nullable bool) (any, func() any) {
	if nullable {
		var p decimal.NullDecimal
		return &p, func() any {
			if p.Valid {
				val, _ := p.Decimal.Float64()
				return val
			}
			return nil
		}
	}
	var v decimal.Decimal
	return &v, func() any {
		val, _ := v.Float64()
		return val
	}
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

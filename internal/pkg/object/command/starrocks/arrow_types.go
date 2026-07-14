package starrocks

import (
	"encoding/json"
	"fmt"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"

	"github.com/patterninc/heimdall/pkg/result/column"
)

var arrowTypeToResultTypeName = map[arrow.Type]string{
	arrow.INT8:         "int",
	arrow.INT16:        "int",
	arrow.INT32:        "int",
	arrow.INT64:        "long",
	arrow.UINT8:        "int",
	arrow.UINT16:       "int",
	arrow.UINT32:       "long",
	arrow.UINT64:       "long",
	arrow.FLOAT32:      "float",
	arrow.FLOAT64:      "double",
	arrow.DECIMAL128:   "double",
	arrow.DECIMAL256:   "double",
	arrow.STRING:       "string",
	arrow.LARGE_STRING: "string",
	arrow.BINARY:       "string",
	arrow.LARGE_BINARY: "string",
	arrow.BOOL:         "boolean",
	arrow.DATE32:       "string",
	arrow.DATE64:       "string",
	arrow.TIMESTAMP:    "string",
	arrow.TIME32:       "string",
	arrow.TIME64:       "string",
	arrow.LIST:         "string",
	arrow.LARGE_LIST:   "string",
	arrow.MAP:          "string",
	arrow.STRUCT:       "string",
}

func columnsFromSchema(schema *arrow.Schema) []*column.Column {

	fields := schema.Fields()
	columns := make([]*column.Column, len(fields))

	for i, f := range fields {
		typeName, ok := arrowTypeToResultTypeName[f.Type.ID()]
		if !ok {
			typeName = "string"
		}
		columns[i] = &column.Column{
			Name: f.Name,
			Type: column.Type(typeName),
		}
	}

	return columns

}

func recordToRows(rec arrow.Record) [][]any {

	numRows := int(rec.NumRows())
	numCols := int(rec.NumCols())

	rows := make([][]any, numRows)
	for r := 0; r < numRows; r++ {
		row := make([]any, numCols)
		for c := 0; c < numCols; c++ {
			row[c] = arrowValue(rec.Column(c), r)
		}
		rows[r] = row
	}

	return rows

}

func arrowValue(col arrow.Array, i int) any {

	if col.IsNull(i) {
		return nil
	}

	switch v := col.(type) {
	case *array.Int8:
		return v.Value(i)
	case *array.Int16:
		return v.Value(i)
	case *array.Int32:
		return v.Value(i)
	case *array.Int64:
		return v.Value(i)
	case *array.Uint8:
		return v.Value(i)
	case *array.Uint16:
		return v.Value(i)
	case *array.Uint32:
		return v.Value(i)
	case *array.Uint64:
		return v.Value(i)
	case *array.Float32:
		return v.Value(i)
	case *array.Float64:
		return v.Value(i)
	case *array.Decimal128:
		return v.Value(i).ToFloat64(int32(v.DataType().(*arrow.Decimal128Type).Scale))
	case *array.Decimal256:
		return v.Value(i).ToFloat64(int32(v.DataType().(*arrow.Decimal256Type).Scale))
	case *array.String:
		return v.Value(i)
	case *array.LargeString:
		return v.Value(i)
	case *array.Binary:
		return string(v.Value(i))
	case *array.LargeBinary:
		return string(v.Value(i))
	case *array.Boolean:
		return v.Value(i)
	case *array.Date32:
		return v.Value(i).ToTime().Format("2006-01-02")
	case *array.Date64:
		return v.Value(i).ToTime().Format("2006-01-02")
	case *array.Timestamp:
		unit := v.DataType().(*arrow.TimestampType).Unit
		return v.Value(i).ToTime(unit).Format("2006-01-02T15:04:05.999999999")
	case *array.Time32:
		return v.Value(i).ToTime(v.DataType().(*arrow.Time32Type).Unit).Format("15:04:05")
	case *array.Time64:
		return v.Value(i).ToTime(v.DataType().(*arrow.Time64Type).Unit).Format("15:04:05.999999999")
	case *array.List, *array.LargeList, *array.Map, *array.Struct:
		b, err := json.Marshal(col.GetOneForMarshal(i))
		if err != nil {
			return fmt.Sprintf("%v", col.GetOneForMarshal(i))
		}
		return string(b)
	default:
		return fmt.Sprintf("%v", col.GetOneForMarshal(i))
	}

}

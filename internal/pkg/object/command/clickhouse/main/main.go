package main

import (
	"context"
	"fmt"
	"log"

	cvc "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/patterninc/heimdall/internal/pkg/object/command/clickhouse"
)

func main() {
	// Create connection to ClickHouse
	conn, err := cvc.Open(&cvc.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: cvc.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Debug: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to ClickHouse:", err)
	}
	defer conn.Close()

	// Test connection
	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal("Failed to ping ClickHouse:", err)
	}

	// Execute SELECT query
	ctx := context.Background()

	// if err:= conn.Exec(ctx, "DROP TABLE IF EXISTS all_types_demo"); err != nil {
	// 	log.Fatal("Failed to drop table:", err)
	// }
	// Create table (raw string, Exec returns only error)
	// if err := conn.Exec(ctx, `CREATE TABLE all_types_demo (
    // int8      Int8,
    // int16     Int16,
    // int32     Int32,
    // int64     Int64,
    // uint8     UInt8,
    // uint16    UInt16,
    // uint32    UInt32,
    // uint64    UInt64,
    // float32   Float32,
    // float64   Float64,
    // decimal32 Decimal32(3),
    // decimal64 Decimal64(6),
    // decimal128 Decimal128(10),
    // date_col Date,
    // date32_col Date32,
    // datetime_col DateTime,
    // string_col String,
    // nullable_col Nullable(Int32),
	// ) ENGINE = MergeTree()
	// ORDER BY tuple();`); err != nil {
	//     log.Fatal("Failed to create table:", err)
	// }
	
	// Create table (raw string, Exec returns only error)
	// if err := conn.Exec(ctx, `INSERT INTO all_types_demo VALUES (
	//     -8,             -- int8
	//     -16,            -- int16
	//     -32,            -- int32
	//     -64,            -- int64
	//     8,              -- uint8
	//     16,             -- uint16
	//     32,             -- uint32
	//     64,             -- uint64
	//     1.23,           -- float32
	//     4.56,           -- float64
	//     123.456,        -- decimal32
	//     7890.123456,    -- decimal64
	//     123456789.1234567890, -- decimal128
	//     '2024-06-17',   -- date_col
	//     '2024-06-17',   -- date32_col
	//     '2024-06-17 10:20:30', -- datetime_col
	//     'hello',        -- string_col
	//     NULL           -- nullable_col
	// );`); err != nil {
	//     log.Fatal("Failed to create table:", err)
	// }





//     nullable_col Nullable(Int32),

// ) ENGINE = MergeTree()
// ORDER BY tuple();
	rows, err := conn.Query(ctx, "SELECT * FROM default.all_types_demo")

	result, err := clickhouse.CollectResults(rows)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result: %+v\n", result)
}

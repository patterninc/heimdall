package result

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/linkedin/goavro"

	"github.com/patterninc/heimdall/pkg/result/column"
)

const (
	dataInitialSize   = 10000
	messageColumn     = `message`
	messageColumnType = `string`
	avroExtension     = `.avro`
)

var (
	ctx  = context.Background()
	rxS3 = regexp.MustCompile(`^s3://([^/]+)/(.*)$`)
)

type Result struct {
	Columns []*column.Column `yaml:"columns,omitempty" json:"columns,omitempty"`
	Data    [][]any          `yaml:"data,omitempty" json:"data,omitempty"`
}

type avroFields struct {
	Fields []*column.Column `yaml:"fields,omitempty" json:"fields,omitempty"`
}

func FromRows(rows *sql.Rows) (*Result, error) {

	defer rows.Close()

	// let's get columns metadata
	columnsTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	columns := make([]*column.Column, 0, len(columnsTypes))
	for _, c := range columnsTypes {
		columns = append(columns, &column.Column{
			Name: c.Name(),
			Type: column.Type(c.DatabaseTypeName()),
		})
	}

	// let's get columns count
	columnsCount := len(columns)

	// let's pull data
	data := make([][]any, 0, dataInitialSize)

	for rows.Next() {
		row := make([]any, columnsCount)
		for i := range row {
			row[i] = new(interface{})
		}
		if err := rows.Scan(row...); err != nil {
			return nil, err
		}
		data = append(data, row)
	}

	return &Result{
		Columns: columns,
		Data:    data,
	}, nil

}

func FromMessage(message string) (*Result, error) {

	return &Result{
		Columns: []*column.Column{{
			Name: messageColumn,
			Type: messageColumnType,
		}},
		Data: [][]any{{message}},
	}, nil

}

func FromJson(json map[string]any) (*Result, error) {
	return &Result{
		Columns: []*column.Column{{
			Name: messageColumn,
			Type: messageColumnType,
		}},
		Data: [][]any{{json}},
	}, nil
}

func FromAvro(uri string) (*Result, error) {

	var nothing *string

	// upload file
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Create an S3 client
	svc := s3.NewFromConfig(awsConfig)

	// get bucket name and prefix
	s3Parts := rxS3.FindAllStringSubmatch(uri, -1)
	if len(s3Parts) == 0 || len(s3Parts[0]) < 3 {
		return nil, fmt.Errorf("unexpected queries key: %v", s3Parts)
	}

	listObjectsOutput, err := svc.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: &s3Parts[0][1],
		Prefix: &s3Parts[0][2],
	})
	if err != nil {
		return nil, err
	}

	r := &Result{
		Data: make([][]any, 0, 10),
	}

	for _, c := range listObjectsOutput.Contents {
		if key := *c.Key; strings.HasSuffix(key, avroExtension) {
			data, err := getS3Object(svc, &s3Parts[0][1], c.Key)
			if err != nil {
				return nil, err
			}
			ocf, err := goavro.NewOCFReader(bytes.NewReader(data))
			if err != nil {
				return nil, err
			}
			// if we do not have columns metadata, let's pull it....
			if len(r.Columns) == 0 {
				fields := &avroFields{}
				if err := json.Unmarshal([]byte(ocf.Codec().Schema()), fields); err != nil {
					return nil, err
				}
				r.Columns = fields.Fields
			}
			// now let's process data
			for ocf.Scan() {
				recordObject, err := ocf.Read()
				if err != nil {
					return nil, err
				}
				row := make([]any, 0, len(r.Columns))
				record, ok := recordObject.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("failed to parse record. unexpected type: %T", recordObject)
				}
				for _, c := range r.Columns {
					if m, ok := record[c.Name].(map[string]any); ok && c.Type.IsPrimitive() {
						if v, found := m[string(c.Type)]; found {
							row = append(row, v)
						} else {
							row = append(row, nothing)
						}
					} else {
						row = append(row, record[c.Name])
					}
				}
				r.Data = append(r.Data, row)
			}
		}
	}

	return r, nil

}

func FromDynamo(items []map[string]types.AttributeValue) (*Result, error) {

	// nothing in, nothing out...
	if len(items) == 0 {
		return nil, nil
	}

	// let's get our columns
	columns := make([]*column.Column, 0, len(items[0]))
	for name, value := range items[0] {
		ct := column.Type(``)
		switch vt := value.(type) {
		case *types.AttributeValueMemberS:
			ct = `string`
		case *types.AttributeValueMemberN:
			ct = `double`
		case *types.AttributeValueMemberBOOL:
			ct = `boolean`
			/* FIXME: add support for lists and maps
			case *types.AttributeValueMemberL:
				fmt.Printf("Attribute: %s, Type: list, Value: %+v\n", name, vt.Value)
			case *types.AttributeValueMemberM:
				fmt.Printf("Attribute: %s, Type: map, Value: %+v\n", name, vt.Value)
			*/
		default:
			return nil, fmt.Errorf("unsupported data type in column %s: %+v", name, vt)
		}
		columns = append(columns, &column.Column{
			Name: name,
			Type: ct,
		})
	}

	r := &Result{
		Columns: columns,
		Data:    make([][]any, 0, len(items)),
	}

	for _, item := range items {
		row := make([]any, len(columns))
		for i, c := range columns {
			if v, ok := item[c.Name]; ok {
				switch vt := v.(type) {
				case *types.AttributeValueMemberS:
					row[i] = vt.Value
				case *types.AttributeValueMemberN:
					floatValue, err := strconv.ParseFloat(vt.Value, 64)
					if err != nil {
						return nil, err
					}
					row[i] = floatValue
				case *types.AttributeValueMemberBOOL:
					row[i] = vt.Value
				default:
					return nil, fmt.Errorf("unsupported data type in column %s: %+v", c.Name, vt)
				}
			}
		}
		r.Data = append(r.Data, row)
	}

	return r, nil

}

func getS3Object(svc *s3.Client, bucket, key *string) ([]byte, error) {

	resp, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the file content
	return io.ReadAll(resp.Body)

}

package factory

import (
	"fmt"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser"
	"github.com/patterninc/heimdall/internal/pkg/sql/parser/trino"
)

type parserType string

const (
	parserTypeTrino parserType = "trino"
)

func CreateParserByType(typ string, defaultCatalog string) (parser.AccessReceiver, error) {
	switch typ {
	case string(parserTypeTrino):
		return trino.NewTrinoAccessReceiver(defaultCatalog), nil
	default:
		return nil, fmt.Errorf("unknown parser type: %s", typ)
	}
}
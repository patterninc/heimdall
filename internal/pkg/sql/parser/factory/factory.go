package factory

import (
	"fmt"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser"
	"github.com/patterninc/heimdall/internal/pkg/sql/parser/trino"
)

const (
	parserTypeTrino = "trino"
)

func CreateParserByType(typ string, defaultCatalog string) (parser.AccessReceiver, error) {
	switch typ {
	case parserTypeTrino:
		return trino.NewTrinoAccessReceiver(defaultCatalog), nil
	default:
		return nil, fmt.Errorf("unknown parser type: %s", typ)
	}
}

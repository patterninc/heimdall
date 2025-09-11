package factory

import (
	"fmt"

	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

func CreateParserByType(typ string, defaultCatalog string) (parser.AccessReceiver, error) {
	switch typ {
	case "trino":
		return trino.NewTrinoAccessReceiver(defaultCatalog), nil
	default:
		return nil, fmt.Errorf("unknown parser type: %s", typ)
	}
}

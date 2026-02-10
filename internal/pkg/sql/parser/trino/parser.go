package trino

import (
	"time"

	"github.com/antlr4-go/antlr/v4"
	"github.com/hladush/go-telemetry/pkg/telemetry"

	"github.com/patterninc/heimdall/internal/pkg/sql/parser"
	"github.com/patterninc/heimdall/internal/pkg/sql/parser/trino/grammar"
)

var (
	parseAccessMethod = telemetry.NewMethod("trino_parse_access", "heimdall")
)

type TrinoAccessReceiver struct {
	defaultCatalog string
}

func NewTrinoAccessReceiver(defaultCatalog string) *TrinoAccessReceiver {
	return &TrinoAccessReceiver{defaultCatalog: defaultCatalog}
}

func (t *TrinoAccessReceiver) ParseAccess(sql string) ([]parser.Access, error) {
	parseAccessMethod.CountRequest()
	defer parseAccessMethod.RecordLatency(time.Now())
	is := antlr.NewInputStream(sql)
	lexer := grammar.NewTrinoLexer(is)
	tokens := antlr.NewCommonTokenStream(lexer, 0)
	p := grammar.NewTrinoParser(tokens)
	tree := p.Statements()

	walker := antlr.NewParseTreeWalker()
	col := newListener(t.defaultCatalog, sql)
	walker.Walk(col, tree)

	return col.collected, nil
}

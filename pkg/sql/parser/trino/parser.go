package trino

import (
	"log"

	"github.com/antlr4-go/antlr/v4"
	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino/grammar"
)

type TrinoAccessReceiver struct {
	defaultCatalog string
}

func NewTrinoAccessReceiver(defaultCatalog string) *TrinoAccessReceiver {
	return &TrinoAccessReceiver{defaultCatalog: defaultCatalog}
}

func (t *TrinoAccessReceiver) ParseAccess(sql string) ([]parser.Access, error) {
	log.Printf("Parsing SQL query: %s", sql)
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


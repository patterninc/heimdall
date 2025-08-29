package trino

import (
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

func (t *TrinoAccessReceiver) ParseTableAccess(query string) ([]*parser.TableAccess, error) {
	is := antlr.NewInputStream(query)
	lexer := grammar.NewTrinoLexer(is)
	tokens := antlr.NewCommonTokenStream(lexer, 0)
	p := grammar.NewTrinoParser(tokens)
	tree := p.Statements()

	walker := antlr.NewParseTreeWalker()
	col := newListener(t.defaultCatalog, query)
	walker.Walk(col, tree)

	return col.collected, nil
}

package trino

import (
	"github.com/patterninc/heimdall/pkg/sql/parser"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino/grammar"
)

type trinoListener struct {
	*grammar.BaseTrinoParserListener
	defaultCatalog string
	collected      []parser.Access
	query          string
}

func newListener(defaultCatalog string, query string) *trinoListener {
	return &trinoListener{
		collected:      []parser.Access{},
		defaultCatalog: defaultCatalog,
		query:          query,
	}
}

func (l *trinoListener) addTableAccess(name grammar.IQualifiedNameContext, access parser.Action) {
	parts := []string{}
	// Trino requires tables to be specified as schema.name or catalog.schema.name.
	// If only a single identifier is present, it's likely an alias or an invalid query.
	if len(name.AllIdentifier()) < 2 {
		// Ignore single identifiers as they are not valid table references.
		return
	}
	if len(name.AllIdentifier()) == 2 {
		parts = append(parts, l.defaultCatalog)
	}
	for _, id := range name.AllIdentifier() {
		parts = append(parts, id.GetText())
	}

	l.collected = append(l.collected, &parser.TableAccess{
		Table:   parts[2],
		Schema:  parts[1],
		Catalog: parts[0],
		Act:     access,
	})
}

func (l *trinoListener) EnterInsertInto(ctx *grammar.InsertIntoContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.INSERT)

}

func (l *trinoListener) EnterDelete(ctx *grammar.DeleteContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.DELETE)
}

func (l *trinoListener) EnterTableName(ctx *grammar.TableNameContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.SELECT)
}

func (l *trinoListener) EnterUpdate(ctx *grammar.UpdateContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.UPDATE)
}

func (l *trinoListener) EnterMerge(ctx *grammar.MergeContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.INSERT)
}

func (l *trinoListener) EnterCreateTableAsSelect(ctx *grammar.CreateTableAsSelectContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.CREATE)
}

func (l *trinoListener) EnterCreateTable(ctx *grammar.CreateTableContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.CREATE)
}

func (l *trinoListener) EnterDropTable(ctx *grammar.DropTableContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.DROP)
}

// EnterAddColumn is called when production addColumn is entered.
func (l *trinoListener) EnterAddColumn(ctx *grammar.AddColumnContext) {
	l.addTableAccess(ctx.QualifiedName(), parser.ALTER)
}

// EnterDropColumn is called when production dropColumn is entered.
func (l *trinoListener) EnterDropColumn(ctx *grammar.DropColumnContext) {
	l.addTableAccess(ctx.QualifiedName(0), parser.ALTER)
}

func (l *trinoListener) EnterRenameTable(ctx *grammar.RenameTableContext) {
	l.addTableAccess(ctx.QualifiedName(0), parser.ALTER)
}

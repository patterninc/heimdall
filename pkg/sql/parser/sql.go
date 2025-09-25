package parser

type AccessKind int

const (
	TableAccessKind AccessKind = iota
	SchemaAccessKind
	CatalogAccessKind
)

type Action int

const (
	SELECT Action = iota
	INSERT
	UPDATE
	DELETE
	CREATE
	DROP
	ALTER
)

type Access interface {
	Kind() AccessKind
	Action() Action
}

type TableAccess struct {
	Catalog string
	Schema  string
	Table   string
	Act     Action
}

func (t *TableAccess) Kind() AccessKind { return TableAccessKind }
func (t *TableAccess) Action() Action   { return t.Act }

type AccessReceiver interface {
	ParseAccess(sql string) ([]Access, error)
}

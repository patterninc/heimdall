package parser

import "fmt"

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
	CREATE
	DROP
	DELETE
	USE
	ALTER
	GRANT
	REVOKE
	SHOW
	IMPERSONATE
	EXECUTE
	UPDATE
	READ_SYSTEM_INFORMATION
	WRITE_SYSTEM_INFORMATION
)

type Access interface {
	Kind() AccessKind
	Action() Action
	QualifiedName() string
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


func (t *TableAccess) QualifiedName() string {
	return fmt.Sprintf("%s.%s.%s", t.Catalog, t.Schema, t.Table)
}
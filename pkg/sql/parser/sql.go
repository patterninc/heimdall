package parser

type AccessType int

const (
	SELECT AccessType = iota
	INSERT
	UPDATE
	DELETE
	CREATE
)

type TableAccess struct {
	Name    string
	Schema  string
	Catalog string
	Access  AccessType
}

type AccessReceiver interface {
	ParseTableAccess(sql string) ([]*TableAccess, error)
}

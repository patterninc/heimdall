package rbac

type RBAC interface {
	Init() error
	HasAccess(user string, query string) (bool, error)
	GetName() string
}
package rbac

import (
	"context"
)

type RBAC interface {
	Init(ctx context.Context) error
	HasAccess(user string, query string) (bool, error)
}

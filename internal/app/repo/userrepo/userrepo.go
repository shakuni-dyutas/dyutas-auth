package userrepo

import (
	"context"

	"github.com/shakuni-dyutas/dyutas-auth/internal/domain/user"
)

// TODO how to handle transactions?
type UserRepo interface {
	CreateUserByGoogleId(ctx context.Context, googleIdSub string, user *user.User) (*user.User, error)
	GetUserByGoogleId(ctx context.Context, googleSub string) (*user.User, error)
}

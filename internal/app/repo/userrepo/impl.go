package userrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shakuni-dyutas/dyutas-auth/internal/db"
	"github.com/shakuni-dyutas/dyutas-auth/internal/domain/user"
)

func New(db *db.Conn) UserRepo {
	return &UserRepoImpl{db: db}
}

type UserRepoImpl struct {
	db *db.Conn
}

func (ur *UserRepoImpl) CreateUserByGoogleId(ctx context.Context, googleIdSub string, newUser *user.User) (*user.User, error) {
	newUserRecord, err := ur.db.Qrs.CreateUser(ctx, db.CreateUserParams{
		Code:            newUser.Code,
		GoogleID:        googleIdSub,
		Email:           newUser.Email,
		ProfileImageUrl: pgtype.Text{String: *newUser.ProfileImageURL, Valid: true},
		Username:        pgtype.Text{String: *newUser.Username, Valid: true},
		SignedUpAt:      pgtype.Timestamptz{Time: newUser.SignedUpAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	var usernameValue *string
	if newUserRecord.Username.Valid {
		usernameValue = &newUserRecord.Username.String
	}

	return user.From(newUserRecord.ID, newUserRecord.Code, newUserRecord.Email, nil, usernameValue, newUserRecord.SignedUpAt.Time), nil
}

func (ur *UserRepoImpl) GetUserByGoogleId(ctx context.Context, googleSub string) (*user.User, error) {
	userRecord, err := ur.db.Qrs.GetUserByGoogleID(ctx, googleSub)
	if err == pgx.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var usernameColumn *string
	if userRecord.Username.Valid {
		usernameColumn = &userRecord.Username.String
	}

	userEntity := user.From(userRecord.ID, userRecord.Code, userRecord.Email, nil, usernameColumn, userRecord.SignedUpAt.Time)

	return userEntity, nil
}

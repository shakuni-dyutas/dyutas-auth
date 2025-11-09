package auth

import "time"

type UserRefreshToken struct {
	UserId    int64
	TokenHash string
	CreatedAt time.Time
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func (t *UserRefreshToken) IsExpired() bool {
	return t.ExpiresAt.Before(time.Now())
}

func (t *UserRefreshToken) IsRevoked() bool {
	return t.RevokedAt != nil
}

package authsvc

import (
	"context"
	"time"
)

type AuthService interface {
	// SignWithGoogle authenticates and identifies a user with Google Id.
	// If the user is not found, registers them.
	// In result, returns the newly signed tokens and the user info.
	SignWithGoogle(ctx context.Context, googleAuthCode string) (*SignResult, error)
	// RefreshAuthState validates the refresh token and returns newly signed tokens.
	// It revokes given refresh token and returns extended refresh token.
	RefreshAuth(ctx context.Context, refreshToken string) (*AuthResult, error)
	// SignOut revokes the refresh token.
	// Cookie removal must be handled by HTTP layer.
	SignOut(ctx context.Context, refreshToken string) error
}

const RefreshTokenTTL = time.Hour * 24 * 6
const AccessTokenTTL = time.Minute * 30

type SignResult struct {
	AuthResult
	SigningUserInfo AuthUserInfo
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
}

type AuthUserInfo struct {
	Code            string
	Email           string
	Username        *string
	ProfileImageURL *string
}

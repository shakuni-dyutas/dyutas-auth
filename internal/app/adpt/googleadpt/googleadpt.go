package googleadpt

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type GoogleAdpt interface {
	VerifyAuthCode(ctx context.Context, googleAuthCode string) (*GoogleUserInfo, error)
}

type GoogleUserInfo struct {
	Sub             string
	Email           string
	ProfileImageURL string
}

// --- external API DTOs below

// googleSignRequestDTO is id token verification API request DTO for
// 'https://oauth2.googleapis.com/token'
type GoogleSignAPIRequestDTO struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type GoogleSignAPIResponseDTO struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	IdToken      string `json:"id_token"`
}

type GoogleIdTokenPayload struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`

	jwt.RegisteredClaims
}

package authsvc

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shakuni-dyutas/dyutas-auth/internal/helper"
)

type RefreshTokenClaims struct {
	Tid string  `json:"tid"`
	Ucd string  `json:"ucd"`
	Exp float64 `json:"exp"`
}

func (svc *AuthServiceImpl) signNewRefreshTokenFor(ctx context.Context, usercode string) (token string, tokenHash string, err error) {
	randomCode := helper.GenerateRandomCode(10)

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"tid": randomCode,
		"ucd": usercode,
		"exp": time.Now().Add(RefreshTokenTTL).Unix(),
	})

	refreshTokenStr, err := refreshToken.SignedString(svc.appJwtKey)
	if err != nil {
		return "", "", err
	}

	tokenHash, err = helper.Hash(refreshTokenStr)
	if err != nil {
		return "", "", err
	}

	// TODO store it in the db

	return refreshTokenStr, tokenHash, nil
}

type AccessTokenClaims struct {
	Rhs string  `json:"rhs"`
	Ucd string  `json:"ucd"`
	Exp float64 `json:"exp"`
}

func (svc *AuthServiceImpl) signNewAccessTokenFor(ctx context.Context, usercode string, refreshTokenHash string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"rhs": refreshTokenHash,
		"ucd": usercode,
		"exp": time.Now().Add(AccessTokenTTL).Unix(),
	})

	accessTokenStr, err := accessToken.SignedString(svc.appJwtKey)
	if err != nil {
		return "", err
	}

	return accessTokenStr, nil
}

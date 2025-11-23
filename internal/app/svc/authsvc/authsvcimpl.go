package authsvc

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/adpt/googleadpt"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/repo/userrepo"
	"github.com/shakuni-dyutas/dyutas-auth/internal/helper"

	"github.com/shakuni-dyutas/dyutas-auth/internal/domain/user"
)

func New(
	googleAdpt googleadpt.GoogleAdpt,
	userRepo userrepo.UserRepo,
	appJwtKey string,
) AuthService {
	return &AuthServiceImpl{
		userRepo:   userRepo,
		googleAdpt: googleAdpt,
		appJwtKey:  []byte(appJwtKey),
	}
}

type AuthServiceImpl struct {
	userRepo   userrepo.UserRepo
	googleAdpt googleadpt.GoogleAdpt
	appJwtKey  []byte
}

func (svc *AuthServiceImpl) SignWithGoogle(ctx context.Context, googleAuthCode string) (*SignResult, error) {
	googleUserInfo, err := svc.googleAdpt.VerifyAuthCode(ctx, googleAuthCode)
	if err != nil {
		return nil, &app.AppError{
			Code:    app.UnauthenticatedError,
			Message: "failed to verify Google auth code",
			Err:     err,
		}
	}

	existingUser, err := svc.userRepo.GetUserByGoogleId(ctx, googleUserInfo.Sub)
	if err != nil {
		return nil, &app.AppError{
			Code:    app.InternalServerError,
			Message: "failed to get user by Google ID",
			Err:     err,
		}
	}

	if existingUser == nil {
		existingUser, err = svc.registerNewUserByGoogleId(ctx, googleUserInfo)
		if err != nil {
			return nil, &app.AppError{
				Code:    app.InternalServerError,
				Message: "failed to register new user by Google ID",
				Err:     err,
			}
		}
	}

	refreshToken, refreshTokenHash, err := svc.signNewRefreshTokenFor(ctx, existingUser.Code)
	if err != nil {
		return nil, &app.AppError{
			Code:    app.InternalServerError,
			Message: "failed to sign new refresh token",
			Err:     err,
		}
	}

	accessToken, err := svc.signNewAccessTokenFor(ctx, existingUser.Code, refreshTokenHash)
	if err != nil {
		return nil, &app.AppError{
			Code:    app.InternalServerError,
			Message: "failed to sign new access token",
			Err:     err,
		}
	}

	signResult := &SignResult{
		AuthResult: AuthResult{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		SigningUserInfo: AuthUserInfo{
			Code:            existingUser.Code,
			Email:           existingUser.Email,
			Username:        existingUser.Username,
			ProfileImageURL: existingUser.ProfileImageURL,
		},
	}

	return signResult, nil
}

const currentUserCodeLength = 10

func (svc *AuthServiceImpl) registerNewUserByGoogleId(ctx context.Context, googleUserInfo *googleadpt.GoogleUserInfo) (*user.User, error) {
	newUserCode := helper.GenerateRandomCode(currentUserCodeLength)

	tempNewUser := user.New(newUserCode, googleUserInfo.Email, nil, nil)

	newUser, err := svc.userRepo.CreateUserByGoogleId(ctx, googleUserInfo.Sub, tempNewUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (svc *AuthServiceImpl) RefreshAuth(ctx context.Context, refreshTokenJWT string) (*AuthResult, error) {
	_, err := helper.Hash(refreshTokenJWT)
	if err != nil {
		return nil, err
	}
	// TODO validate the token hash and revoke previous one

	refreshToken, err := jwt.Parse(refreshTokenJWT, func(token *jwt.Token) (interface{}, error) {
		return svc.appJwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	tokenMapClaims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid refresh token claims")
	}

	refreshTokenClaims := RefreshTokenClaims{
		Tid: tokenMapClaims["tid"].(string),
		Ucd: tokenMapClaims["ucd"].(string),
		Exp: tokenMapClaims["exp"].(float64),
	}

	newRefreshToken, newRefreshTokenHash, err := svc.signNewRefreshTokenFor(ctx, refreshTokenClaims.Ucd)
	if err != nil {
		return nil, err
	}

	newAccessToken, err := svc.signNewAccessTokenFor(ctx, refreshTokenClaims.Ucd, newRefreshTokenHash)
	if err != nil {
		return nil, err
	}

	refreshResult := &AuthResult{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return refreshResult, nil
}

func (svc *AuthServiceImpl) SignOut(ctx context.Context, refreshTokenJWT string) error {
	_, err := helper.Hash(refreshTokenJWT)
	if err != nil {
		return err
	}

	// TODO revoke the token

	return nil
}

func (svc *AuthServiceImpl) GetSelfInfo(ctx context.Context, accessToken string) (*AuthUserInfo, error) {
	_, err := helper.Hash(accessToken)
	if err != nil {
		return nil, err
	}

	accessTokenJWT, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return svc.appJwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	tokenMapClaims, ok := accessTokenJWT.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid access token claims")
	}

	accessTokenClaims := AccessTokenClaims{
		Rhs: tokenMapClaims["rhs"].(string),
		Ucd: tokenMapClaims["ucd"].(string),
		Exp: tokenMapClaims["exp"].(float64),
	}

	user, err := svc.userRepo.GetUserByCode(ctx, accessTokenClaims.Ucd)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &AuthUserInfo{
		Code:            user.Code,
		Email:           user.Email,
		Username:        user.Username,
		ProfileImageURL: user.ProfileImageURL,
	}, nil
}

package rest

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/svc/authsvc"
)

func InitAuthControllerWith(router gin.IRouter, cnf RestConfig, authSvc authsvc.AuthService) {
	if router == nil {
		panic("gin router isn't configured while initializing Auth controller")
	}
	if authSvc == nil {
		panic("auth service isn't configured while initializing Auth controller")
	}

	RegisterHandlersWithOptions(router, &AuthController{logger: cnf.Logger, authSvc: authSvc}, GinServerOptions{BaseURL: "/auth"})
}

type AuthController struct {
	logger *slog.Logger

	authSvc authsvc.AuthService
}

func (ac *AuthController) GoogleSignIn(ctx *gin.Context) {
	var reqBody GoogleSignInRequestDTO
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, RestError{
			ErrorCode: string(app.BadRequestError),
			Message:   "failed to bind request body",
		})
		return
	}

	signResult, err := ac.authSvc.SignWithGoogle(ctx, reqBody.Code)
	var appErr *app.AppError
	if errors.As(err, &appErr) {
		if appErr.Code == app.UnauthenticatedError {
			ctx.JSON(http.StatusUnauthorized, RestError{
				ErrorCode: string(appErr.Code),
				Message:   appErr.Message,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, RestError{
			ErrorCode: string(appErr.Code),
			Message:   "failed to sign with Google",
		})
		return
	}
	if err != nil {
		ac.logger.Error("GoogleSignIn: Failed to sign with Google", "error", err)

		ctx.JSON(http.StatusInternalServerError, RestError{
			ErrorCode: string(app.InternalServerError),
			Message:   "failed to sign with Google",
		})
		return
	}

	signingUserInfo := signResult.SigningUserInfo

	setRefreshTokenCookie(ctx, signResult.RefreshToken)

	resp := GoogleSignInResponseDTO{
		AccessToken: signResult.AccessToken,
		User: UserDTO{
			Email:           signingUserInfo.Email,
			Username:        signingUserInfo.Username,
			ProfileImageURL: nil,
			UserId:          signingUserInfo.Code,
		},
	}

	ctx.JSON(http.StatusOK, resp)
}

func (ac *AuthController) RefreshAuthentication(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(dyutasRefreshTokenCookieName)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, RestError{
			ErrorCode: string(app.UnauthorizedError),
			Message:   "failed to get refresh token from cookie",
		})
		return
	}

	refreshResult, err := ac.authSvc.RefreshAuth(ctx, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, RestError{
			ErrorCode: string(app.UnauthorizedError),
			Message:   "failed to refresh authentication",
		})
		return
	}

	setRefreshTokenCookie(ctx, refreshResult.RefreshToken)

	ctx.JSON(http.StatusCreated, RefreshAuthResponseDTO{
		AccessToken: refreshResult.AccessToken,
	})
}

func (ac *AuthController) SignOut(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(dyutasRefreshTokenCookieName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, RestError{
			ErrorCode: string(app.BadRequestError),
			Message:   "failed to get refresh token from cookie",
		})
		return
	}

	err = ac.authSvc.SignOut(ctx, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, RestError{
			ErrorCode: string(app.InternalServerError),
			Message:   "failed to sign out",
		})
		return
	}

	removeRefreshTokenCookie(ctx)
	ctx.JSON(http.StatusNoContent, nil)
}

func (ac *AuthController) GetSelfInfo(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	if authorizationHeader == "" {
		ctx.JSON(http.StatusBadRequest, RestError{
			ErrorCode: string(app.BadRequestError),
			Message:   "missing authorization header",
		})
		return
	}
	accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

	selfInfo, err := ac.authSvc.GetSelfInfo(ctx, accessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, RestError{
			ErrorCode: string(app.InternalServerError),
			Message:   "failed to get self info",
		})
		return
	}

	ctx.JSON(http.StatusOK, SelfDTO{
		UserId:          selfInfo.Code,
		Email:           selfInfo.Email,
		Username:        selfInfo.Username,
		ProfileImageURL: selfInfo.ProfileImageURL,
	})
}

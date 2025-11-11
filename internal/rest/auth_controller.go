package rest

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

// TODO cookie get/set, responding should be commonized

// TODO domain should be injected
const dyutasRefreshTokenCookieDomain = "local-api.dyutas.com"
const dyutasRefreshTokenCookiePath = "/auth/refresh" // TODO remove hardcoded path
const dyutasRefreshTokenCookieName = "rtk"
const dyutasRefreshTokenCookieTTL = authsvc.RefreshTokenTTL

func (ac *AuthController) GoogleSignIn(ctx *gin.Context) {
	var reqBody GoogleSignInRequestDTO
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to bind request body - " + err.Error()})
		return
	}

	signResult, err := ac.authSvc.SignWithGoogle(ctx, reqBody.Code)
	if err != nil {
		ac.logger.Error("GoogleSignIn: Failed to sign with Google", "error", err)

		ctx.JSON(400, gin.H{"error": "GoogleSignIn: Failed to sign with Google - " + err.Error()})
		return
	}

	signingUserInfo := signResult.SigningUserInfo

	ctx.SetCookie(dyutasRefreshTokenCookieName, signResult.RefreshToken, int(dyutasRefreshTokenCookieTTL.Seconds()), dyutasRefreshTokenCookiePath, dyutasRefreshTokenCookieDomain, true, true)

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
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	refreshResult, err := ac.authSvc.RefreshAuth(ctx, refreshToken)
	if err != nil {
		ac.logger.Error("RefreshAuthentication: Failed to refresh authentication", "error", err)

		ctx.JSON(400, gin.H{"error": "RefreshAuthentication: Failed to refresh authentication - " + err.Error()})
		return
	}

	ctx.SetCookie(dyutasRefreshTokenCookieName, refreshResult.RefreshToken, int(dyutasRefreshTokenCookieTTL.Seconds()), dyutasRefreshTokenCookiePath, dyutasRefreshTokenCookieDomain, true, true)

	ctx.JSON(http.StatusCreated, gin.H{"accessToken": refreshResult.AccessToken})
}

func (ac *AuthController) SignOut(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie(dyutasRefreshTokenCookieName)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	err = ac.authSvc.SignOut(ctx, refreshToken)
	if err != nil {
		ac.logger.Error("SignOut: Failed to sign out", "error", err)

		ctx.JSON(400, gin.H{"error": "SignOut: Failed to sign out - " + err.Error()})
		return
	}

	ctx.SetCookie(dyutasRefreshTokenCookieName, "", -1, dyutasRefreshTokenCookiePath, dyutasRefreshTokenCookieDomain, true, true)
	ctx.JSON(http.StatusNoContent, nil)
}

func (ac *AuthController) GetSelfInfo(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	accessToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

	selfInfo, err := ac.authSvc.GetSelfInfo(ctx, accessToken)
	if err != nil {
		ac.logger.Error("GetSelfInfo: Failed to get self info", "error", err)

		ctx.JSON(400, gin.H{"error": "GetSelfInfo: Failed to get self info - " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, SelfDTO{
		UserId:          selfInfo.Code,
		Email:           selfInfo.Email,
		Username:        selfInfo.Username,
		ProfileImageURL: selfInfo.ProfileImageURL,
	})
}

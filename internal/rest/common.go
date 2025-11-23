package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/svc/authsvc"
)

// TODO domain should be injected
const dyutasRefreshTokenCookieDomain = "local-api.dyutas.com"
const dyutasRefreshTokenCookiePath = "/auth/refresh" // TODO remove hardcoded path
const dyutasRefreshTokenCookieName = "rtk"
const dyutasRefreshTokenCookieTTL = authsvc.RefreshTokenTTL

func setRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie(dyutasRefreshTokenCookieName, refreshToken, int(dyutasRefreshTokenCookieTTL.Seconds()), dyutasRefreshTokenCookiePath, dyutasRefreshTokenCookieDomain, true, true)
}

func removeRefreshTokenCookie(ctx *gin.Context) {
	ctx.SetCookie(dyutasRefreshTokenCookieName, "", -1, dyutasRefreshTokenCookiePath, dyutasRefreshTokenCookieDomain, true, true)
}

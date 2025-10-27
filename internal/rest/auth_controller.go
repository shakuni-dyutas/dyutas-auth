package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shakuni-dyutas/dyutas-auth/internal/db"
)

func InitAuthControllerWith(db *db.Conn, router gin.IRouter) {
	if db == nil {
		panic("database isn't configured while initializing Auth controller")
	}
	if router == nil {
		panic("gin router isn't configured while initializing Auth controller")
	}

	RegisterHandlersWithOptions(router, &AuthController{db: db}, GinServerOptions{BaseURL: "/auth"})
}

type AuthController struct {
	db *db.Conn
}

func (c *AuthController) GoogleSignIn(ctx *gin.Context) {
	fmt.Println("GoogleSignIn")

}

func (c *AuthController) RefreshAuthentication(ctx *gin.Context) {

}

package rest

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shakuni-dyutas/dyutas-auth/internal/app/svc/authsvc"
)

type RestConfig struct {
	Port          string
	AllowOrigins  []string
	AllowMethods  []string
	AllowHeaders  []string
	ExposeHeaders []string
	Logger        *slog.Logger
}

func Run(cnf RestConfig, authSvc authsvc.AuthService) error {
	if cnf.Logger == nil {
		panic("logger isn't configured while initializing REST server")
	}
	if authSvc == nil {
		cnf.Logger.Error("auth service isn't configured while initializing REST server")
		return nil
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     cnf.AllowOrigins,
		AllowMethods:     cnf.AllowMethods,
		AllowHeaders:     cnf.AllowHeaders,
		ExposeHeaders:    cnf.ExposeHeaders,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/auth/v1/openapi", func(c *gin.Context) {
		swagger, err := GetSwagger()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, swagger)
	})

	InitAuthControllerWith(router, cnf, authSvc)

	servicePort := cnf.Port
	if !strings.HasPrefix(servicePort, ":") {
		servicePort = ":" + servicePort
	}

	err := router.Run(servicePort)
	if err != nil {
		return err
	}

	return nil
}

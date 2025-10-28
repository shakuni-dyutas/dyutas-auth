package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shakuni-dyutas/dyutas-auth/internal/db"
	"github.com/shakuni-dyutas/dyutas-auth/internal/rest"
)

func main() {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://local.dyutas.com:8010", "https://local-api.dyutas.com:8010"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/auth/v1/openapi", func(c *gin.Context) {
		swagger, err := rest.GetSwagger()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, swagger)
	})

	rest.InitAuthControllerWith(&db.Conn{}, router)

	router.Run(":8020")
}

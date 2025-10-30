package main

import (
	"os"
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

	dbConn, err := db.NewRDBConnectionPool(db.ConnectionConfig{
		User: os.Getenv("AUTH_RDB_USER"),
		Pw:   os.Getenv("AUTH_RDB_PASSWORD"),
		Host: os.Getenv("AUTH_RDB_HOST"),
		Port: os.Getenv("AUTH_RDB_PORT"),
		DB:   os.Getenv("AUTH_RDB_DB"),
	})
	if err != nil {
		panic(err)
	}

	rest.InitAuthControllerWith(dbConn, router)

	router.Run(":8020")
}

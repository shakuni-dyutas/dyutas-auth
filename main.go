package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var v1 = os.Getenv("AUTH_RDB_PORT")
	var v2 = os.Getenv("FE_DOCKERFILE_PATH")

	fmt.Println(v1, v2)

	router := gin.Default()
	router.GET("/auth/dyutas", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Shakuni!",
		})
	})
	router.Run(":8020")
}

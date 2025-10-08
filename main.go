package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/auth/dyutas", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Shakuni!",
		})
	})
	router.Run(":8020")
}

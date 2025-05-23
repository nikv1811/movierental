package main

import (
	// "fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	router.Run(":8080")
}

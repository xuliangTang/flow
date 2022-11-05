package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Div(v1 int, v2 int) int {
	if v2 == 0 {
		return 0
	}
	return v1 / v2
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}

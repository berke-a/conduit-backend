package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	userServiceProxy := newReverseProxy("http://user-service:8080")
	router.Any("/users/*proxyPath", func(c *gin.Context) {
		userServiceProxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Run() // Listen and serve on 0.0.0.0:8080
}

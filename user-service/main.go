package main

import (
	"user-service/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/users/:id", handler.GetUserByID)

	router.Run()
}

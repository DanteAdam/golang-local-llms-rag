package main

import (
	"golang-llms/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run("0.0.0.0:8080")
}

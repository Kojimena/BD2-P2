package main

import (
	"backend/configs"
	"backend/routes"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	// Reference: https://github.com/gin-contrib/cors/issues/29
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(CORS())

	// Neo4j
	configs.ConnectDB()

	routes.Routes(router)

	router.Run("localhost:8080")
}

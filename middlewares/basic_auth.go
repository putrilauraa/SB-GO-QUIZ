package middlewares

import (
	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	users := gin.Accounts{
		"postgres": "postgres12345",
	}

	return gin.BasicAuth(users)
}
package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func RoleAuth(c *gin.Context) {
	fmt.Println("I'm in role auth")
	// c.JSON(200, gin.H{
	// 	"message": "You're not allowed",
	// })
	// c.AbortWithStatus(http.StatusUnauthorized)
	c.Next()
}

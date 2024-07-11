package middlewares

import (
	"fmt"
	"go-api-server/initializers"
	"go-api-server/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("I'm in RequireAuth")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	fmt.Println("I'm in RequireAuth")
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

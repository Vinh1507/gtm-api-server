package controllers

import (
	"fmt"
	"go-api-server/initializers"
	"go-api-server/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {

	// get data from req body

	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}

func Login(c *gin.Context) {

	// get data from req body

	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var user models.User
	initializers.DB.Where("email = ?", body.Email).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // exp in one month
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	fmt.Println(tokenString, err)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func Validate(c *gin.Context) {
	user, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusBadGateway, gin.H{
			"message": "Not Authorized",
		})
	}
	c.JSON(200, gin.H{
		"message": "I'm validating",
		"user":    user,
	})
}

package main

import (
	"go-api-server/controllers"
	"go-api-server/initializers"
	"go-api-server/middlewares"
	"go-api-server/migrate"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.ConnectToEtcd()
	migrate.Migrate()
}
func main() {
	r := gin.Default()

	// Cấu hình CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Thay thế bằng các nguồn mà bạn muốn cho phép
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middlewares.RequireAuth, middlewares.RoleAuth, controllers.Validate)

	r.POST("/etcd-put", controllers.EtcdPut)
	r.GET("/etcd-get-prefix", controllers.EtcdGetByPrefix)

	r.POST("/api/v1/gtm/config-gtm", controllers.ConfigGtm)
	r.GET("/api/v1/gtm/history", controllers.GetDataCenterHistory)

	r.GET("/api/v1/gtm/resolvers/:id", controllers.GetResolverDetail)
	r.GET("/api/v1/gtm/resolvers", controllers.GetResolverList)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Route not found"})
	})

	port := os.Getenv("PORT")
	r.Run(port)
}

package main

import (
	"go-api-server/controllers"
	"go-api-server/initializers"
	"go-api-server/middlewares"
	"go-api-server/migrate"
	"os"

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

	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middlewares.RequireAuth, middlewares.RoleAuth, controllers.Validate)

	r.POST("/etcd-put", controllers.EtcdPut)
	r.GET("/etcd-get-prefix", controllers.EtcdGetByPrefix)

	port := os.Getenv("PORT")
	r.Run(port)
}

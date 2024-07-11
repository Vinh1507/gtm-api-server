package migrate

import (
	"go-api-server/initializers"
	"go-api-server/models"
)

// func init() {
// 	initializers.LoadEnvVariables()
// 	initializers.ConnectToDB()

// }

// func main() {
// 	initializers.DB.AutoMigrate(&models.Post{})
// 	initializers.DB.AutoMigrate(&models.User{})
// }

func Migrate() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.DB.AutoMigrate(&models.User{})
}

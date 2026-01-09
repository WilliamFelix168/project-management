package routes

import (
	"log"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// s kecil dipakai biar ga bisa diakses dari luar package, begitu sebaliknya
func Setup(app *fiber.App, uc *controllers.UserController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	app.Post("/v1/auth/register", uc.Register)
}

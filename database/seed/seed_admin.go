package seed

import (
	"log"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("admin123")

	admin := models.User{
		Name:     "Super Admin",
		Email:    "admin@example.com",
		Password: password,
		Role:     "admin",
	}

	if err := config.DB.FirstOrCreate(&admin, models.User{Email: admin.Email}).Error;err != nil{
		log.Println("Failed to seed admin user:", err);
	}else {
		log.Println("Admin user seeded successfully");
	}

}

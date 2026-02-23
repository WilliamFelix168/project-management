package main

import (
	"log"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/controllers"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/database/seed"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/repositories"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/routes"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()

	app := fiber.New()

	//User
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	//Board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	//List
	listRepo := repositories.NewListRepository()
	listPosRepo := repositories.NewListPositionRepository()
	listService := services.NewListService(listRepo, boardRepo, listPosRepo)
	listController := controllers.NewListController(listService)

	routes.Setup(app, userController, boardController, listController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port :", port)
	log.Fatal(app.Listen(":" + port))
}

package routes

import (
	"log"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/controllers"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	jwtware "github.com/gofiber/jwt/v3"
)

// s kecil dipakai biar ga bisa diakses dari luar package, begitu sebaliknya
func Setup(app *fiber.App,
	uc *controllers.UserController,
	bc *controllers.BoardController,
	lc *controllers.ListController) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	app.Post("/v1/auth/register", uc.Register)
	app.Post("/v1/auth/login", uc.Login)

	//JWT proctected routes
	api := app.Group("/api/v1", jwtware.New(jwtware.Config{
		SigningKey: []byte(config.AppConfig.JWTSecret),
		// ContextKey untuk setelah token terverifikasi, paybload tokennya disimpen di Context dengan key "user",
		ContextKey: "user",
		//ErrorHandler digunakan untuk menangani error yang terjadi saat verifikasi token JWT gagal
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.Unauthorized(c, "Error unauthorized", err.Error())
		},
	}))

	userGroup := api.Group("/users")
	userGroup.Get("/page", uc.GetUserPagination)
	userGroup.Get("/:id", uc.GetUser) // /api/v1/users/:id
	userGroup.Put("/:id", uc.UpdateUser)
	userGroup.Delete("/:id", uc.DeleteUser)

	boardGroup := api.Group("/boards")
	boardGroup.Post("/", bc.CreateBoard)
	boardGroup.Put("/:id", bc.UpdateBoard)
	boardGroup.Post("/:id/members", bc.AddBoardMembers)
	boardGroup.Delete("/:id/members", bc.RemoveBoardMembers)
	boardGroup.Get("/my", bc.GetMyBoardPaginate)

	//list
	listGroup := api.Group("/lists")
	listGroup.Post("/", lc.CreateList)
	// listGroup.Get("/board/:boardPublicId", lc.GetListsByBoardID)
	// listGroup.Get("/:id", lc.GetListByID)
	// listGroup.Get("/public/:publicId", lc.GetListByPublicID)
	// listGroup.Put("/:id", lc.UpdateList)
	// listGroup.Delete("/:id", lc.DeleteList)
	// listGroup.Put("/board/:boardPublicId/positions", lc.UpdateListPositions)

}

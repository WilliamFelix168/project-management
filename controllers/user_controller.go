package controllers

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/services"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/gofiber/fiber/v2"
)

// membuat struct UserController yang berisi dependency ke UserService
// membuat request handler untuk user
type UserController struct {
	service services.UserService
	//dependency yang dibutuhkan controller utk logika bisnis
}

// fungsi untuk mengembalikan objek UserController dengan dependency UserService
// tujuannya untuk menginisialisasi controller dengan service yang diberikan
func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	//handler untuk register user
	user := new(models.User)

	// untuk memparsing body request ke struct user
	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal memparsing data user", err.Error())
	}

	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	return utils.Success(ctx, "Register Success", user)
}
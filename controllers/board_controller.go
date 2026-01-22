package controllers

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/services"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController(s services.BoardService) *BoardController {
	return &BoardController{
		service: s,
	}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	board := new(models.Board)

	var userID uuid.UUID
	var err error
	//fungsi untuk mengambil data user dari context setelah token JWT terverifikasi
	user := ctx.Locals("user").(*jwt.Token)
	//fungsi untuk mengakses klaim (claims) dari token JWT
	claims = user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Error Read Request", err.Error())
	}

	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Failed to create board", err.Error())
	}

	return utils.Success(ctx, "Board Created Successfully", board)
}

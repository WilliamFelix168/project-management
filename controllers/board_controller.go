package controllers

import (
	"math"
	"strconv"

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
	claims := user.Claims.(jwt.MapClaims)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Error Read Request", err.Error())
	}

	userID, err = uuid.Parse(claims["pub_id"].(string))
	if err != nil {
		return utils.BadRequest(ctx, "Invalid User ID", err.Error())
	}

	board.OwnerPublicID = userID

	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Failed to create board", err.Error())
	}

	return utils.Success(ctx, "Board Created Successfully", board)
}

func (c *BoardController) UpdateBoard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Error Parsing Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID Not Valid", err.Error())
	}

	existingBoard, err := c.service.GetByPublicID(publicID)

	if err != nil {
		return utils.NotFound(ctx, "Board Not Found", err.Error())
	}

	board.InternalID = existingBoard.InternalID
	board.PublicID = existingBoard.PublicID
	board.OwnerID = existingBoard.OwnerID
	board.OwnerPublicID = existingBoard.OwnerPublicID
	board.CreatedAt = existingBoard.CreatedAt

	if err := c.service.Update(board); err != nil {
		return utils.BadRequest(ctx, "Failed to update board", err.Error())
	}

	return utils.Success(ctx, "Board Updated Successfully", board)

}

func (c *BoardController) AddBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Error Parsing Data", err.Error())
	}
	if err := c.service.AddMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed to add board members", err.Error())
	}
	return utils.Success(ctx, "Board Members Added Successfully", nil)
}

func (c *BoardController) RemoveBoardMembers(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	var userIDs []string
	if err := ctx.BodyParser(&userIDs); err != nil {
		return utils.BadRequest(ctx, "Error Parsing Data", err.Error())
	}
	if err := c.service.RemoveMembers(publicID, userIDs); err != nil {
		return utils.BadRequest(ctx, "Failed to remove board members", err.Error())
	}
	return utils.Success(ctx, "Board Members Removed Successfully", nil)
}

func (c *BoardController) GetMyBoardPaginate(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["pub_id"].(string)

	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	boards, total, err := c.service.GetAllByUserPaginate(userID, filter, sort, limit, offset)
	if err != nil {
		return utils.BadRequest(ctx, "Failed to get boards", err.Error())
	}

	meta := utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		Filter:     filter,
		Sorting:    sort,
	}
	return utils.SuccessPagination(ctx, "Boards Retrieved Successfully", boards, meta)

}

package controllers

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/services"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ListController struct {
	services services.ListService
}

func NewListController(services services.ListService) *ListController {
	return &ListController{services}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Error Read Request", err.Error())
	}

	if err := c.services.Create(list); err != nil {
		return utils.BadRequest(ctx, "Failed to create list", err.Error())
	}

	return utils.Success(ctx, "List created successfully", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Error Data Parsing", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	existingList, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.services.Update(list); err != nil {
		return utils.BadRequest(ctx, "Failed to update list", err.Error())
	}

	updatedList, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found after update", err.Error())
	}

	return utils.Success(ctx, "List updated successfully", updatedList)
}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")

	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "Invalid Board ID", err.Error())
	}

	listsWithOrder, err := c.services.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Lists not found for the board", err.Error())
	}
	return utils.Success(ctx, "Lists retrieved successfully", listsWithOrder)
}

func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	list, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	if err := c.services.Delete(uint(list.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Failed to delete list", err.Error())
	}

	return utils.Success(ctx, "List deleted successfully", publicID)
}

func (c *ListController) UpdateListPositions(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	var body struct {
		Positions []uuid.UUID
	}

	if err := ctx.BodyParser(&body); err != nil {
		return utils.BadRequest(ctx, "Error Data Parsing", err.Error())
	}

	if len(body.Positions) == 0 {
		return utils.BadRequest(ctx, "Positions is required", "positions cannot be empty")
	}

	existingList, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List not found", err.Error())
	}

	if err := c.services.UpdatePositions(existingList.BoardPublicID.String(), body.Positions); err != nil {
		return utils.BadRequest(ctx, "Failed to update list positions", err.Error())
	}

	return utils.Success(ctx, "List positions updated successfully", body.Positions)
}

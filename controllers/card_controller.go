package controllers

import (
	"time"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/services"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CardController struct {
	service services.CardService
}

func NewCardController(s services.CardService) *CardController {
	return &CardController{service: s}
}

func (c *CardController) CreateCard(ctx *fiber.Ctx) error {
	type CreateCardRequest struct {
		ListPublicID string    `json:"list_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		DueDate      time.Time `json:"due_date"`
		Position     int       `json:"position"`
	}

	var req CreateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Failed get Data", err.Error())
	}

	card := &models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Position:    req.Position,
	}

	if err := c.service.Create(card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Failed to create card", err.Error())
	}

	return utils.Success(ctx, "Card created successfully", card)

}

func (c *CardController) UpdateCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	type updateCardRequest struct {
		ListPublicID string     `json:"list_id"`
		Title        string     `json:"title"`
		Description  string     `json:"description"`
		DueDate      *time.Time `json:"due_date"`
		Position     int        `json:"position"`
	}

	var req updateCardRequest
	if err := ctx.BodyParser(&req); err != nil {
		return utils.BadRequest(ctx, "Failed to parse Data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	card := models.Card{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     time.Time{},
		Position:    req.Position,
		PublicID:    uuid.MustParse(publicID),
	}

	if err := c.service.Update(&card, req.ListPublicID); err != nil {
		return utils.InternalServerError(ctx, "Failed Update Data", err.Error())
	}

	return utils.Success(ctx, "Card updated successfully", card)
}

func (c *CardController) DeleteCard(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "Invalid ID", err.Error())
	}

	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to find card", err.Error())
	}

	if err := c.service.Delete(uint(card.InternalID)); err != nil {
		return utils.BadRequest(ctx, "Failed to delete card", err.Error())
	}

	return utils.Success(ctx, "Card deleted successfully", nil)
}

func (c *CardController) GetCardDetail(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	card, err := c.service.GetByPublicID(publicID)
	if err != nil {
		return utils.InternalServerError(ctx, "Failed to find card", err.Error())
	}

	if card == nil {
		return utils.NotFound(ctx, "Card not found", err.Error())
	}

	return utils.Success(ctx, "Card retrieved successfully", card)
}

func (c *CardController) GetCardOnList(ctx *fiber.Ctx) error {
	listPublicID := ctx.Params("list_id")

	if _, err := uuid.Parse(listPublicID); err != nil {
		return utils.BadRequest(ctx, "Invalid List ID", err.Error())
	}

	cards, err := c.service.GetByListID(listPublicID)
	if err != nil {
		return utils.NotFound(ctx, "Cards not found for the list", err.Error())
	}
	return utils.Success(ctx, "Cards retrieved successfully", cards)
}

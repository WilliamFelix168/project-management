package controllers

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/services"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/gofiber/fiber/v2"
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

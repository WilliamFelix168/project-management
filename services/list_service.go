package services

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/repositories"
	"github.com/google/uuid"
)

type listService struct {
	listRepo    repositories.ListRepository
	boardRepo   repositories.BoardRepository
	listPosRepo repositories.ListPositionRepository
}

type ListWithOrder struct {
	Positions []uuid.UUID
	Lists     []models.List
}

type ListService interface {
	GetByBoardID(boardPublicId string) (*ListWithOrder, error)
	GetByID(id uint) (*models.List, error)
	GetByPublicID(publicId string) (*models.List, error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdateListPosition(boardPublicId string, position []string) error
}

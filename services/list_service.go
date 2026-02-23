package services

import (
	"errors"
	"fmt"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models/types"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/repositories"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
	UpdatePositions(boardPublicId string, position []uuid.UUID) error
}

func NewListService(
	listRepo repositories.ListRepository,
	boardRepo repositories.BoardRepository,
	listPosRepo repositories.ListPositionRepository,
) ListService {
	return &listService{listRepo, boardRepo, listPosRepo}
}

func (s *listService) GetByBoardID(boardPublicId string) (*ListWithOrder, error) {
	//verifikasi board
	_, err := s.boardRepo.FindByPublicID(boardPublicId)
	if err != nil {
		return nil, errors.New("board not found")
	}

	position, err := s.listPosRepo.GetListOrder(boardPublicId)
	if err != nil {
		return nil, errors.New("Failed to get list order: " + err.Error())
	}

	lists, err := s.listRepo.FindByBoardID(boardPublicId)
	if err != nil {
		return nil, errors.New("Failed to get lists: " + err.Error())
	}

	//sorting by position
	orderedList := utils.SortListsByPosition(lists, position)

	return &ListWithOrder{
		Positions: position,
		Lists:     orderedList,
	}, nil
}

func (s *listService) GetByID(id uint) (*models.List, error) {
	return s.listRepo.FindByID(id)
}

func (s *listService) GetByPublicID(publicId string) (*models.List, error) {
	return s.listRepo.FindByPublicID(publicId)
}

func (s *listService) Create(list *models.List) error {
	//validasi board
	//trasaction
	//update position
	//commit trx
	board, err := s.boardRepo.FindByPublicID(list.BoardPublicID.String())
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return errors.New("board not found")
		}

		return fmt.Errorf("failed to get board: %w", err)

	}

	list.BoardInternalID = board.InternalID

	if list.PublicID == uuid.Nil {
		list.PublicID = uuid.New()
	}
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(list).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list: %w", err)
	}

	var position models.ListPosition
	res := tx.Where("board_internal_id = ?", board.InternalID).First(&position)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		//create new position
		position = models.ListPosition{
			PublicID:  uuid.New(),
			BoardID:   board.InternalID,
			ListOrder: types.UUIDArray{list.PublicID},
		}
		if err := tx.Create(&position).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create list position: %w", err)
		}
	} else if res.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create list position: %w", res.Error)
	} else {
		position.ListOrder = append(position.ListOrder, list.PublicID)

		if err := tx.Model(&position).Update("list_order", position.ListOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update list position: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *listService) Update(list *models.List) error {
	//validasi board
	return s.listRepo.Update(list)
}

func (s *listService) Delete(id uint) error {
	return s.listRepo.Delete(id)
}

func (s *listService) UpdatePositions(boardPublicId string, positions []uuid.UUID) error {
	//verifikasi board
	board, err := s.boardRepo.FindByPublicID(boardPublicId)
	if err != nil {
		return errors.New("board not found")
	}

	//get list position
	position, err := s.listPosRepo.GetByBoard(board.PublicID.String())
	if err != nil {
		return errors.New("failed to get list position: " + err.Error())
	}

	//update list order
	position.ListOrder = positions
	return s.listPosRepo.UpdateListOrder(position)
}

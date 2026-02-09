package repositories

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/google/uuid"
)

type listPositionRepository struct{}

type ListPositionRepository interface {
	GetByBoard(boardPublicId string) (*models.ListPosition, error)
	CreateOrUpdate(boardPublicId string, listOrder []uuid.UUID) error
	GetListOrder(boardPublicId string) ([]uuid.UUID, error)
	UpdateListOrder(position *models.ListPosition) error
}

func NewListPositionRepository() ListPositionRepository {
	return &listPositionRepository{}
}

func (r *listPositionRepository) GetByBoard(boardPublicId string) (*models.ListPosition, error) {
	var position models.ListPosition
	err := config.DB.
		Joins("JOIN boards on boards.internal_id = list_positions.board_internal_id").
		Where("boards.public_id = ?", boardPublicId).
		First(&position).Error
	return &position, err
}

func (r *listPositionRepository) CreateOrUpdate(boardPublicId string, listOrder []uuid.UUID) error {

	// On conflict, update the existing record
	return config.DB.Exec(`
		INSERT INTO list_positions (board_internal_id, list_order)
		SELECT internal_id, ? FROM boards WHERE public_id = ?
		ON CONFLICT (board_internal_id) DO UPDATE SET list_order = EXCLUDED.list_order
	`, listOrder, boardPublicId).Error
}

func (r *listPositionRepository) GetListOrder(boardPublicId string) ([]uuid.UUID, error) {
	position, err := r.GetByBoard(boardPublicId)
	if err != nil {
		return nil, err
	}
	return position.ListOrder, nil
}

func (r *listPositionRepository) UpdateListOrder(position *models.ListPosition) error {
	return config.DB.Model(position).Where("internal_id = ?",position.InternalID).Update("list_order", position.ListOrder).Error
}

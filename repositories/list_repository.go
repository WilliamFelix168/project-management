package repositories

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/google/uuid"
)

type ListRepository interface {
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPublicId string, position []string) error
	GetCardPosition(listPublicId string) ([]uuid.UUID, error)
}

type listRepository struct {
}

func NewListRepository() ListRepository {
	return &listRepository{}
}

func (r *listRepository) Create(list *models.List) error {
	return config.DB.Create(list).Error
}

func (r *listRepository) Update(list *models.List) error {
	return config.DB.Model(&models.List{}).Where("public_id = ?", list.PublicID).Updates(
		map[string]interface{}{
			"title": list.Title,
		},
	).Error
}

func (r *listRepository) Delete(id uint) error {
	return config.DB.Delete(&models.List{}, id).Error
}

func (r *listRepository) UpdatePosition(boardPublicId string, position []string) error {
	return config.DB.Model(&models.ListPosition{}).
		Where("board_internal_id = (select internal_id from boards where public_id = ?)", boardPublicId).
		Update("list_order", position).Error
}

func (r *listRepository) GetCardPosition(listPublicId string) ([]uuid.UUID, error) {
	var position models.CardPosition
	err := config.DB.Joins("JOIN lists on lists.internal_id = card_positions.list_internal_id").Where("lists.public_id = ?", listPublicId).First(&position).Error
	return position.CardOrder, err
}

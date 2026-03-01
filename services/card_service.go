package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models/types"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card) error
	Delete(id uint) error

	GetByListID(listPublicID string) ([]*models.Card, error)
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)
}

type CardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
}

func NewCardService(cardRepo repositories.CardRepository, listRepo repositories.ListRepository, userRepo repositories.UserRepository) *CardService {
	return &CardService{cardRepo, listRepo, userRepo}
}

func (s *CardService) Create(card *models.Card, listPublicID string) error {
	// 1. Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	//2. Set list_internal_id
	card.ListId = list.InternalId

	//3. Generate public_id jika belum ada
	if card.PublicId == uuid.Nil {
		card.PublicId = uuid.New()
	}
	card.CreatedAt = time.Now()

	//4. Mulai Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	//5. Simpan card
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w", err)
	}

	//6. Update atau buat card_position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).Where("list_internal_id = ?", list.InternalId).First(&position).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			position = models.CardPosition{
				PublicId:  uuid.New(),
				ListID:    list.InternalId,
				CardOrder: types.UUIDArray{card.PublicId},
			}
			if err := tx.Create(&position).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to find card position: %w", err)
		}
	} else {
		position.CardOrder = append(position.CardOrder, card.PublicId)
		if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", position.InternalId).Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update card position: %w", err)
		}
	}

	//7. Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

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
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error

	GetByListID(listPublicID string) ([]*models.Card, error)
	GetByID(id uint) (*models.Card, error)
	GetByPublicID(publicID string) (*models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
}

func NewCardService(cardRepo repositories.CardRepository, listRepo repositories.ListRepository, userRepo repositories.UserRepository) *cardService {
	return &cardService{cardRepo, listRepo, userRepo}
}

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	// 1. Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	//2. Set list_internal_id
	card.ListID = list.InternalID

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
	if err := tx.Model(&models.CardPosition{}).Where("list_internal_id = ?", list.InternalID).First(&position).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			position = models.CardPosition{
				PublicId:  uuid.New(),
				ListID:    list.InternalID,
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

func (s *cardService) Update(card *models.Card, listPublicID string) error {

	//ambil card lama
	existingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	//ambil list baru
	newList, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// Mulai Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// jika pindah list -> hapus dari posisi lama dan tambahkan ke list baru
	if existingCard.ListID != newList.InternalID {
		// Hapus dari posisi lama
		var oldPosition models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPosition).Error; err != nil {
			// Jika posisi lama tidak ditemukan, lanjutkan saja karena mungkin belum ada kartu lain di list tersebut
			filtered := make(types.UUIDArray, 0, len(oldPosition.CardOrder))
			for _, id := range oldPosition.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}

			//update
			if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", oldPosition.InternalId).Update("card_order", filtered).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update old card position: %w", err)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return fmt.Errorf("failed to find old card position: %w", err)
		}

		// Tambahkan ke posisi baru
		var newPosition models.CardPosition
		err := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPosition)
		if errors.Is(err.Error, gorm.ErrRecordNotFound) {
			newPosition = models.CardPosition{
				PublicId:  uuid.New(),
				ListID:    newList.InternalID,
				CardOrder: types.UUIDArray{existingCard.PublicID},
			}

			if err := tx.Create(&newPosition).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create new card position for new list : %w", err)
			}
		} else if err.Error == nil {
			//append
			updateOrder := append(newPosition.CardOrder, existingCard.PublicID)

			if err := tx.Model(&models.CardPosition{}).Where("internal_id = ?", newPosition.InternalId).Update("card_order", types.UUIDArray(updateOrder)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update new card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to find new card position: %w", err)
		}
	}

	//update card
	card.InternalId = existingCard.InternalId
	card.PublicID = existingCard.PublicID
	card.ListID = existingCard.ListID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update card: %w", err)
	}

	// Commit Transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

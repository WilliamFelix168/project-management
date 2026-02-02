package repositories

import (
	"time"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
)

// fungsi untuk mengelola operasi database terkait board
type BoardRepository interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
	AddMember(boardID uint, userIDs []uint) error
	RemoveMembers(boardID uint, userIDs []uint) error
	FindAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error)
}

// fungsi untuk mengimplementasi BoardRepository
type boardRepository struct {
}

// fungsi untuk mengembalikan objek boardRepository yang mengimplementasi BoardRepository
// terkait pemisahan interface dan impelmentasi
// agar lebih mudah dalam pengujian (testing)
// untuk encapsulasi kode, agar kode lebih terstruktur dan tidak bisa diakses langsung dari luar package, kalau mau pakai cukup NewBoardRepository()
// constructor
func NewBoardRepository() BoardRepository {
	return &boardRepository{}
}

// implementasi method dari interface BoardRepository
func (r *boardRepository) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}

func (r *boardRepository) Update(board *models.Board) error {
	return config.DB.Model(&models.Board{}).Where("public_id = ?", board.PublicID).Updates(map[string]interface{}{
		"title":       board.Title,
		"description": board.Description,
		"due_date":    board.DueDate,
	}).Error
}

func (r *boardRepository) FindByPublicID(publicID string) (*models.Board, error) {
	var board models.Board
	err := config.DB.Where("public_id = ?", publicID).First(&board).Error
	return &board, err
}

func (c *boardRepository) AddMember(boardID uint, userIDs []uint) error {

	if len(userIDs) == 0 {
		return nil
	}

	now := time.Now()
	var members []models.BoardMember
	for _, userID := range userIDs {
		members = append(members, models.BoardMember{
			BoardID:  int64(boardID),
			UserID:   int64(userID),
			JoinedAt: now,
		})
	}
	return config.DB.Create(&members).Error
}

func (r *boardRepository) RemoveMembers(boardID uint, userIDs []uint) error {

	if len(userIDs) == 0 {
		return nil
	}

	return config.DB.Where("board_internal_id = ? AND user_internal_id IN (?)", boardID, userIDs).Delete(&models.BoardMember{}).Error
}

func (r *boardRepository) FindAllByUserPaginate(userPublicID, filter, sort string, limit, offset int) ([]models.Board, int64, error) {
	var boards []models.Board
	var total int64

	query := config.DB.Model(&models.Board{}).Where("owner_public_id = ? or internal_id IN"+
		"(SELECT board_members.board_internal_id FROM board_members"+
		" JOIN users ON board_members.user_internal_id = users.internal_id "+
		"WHERE users.public_id = ?)", userPublicID, userPublicID)

	if filter != "" {
		query = query.Where("title ILIKE ?", "%"+filter+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//sorting
	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Limit(limit).Offset(offset).Find(&boards).Error; err != nil {
		return nil, 0, err
	}

	return boards, total, nil
}

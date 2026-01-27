package repositories

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
)

// fungsi untuk mengelola operasi database terkait board
type BoardRepository interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	FindByPublicID(publicID string) (*models.Board, error)
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

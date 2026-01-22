package repositories

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
)

// fungsi untuk mengelola operasi database terkait board
type BoardRepository interface {
	Create(board *models.Board) error
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

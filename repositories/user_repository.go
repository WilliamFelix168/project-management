package repositories

//utk mengelola operasi database terkait user

import (
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/config"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
)

// interface
// method Create untuk menambahkan user baru
type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
}

// struct yang mengimplementasi UserRepository
type userRepository struct{}

// fungsi untuk mengembalikan objek userRepository yang mengimplementasi UserRepository
// terkait pemisahan interface dan impelmentasi
// agar lebih mudah dalam pengujian (testing)
// untuk encapsulasi kode, agar kode lebih terstruktur dan tidak bisa diakses langsung dari luar package, kalau mau pakai cukup NewUserRepository()
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// implementasi method dari interface UserRepository
func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id = ?", publicID).First(&user).Error
	return &user, err
}

package repositories

//utk mengelola operasi database terkait user

import (
	"strings"

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
	FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint) error
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

func (r *userRepository) FindAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	//tujuan untuk mempermudah query
	db := config.DB.Model(&models.User{})

	//filtering
	if filter != "" {
		filterpattern := "%" + filter + "%"
		//ILIKE untuk case insensitive
		db = db.Where("name ILIKE ? OR email ILIKE ?", filterpattern, filterpattern)
	}

	//count total data sesuai filter
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//sorting
	if sort != "" {
		// misalnya sort = "name asc" atau  "-name (desc)""
		if sort == "-id" {
			sort = "-internal_id"
		} else if sort == "id" {
			sort = "internal_id"
		}

		if strings.HasPrefix(sort, "-") {
			sort = strings.TrimPrefix(sort, "-") + " DESC"
		} else {
			sort += " ASC"
		}

		db = db.Order(sort)

	}

	//pagination
	err := db.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err

}

func (r *userRepository) Update(user *models.User) error {
	return config.DB.Model(&models.User{}).Where("public_id = ?", user.PublicID).Updates(map[string]interface{}{
		"name": user.Name,
	}).Error
}

func (r *userRepository) Delete(id uint) error {
	return config.DB.Delete(&models.User{}, id).Error
}

package services

import (
	"errors"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/repositories"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/utils"
	"github.com/google/uuid"
)

//logika bisnis untuk user

type UserService interface {
	Register(user *models.User) error
	//method-method untuk user
}

type userService struct {
	repo repositories.UserRepository
	//dependency yang dibutuhkan service
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// (s *userService) artinya method ini milik struct userService
func (s *userService) Register(user *models.User) error {
	//logika bisnis untuk register user

	//cek email sudah terdaftar atau belum
	existingUser, _ := s.repo.FindByEmail(user.Email)

	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}

	//hashing password

	hased, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hased
	//Set role
	user.Role = "user"
	user.PublicID = uuid.New()

	// simpan user ke database lewat repository

	return s.repo.Create(user)
}

/*
	Controller
	↓
	UserService.Register()
	↓
	UserRepository.Create()
	↓
	Database

	Service = mikir
	Repo = eksekusi
	DB = nyimpen

	ANALOGI PALING WARAS 🧠

	Bayangin restoran:

	Controller = kasir
	Service = koki (mikir resep)
	Repository = tukang belanja bahan
	Database = gudang


	Kasir: “Bang, ada pesanan!”

	Koki: “Oke, resepnya gini, tapi bahannya beli dulu.”

	Tukang belanja: “Siap, gue ke gudang.”

*/

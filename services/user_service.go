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
	Login(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByPublicID(publicID string) (*models.User, error)
	GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
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

func (s *userService) Login(email, password string) (*models.User, error) {
	//cari user berdasarkan email
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("Invalid Credentials")
	}

	//bandingkan password
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("Invalid Credentials")
	}

	return user, nil
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetByPublicID(publicID string) (*models.User, error) {
	return s.repo.FindByPublicID(publicID)
}

func (s *userService) GetAllPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	return s.repo.FindAllPagination(filter, sort, limit, offset)
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

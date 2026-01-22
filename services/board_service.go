package services

import (
	"errors"

	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/models"
	"github.com/WilliamFelix168/learning-journey/tree/main/Golang/WPU/Project/project-management/repositories"
	"github.com/google/uuid"
)

type BoardService interface {
	Create(board *models.Board) error
}

type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo  repositories.UserRepository
}

func NewBoardService(boardRepo repositories.BoardRepository, userRepo repositories.UserRepository) BoardService {
	return &boardService{
		boardRepo, userRepo,
	}
}

func (s *boardService) Create(board *models.Board) error {
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("Owner Not Found")
	}

	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepo.Create(board)
}

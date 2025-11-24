package services

import (
	"errors"
	"monetix-be-api/internal/models"
	"monetix-be-api/internal/repositories"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(req models.CreateUserRequest) (*models.UserResponse, error)
	GetUser(id uuid.UUID) (*models.UserResponse, error)
	Authenticate(email, password string) (*models.User, error)
}

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s UserService) CreateUser(req models.CreateUserRequest) (*models.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

package service

import (
	"errors"
	"fmt"
	"otus/go-server-project/internal/models"
	"strconv"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserRepository interface {
	Login(string, string) (string, error)
	RegisterUser(u models.UserDTO) error
	Get(id int) (models.UserDTO, error)
}

type PasswordHasher interface {
	Hash(password string) string
}

type userService struct {
	hasher PasswordHasher
	repo   UserRepository
}

func NewUserService(r UserRepository, h PasswordHasher) *userService {
	return &userService{
		repo:   r,
		hasher: h,
	}
}

// Login authenticates a user with the given username and password.
// It returns a token if the credentials are valid, or an error if they are not.
func (s *userService) Login(login, password string) (string, error) {
	fmt.Printf("Login attempt with username: %s\n", login)
	if login == "" || password == "" {
		return "", ErrInvalidCredentials
	}
	pwd_hash := s.hasher.Hash(password)

	token, err := s.repo.Login(login, pwd_hash)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) RegisterUser(u models.User) (string, error) {
	m := models.MustConvertUserModelToDTO(u)
	m.PasswordHash = s.hasher.Hash(u.Password)
	if err := s.repo.RegisterUser(m); err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}
	return strconv.Itoa(m.ID), nil
}

func (s *userService) Get(id int) (models.User, error) {
	user, err := s.repo.Get(id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return models.ConvertUserDTOToModel(user), nil
}

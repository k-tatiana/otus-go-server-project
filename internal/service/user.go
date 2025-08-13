package service

import (
	"context"
	"errors"
	"fmt"

	"otus/go-server-project/internal/models"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserRepository interface {
	Login(context.Context, string, string) (string, error)
	Get(context.Context, string) (models.UserDTO, error)
	SearchUser(context.Context, string, string) ([]models.UserDTO, error)
	RegisterUser(context.Context, models.UserDTO) (string, error)
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
func (s *userService) Login(ctx context.Context, login, password string) (string, error) {
	fmt.Printf("Login attempt with username: %s\n", login)
	if login == "" || password == "" {
		return "", ErrInvalidCredentials
	}
	pwd_hash := s.hasher.Hash(password)

	token, err := s.repo.Login(ctx, login, pwd_hash)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) RegisterUser(ctx context.Context, u models.User) (string, error) {
	m := models.MustConvertUserModelToDTO(u)
	m.PasswordHash = s.hasher.Hash(u.Password)
	token, err := s.repo.RegisterUser(ctx, m)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}
	return token, nil
}

func (s *userService) Get(ctx context.Context, id string) (models.User, error) {
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return models.ConvertUserDTOToModel(user), nil
}

func (s *userService) SearchUser(ctx context.Context, name, surname string) ([]models.User, error) {
	users := make([]models.User, 0)
	usersDTO, err := s.repo.SearchUser(ctx, name, surname)
	if err != nil {
		return nil, err
	}
	for _, user := range usersDTO {
		u := models.ConvertUserDTOToModel(user)
		users = append(users, u)
	}

	return users, nil
}

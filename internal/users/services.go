package users

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type UserService struct {
	repository RepositoryInterface
	logger     *slog.Logger
}

func NewUserService(repo RepositoryInterface) *UserService {
	return &UserService{repository: repo}
}

// CreateUser - создает пользователя с хешированным паролем
func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*User, error) {
	const op = "internal.users.services.CreateUser"
	s.logger.With(slog.String("op", op))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Ошибка хеширования пароля в сервисе")
		return nil, err
	}

	user := &User{
		Login:        dto.Login,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repository.Create(ctx, user)
	if err != nil {
		s.logger.Error("Ошибка при создании пользователя в репозитории", slog.Any("err", err))
		return nil, err
	}
	return user, nil
}

// CheckPassword - проверяет, совпадает ли пароль с хешем
func (s *UserService) CheckPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

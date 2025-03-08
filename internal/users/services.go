package users

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type UserService struct {
	repository RepositoryInterface
	logger     *slog.Logger
}

func NewUserService(repo RepositoryInterface, logger *slog.Logger) *UserService {
	return &UserService{repository: repo, logger: logger}
}

// RegisterUser - создает пользователя с хешированным паролем
func (s *UserService) RegisterUser(ctx context.Context, dto UsersDTO) (*User, error) {
	const op = "internal.users.services.RegisterUser"
	s.logger = s.logger.With(slog.String("op", op))
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

// LoginUser Возвращает валидный ли пароль, и успешный ли запрос
func (s *UserService) LoginUser(ctx context.Context, userDTO UsersDTO) (isValid bool, ok bool) {
	const op = "internal.users.services.RegisterUser"
	s.logger = s.logger.With(slog.String("op", op))
	currentUser, err := s.repository.FindOne(ctx, userDTO.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Info(fmt.Sprintf("Пользователя %s не существует", userDTO.Login))
			return false, true
		}
		s.logger.Error(fmt.Sprintf("ошибка при поиске пользователя %s", userDTO.Login))
		return false, false
	}

	isValid = s.checkPassword(currentUser, userDTO.Password)
	if isValid {
		return isValid, true
	}
	return false, true

}

// CheckPassword - проверяет, совпадает ли пароль с хешем
func (s *UserService) checkPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

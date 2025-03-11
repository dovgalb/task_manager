package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"task-manager/internal/users/repo"
	"task-manager/pkg/clients/kafka"
	"time"
)

type UserService struct {
	logger     *slog.Logger
	repository repo.RepositoryInterface
	producer   *kafka.Producer
}

func NewUserService(logger *slog.Logger, repo repo.RepositoryInterface, producer *kafka.Producer) *UserService {
	return &UserService{repository: repo, logger: logger, producer: producer}
}

// RegisterUser - создает пользователя с хешированным паролем
func (s *UserService) RegisterUser(ctx context.Context, dto UsersDTO) (*repo.User, error) {
	const op = "internal.users.services.RegisterUser"
	s.logger = s.logger.With(slog.String("op", op))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Ошибка хеширования пароля в сервисе")
		return nil, err
	}

	user := &repo.User{
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

	message := fmt.Sprintf("Пользователь %s зарегестрирован", user.Login)
	if err := s.producer.SendMessage("key", message); err != nil {
		s.logger.Error("Ошибка отправки сообщения о зарегистрированном пользователе", slog.Any("err", err))
	}

	return user, nil
}

// AuthenticateUser Возвращает валидный ли пароль, и успешный ли запрос по логину
func (s *UserService) AuthenticateUser(ctx context.Context, userDTO UsersDTO) (*repo.User, error) {
	const op = "internal.users.services.RegisterUser"
	s.logger = s.logger.With(slog.String("op", op))
	currentUser, err := s.repository.FindOne(ctx, userDTO.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Info(fmt.Sprintf("Пользователя %s не существует", userDTO.Login))
			return nil, errors.New("пользователь не найден")
		}
		s.logger.Error(fmt.Sprintf("ошибка при поиске пользователя %s", userDTO.Login))
		return nil, err
	}

	isValidHash := s.checkPasswordHash(currentUser, userDTO.Password)
	if !isValidHash {
		return nil, errors.New("неверный пароль или логин")
	}

	message := fmt.Sprintf("Пользователь %s успешно аутентифицирован", currentUser.Login)
	if err := s.producer.SendMessage("key", message); err != nil {
		s.logger.Error("Ошибка отправки сообщения", slog.Any("err", err))
	}

	return currentUser, nil

}

// AuthenticateUserByID Возвращает валидный ли пароль, и успешный ли запрос по id
func (s *UserService) GetUserByID(ctx context.Context, id float64, password string) (*repo.User, error) {
	const op = "internal.users.services.RegisterUser"
	intID := int(id)

	s.logger = s.logger.With(slog.String("op", op))
	currentUser, err := s.repository.FindOneByID(ctx, intID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.logger.Info(fmt.Sprintf("Пользователя %d не существует", intID))
			return nil, errors.New("пользователь не найден")
		}
		s.logger.Error(fmt.Sprintf("ошибка при поиске пользователя %d", intID))
		return nil, err
	}

	isValidHash := s.checkPasswordHash(currentUser, password)
	if !isValidHash {
		return nil, errors.New("неверный пароль или логин")
	}
	return currentUser, nil

}

func (s *UserService) DeleteUser(ctx context.Context, user *repo.User) error {
	const op = "internal.users.services.RegisterUser"

	err := s.repository.Delete(ctx, user.ID)
	if err != nil {
		s.logger.Error("Ошибка удаления пользователя", slog.Any("err", err))
		return err
	}

	return nil

}

// CheckPassword - проверяет, совпадает ли пароль с хешем
func (s *UserService) checkPasswordHash(user *repo.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

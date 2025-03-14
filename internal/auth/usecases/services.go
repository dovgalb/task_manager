package usecases

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"task-manager/internal/auth/repo"
	"task-manager/pkg/logger/sl"
	"time"
)

var (
	ErrIncorrectCredentials = errors.New("неправильный логин или пароль")
)

type UserService struct {
	logger     *slog.Logger
	repository RepositoryInterface
	producer   Producer
}

func NewUserService(logger *slog.Logger, repo RepositoryInterface, producer Producer) *UserService {
	return &UserService{repository: repo, logger: logger, producer: producer}
}

// RegisterUser - создает пользователя с хешированным паролем
func (s *UserService) RegisterUser(ctx context.Context, dto UsersDTO) (*repo.User, error) {
	const op = "internal.users.services.RegisterUser"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("login", dto.Login),
	)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Ошибка хеширования пароля", sl.Err(err))
		return nil, fmt.Errorf("%s :%w", op, err)
	}

	user := &repo.User{
		Login:        dto.Login,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.repository.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repo.ErrUserExists) {
			return nil, fmt.Errorf("%s: %w", op, repo.ErrUserExists)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	message := fmt.Sprintf("Пользователь зарегестрирован")

	if err := s.producer.SendMessage("key", message); err != nil {
		log.Error("Ошибка отправки сообщения о зарегистрированном пользователе", sl.Err(err))
	}

	return user, nil
}

// AuthenticateUser Возвращает валидный ли пароль, и успешный ли запрос по логину
func (s *UserService) AuthenticateUser(ctx context.Context, userDTO UsersDTO) (*repo.User, error) {
	const op = "internal.users.services.RegisterUser"

	log := s.logger.With(slog.String("op", op))
	currentUser, err := s.repository.FindOne(ctx, userDTO.Login)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, repo.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	isValidHash := s.checkPasswordHash(currentUser, userDTO.Password)
	if !isValidHash {
		return nil, fmt.Errorf("%s: %w", op, ErrIncorrectCredentials)
	}

	message := fmt.Sprintf("Пользователь %s успешно аутентифицирован", currentUser.Login)
	if err := s.producer.SendMessage("key", message); err != nil {
		log.Error("Ошибка отправки сообщения", sl.Err(err))
	}

	return currentUser, nil

}

// GetUserByID Возвращает валидный ли пароль, и успешный ли запрос по id
func (s *UserService) GetUserByID(ctx context.Context, id float64, password string) (*repo.User, error) {
	const op = "internal.users.services.RegisterUser"
	intID := int(id)

	currentUser, err := s.repository.FindOneByID(ctx, intID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, fmt.Errorf("%s :%w", op, repo.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s :%w", op, err)
	}

	isValidHash := s.checkPasswordHash(currentUser, password)
	if !isValidHash {
		return nil, fmt.Errorf("%s :%w", op, ErrIncorrectCredentials)
	}
	return currentUser, nil

}

func (s *UserService) DeleteUser(ctx context.Context, user *repo.User) error {
	const op = "internal.users.services.DeleteUser"

	err := s.repository.Delete(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return fmt.Errorf("%s :%w", op, repo.ErrUserNotFound)
		}
		return fmt.Errorf("%s :%w", op, err)
	}

	return nil

}

// CheckPassword - проверяет, совпадает ли пароль с хешем
func (s *UserService) checkPasswordHash(user *repo.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

package usersUsecase

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	httpModels "github.com/cyber_bed/internal/models/http"
	"github.com/cyber_bed/internal/utils/crypto"
)

type UsersUsecase struct {
	usersRepository domain.UsersRepository
}

func NewUsersUsecase(r domain.UsersRepository) domain.UsersUsecase {
	return UsersUsecase{
		usersRepository: r,
	}
}

func (u UsersUsecase) CreateUser(user httpModels.User) (uint64, error) {
	if _, err := u.usersRepository.GetByUsername(user.Username); err == nil {
		return 0, errors.Wrapf(
			httpModels.ErrUserExists,
			"user already exists with username: %s",
			user.Username,
		)
	}

	hash, err := crypto.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = hash

	id, err := u.usersRepository.Create(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u UsersUsecase) GetBySessionID(sessionID string) (httpModels.User, error) {
	user, err := u.usersRepository.GetBySessionID(sessionID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return httpModels.User{}, errors.Wrapf(
				httpModels.ErrNotFound,
				"session with value: {%s} not found",
				sessionID,
			)
		}
		return httpModels.User{}, err
	}
	return httpModels.UserGormToHttp(user), nil
}

func (u UsersUsecase) GetUserIDBySessionID(sessionID string) (uint64, error) {
	usrID, err := u.usersRepository.GetUserIDBySessionID(sessionID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return 0, errors.Wrapf(
				httpModels.ErrNotFound,
				"session with value: {%s} not found",
				sessionID,
			)
		}
		return 0, err
	}
	return usrID, nil
}

func (u UsersUsecase) GetByUsername(username string) (httpModels.User, error) {
	user, err := u.usersRepository.GetByUsername(username)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return httpModels.User{}, errors.Wrapf(
				httpModels.ErrNotFound,
				"username with value: {%s} not found",
				username,
			)
		}
		return httpModels.User{}, err
	}
	return httpModels.UserGormToHttp(user), nil
}

func (u UsersUsecase) GetByID(userID uint64) (httpModels.User, error) {
	user, err := u.usersRepository.GetByID(userID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return httpModels.User{}, errors.Wrapf(
				httpModels.ErrNotFound,
				"username with id: {%d} not found",
				userID,
			)
		}
		return httpModels.User{}, err
	}
	return httpModels.UserGormToHttp(user), nil
}

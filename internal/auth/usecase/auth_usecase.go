package authUsecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/cyber_bed/internal/config"
	"github.com/cyber_bed/internal/domain"
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
	"github.com/cyber_bed/internal/utils/crypto"
)

type AuthUsecase struct {
	authRepository   domain.AuthRepository
	usersRepoisitory domain.UsersRepository

	config config.CookieSettings
}

func NewAuthUsecase(
	r domain.AuthRepository,
	u domain.UsersRepository,
	c config.CookieSettings,
) domain.AuthUsecase {
	return AuthUsecase{
		authRepository:   r,
		usersRepoisitory: u,
		config:           c,
	}
}

func (u AuthUsecase) generateCookie(userID uint64) gormModels.Cookie {
	return gormModels.Cookie{
		UserID: userID,
		Value:  uuid.New().String(),
		ExpireDate: time.Now().AddDate(
			int(u.config.ExpireDate.Years),
			int(u.config.ExpireDate.Months),
			int(u.config.ExpireDate.Days),
		),
	}
}

func (u AuthUsecase) Auth(sessionID string) error {
	if _, err := u.usersRepoisitory.GetBySessionID(sessionID); err != nil {
		return err
	}
	return nil
}

func (u AuthUsecase) SignUpByID(userID uint64) (string, error) {
	session, err := u.authRepository.CreateSession(u.generateCookie(userID))
	if err != nil {
		return "", err
	}
	return session, nil
}

func (u AuthUsecase) Login(login, password string) (string, uint64, error) {
	user, err := u.usersRepoisitory.GetByUsername(login)
	if err != nil {
		return "", 0, err
	}

	if !crypto.CheckHash(user.Password, password) {
		return "", 0, errors.Wrapf(
			httpModels.ErrIncorrectPassword,
			"incorrect password from user with login: %s",
			login,
		)
	}

	session, err := u.authRepository.CreateSession(u.generateCookie(user.ID))
	if err != nil {
		return "", 0, errors.Wrapf(err, "failed to create session for user with id: %d", user.ID)
	}
	return session, user.ID, nil
}

func (u AuthUsecase) Logout(sessionID string) error {
	return errors.Wrapf(
		u.authRepository.DeleteBySessionID(sessionID),
		"failed to delete user by session id %s",
		sessionID,
	)
}

package authMiddlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	httpAuth "github.com/cyber_bed/internal/auth/delivery"
	"github.com/cyber_bed/internal/domain"
)

type Middlewares struct {
	authUsecase  domain.AuthUsecase
	usersUsecase domain.UsersUsecase
}

func New(a domain.AuthUsecase, u domain.UsersUsecase) *Middlewares {
	return &Middlewares{
		authUsecase:  a,
		usersUsecase: u,
	}
}

func (m Middlewares) LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := httpAuth.GetCookie(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		if err = m.authUsecase.Auth(cookie.Value); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}

		return next(c)
	}
}

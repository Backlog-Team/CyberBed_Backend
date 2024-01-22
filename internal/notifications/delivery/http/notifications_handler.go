package httpNotifications

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	httpAuth "github.com/cyber_bed/internal/auth/delivery"
	"github.com/cyber_bed/internal/domain"
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type NotificationsHandler struct {
	notificaionsUsecase domain.NotificationsUsecase
	usersUsecase        domain.UsersUsecase
}

func NewNotificationsHandler(
	n domain.NotificationsUsecase,
	u domain.UsersUsecase,
) NotificationsHandler {
	return NotificationsHandler{
		notificaionsUsecase: n,
		usersUsecase:        u,
	}
}

func (h NotificationsHandler) GetNotifications(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	notifications, err := h.notificaionsUsecase.GetNotificationsByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, notifications)
}

func (h NotificationsHandler) DeleteNotification(c echo.Context) error {
	notificationID, err := strconv.ParseUint(c.Param("notificationID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if h.notificaionsUsecase.DeleteNotification(notificationID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h NotificationsHandler) DeleteCategoryNotification(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	stateToDelete := c.QueryParam("state")
	if stateToDelete != "" && stateToDelete != gormModels.NotificationStatusWaiting &&
		stateToDelete != gormModels.NotificationStatusDone &&
		stateToDelete != gormModels.NotificationStatusFinish {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			errors.New("providing status doest't exist"),
		)
	}

	if stateToDelete != "" {
		if err := h.notificaionsUsecase.DeleteNotificationByIDAndStatus(
			userID,
			gormModels.NotificationStatus(stateToDelete),
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

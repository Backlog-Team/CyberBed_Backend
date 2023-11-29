package foldersHandler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	httpAuth "github.com/cyber_bed/internal/auth/delivery"
	"github.com/cyber_bed/internal/domain"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type FoldersHandler struct {
	foldersUsecase       domain.FoldersUsecase
	usersUsecase         domain.UsersUsecase
	notificationsUsecase domain.NotificationsUsecase
}

func NewFoldersHandler(
	f domain.FoldersUsecase,
	a domain.UsersUsecase,
	n domain.NotificationsUsecase,
) FoldersHandler {
	return FoldersHandler{
		foldersUsecase:       f,
		usersUsecase:         a,
		notificationsUsecase: n,
	}
}

func (h FoldersHandler) CreateFolder(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	folderName := c.QueryParam("name")
	id, err := h.foldersUsecase.CreateFolder(httpModels.Folder{
		UserID:     userID,
		FolderName: folderName,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.UserID{ID: id})
}

func (h FoldersHandler) GetFolders(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	folders, err := h.foldersUsecase.GetFoldersByUserID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, folders)
}

func (h FoldersHandler) DeleteFolder(c echo.Context) error {
	folderID, err := strconv.ParseUint(c.Param("folderID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.foldersUsecase.DeleteFolderByID(folderID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h FoldersHandler) GetPlantsFromFolder(c echo.Context) error {
	folderID, err := strconv.ParseUint(c.Param("folderID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	plants, err := h.foldersUsecase.GetPlantsFromFolder(folderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, plants)
}

func (h FoldersHandler) AddPlantToFolder(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	folderID, err := strconv.ParseUint(c.Param("folderID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	duration := c.QueryParam("wateringTime")

	if err := h.foldersUsecase.AddPlantToFolder(folderID, plantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if len(duration) > 0 {
		if _, err = h.notificationsUsecase.CreateNotification(httpModels.Notification{
			UserID:         userID,
			PlantID:        plantID,
			FolderID:       folderID,
			ExpirationTime: duration,
		}); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h FoldersHandler) DeletePlantFromFolder(c echo.Context) error {
	folderID, err := strconv.ParseUint(c.Param("folderID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.foldersUsecase.DeletePlantFromFolder(folderID, plantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h FoldersHandler) UpdatePeriod(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	folderID, err := strconv.ParseUint(c.Param("folderID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	duration := c.QueryParam("wateringTime")
	nf, err := h.notificationsUsecase.UpdatePeriodNotification(httpModels.Notification{
		UserID:         userID,
		FolderID:       folderID,
		PlantID:        plantID,
		ExpirationTime: duration,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nf)
}

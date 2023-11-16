package foldersHandler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	httpAuth "github.com/cyber_bed/internal/auth/delivery"
	"github.com/cyber_bed/internal/domain"
	"github.com/cyber_bed/internal/models"
)

type FoldersHandler struct {
	foldersUsecase domain.FoldersUsecase
	usersUsecase   domain.UsersUsecase
}

func NewFoldersHandler(f domain.FoldersUsecase, a domain.UsersUsecase) FoldersHandler {
	return FoldersHandler{
		foldersUsecase: f,
		usersUsecase:   a,
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
	id, err := h.foldersUsecase.CreateFolder(models.Folder{
		UserID:     userID,
		FolderName: folderName,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, models.UserID{ID: id})
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
	return c.JSON(http.StatusOK, models.EmptyModel{})
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
	folderID, err := strconv.ParseUint(c.Param("folderID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.foldersUsecase.AddPlantToFolder(folderID, plantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, models.EmptyModel{})
}

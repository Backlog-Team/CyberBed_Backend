package httpPlants

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"

	httpAuth "github.com/cyber_bed/internal/auth/delivery"
	"github.com/cyber_bed/internal/domain"
	httpModels "github.com/cyber_bed/internal/models/http"
	"github.com/cyber_bed/internal/utils/decoding"
)

type PlantsHandler struct {
	plantsUsecase domain.PlantsUsecase
	usersUsecase  domain.UsersUsecase
}

func NewPlantsHandler(
	p domain.PlantsUsecase,
	u domain.UsersUsecase,
	pl domain.PlantsAPI,
) PlantsHandler {
	return PlantsHandler{
		plantsUsecase: p,
		usersUsecase:  u,
	}
}

func (h PlantsHandler) GetPlantFromAPI(c echo.Context) error {
	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	plant, err := h.plantsUsecase.GetPlantByID(plantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, plant)
}

func (h PlantsHandler) GetPlantImage(c echo.Context) error {
	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	plant, err := h.plantsUsecase.GetPlantByID(plantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	fileName, err := decoding.DecodeBase64(plant.Image)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer os.Remove(fileName)

	plant.Image = ""
	return c.Attachment(fileName, plant.Image)
}

func (h PlantsHandler) GetPlantsFromAPI(c echo.Context) error {
	pageStr := c.QueryParam("page")
	if pageStr == "" {
		plantName := c.QueryParam("name")
		plants, err := h.plantsUsecase.GetPlantByName(plantName)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusOK, plants)
	}

	pageNum, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	plants, err := h.plantsUsecase.GetPlantsPage(pageNum)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, plants)
}

func (h PlantsHandler) CreatePlant(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var recievedPlant httpModels.Plant
	if err := c.Bind(&recievedPlant); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	recievedPlant.ID = plantID
	recievedPlant.UserID = userID

	if err := h.plantsUsecase.AddPlant(recievedPlant); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h PlantsHandler) GetPlant(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	plant, err := h.plantsUsecase.GetPlant(userID, int64(plantID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	xiaomiPlant, err := h.plantsUsecase.GetPlantByID(plant.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, xiaomiPlant)
}

func (h PlantsHandler) GetPlants(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	plants, err := h.plantsUsecase.GetPlants(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, plants)
}

func (h PlantsHandler) DeletePlant(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	plantID, err := strconv.ParseUint(c.Param("plantID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err = h.plantsUsecase.DeletePlant(userID, plantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

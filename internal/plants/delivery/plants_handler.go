package httpPlants

import (
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"

	httpAuth "github.com/cyber_bed/internal/auth/delivery"
	"github.com/cyber_bed/internal/domain"
	httpModels "github.com/cyber_bed/internal/models/http"
	coder "github.com/cyber_bed/internal/utils/decoding"
	fileUtils "github.com/cyber_bed/internal/utils/files"
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

	return c.JSON(http.StatusOK, httpModels.XiaomiPlantGormToHttp(plant))
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

	fileName, err := coder.DecodeBase64(plant.Image)
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

	return c.JSON(http.StatusOK, httpModels.XiaomiPlantGormToHttp(xiaomiPlant))
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

func (h PlantsHandler) CreateCustomPlant(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	customPlant := httpModels.CustomPlant{
		UserID:    userID,
		PlantName: c.FormValue("plantName"),
		About:     c.FormValue("about"),
	}

	formdata, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	image, isImageProvided := formdata.File["image"]
	var extension string
	if isImageProvided {
		content, err := image[0].Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		defer content.Close()

		mimeType, err := fileUtils.GetMimeType(image[0])
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		extensions, err := mime.ExtensionsByType(mimeType)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if len(extensions) > 0 {
			extension = extensions[0]
		}

		data, err := io.ReadAll(content)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		customPlant.Image = string(data)
	}

	customPlantID, err := h.plantsUsecase.CreateCustomPlant(customPlant, extension)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.UserID{ID: customPlantID})
}

func (h PlantsHandler) UpdateCustomPlant(c echo.Context) error {
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

	customPlant := httpModels.CustomPlant{
		ID:        plantID,
		UserID:    userID,
		PlantName: c.FormValue("plantName"),
		About:     c.FormValue("about"),
	}

	formdata, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	image, isImageProvided := formdata.File["image"]
	var extension string
	if isImageProvided {
		content, err := image[0].Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		defer content.Close()

		mimeType, err := fileUtils.GetMimeType(image[0])
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		extensions, err := mime.ExtensionsByType(mimeType)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		if len(extensions) > 0 {
			extension = extensions[0]
		}

		data, err := io.ReadAll(content)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		customPlant.Image = string(data)
	}

	if err := h.plantsUsecase.UpdateCustomPlant(customPlant, extension); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h PlantsHandler) GetCustomPlants(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	customPlants, err := h.plantsUsecase.GetCustomPlants(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, customPlants)
}

func (h PlantsHandler) GetCustomPlant(c echo.Context) error {
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

	customPlant, err := h.plantsUsecase.GetCustomPlant(userID, plantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, customPlant)
}

func (h PlantsHandler) GetCustomPlantImage(c echo.Context) error {
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

	image, err := h.plantsUsecase.GetCustomPlantImage(userID, plantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	fileName, err := coder.DecodeBase64(image)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	defer os.Remove(fileName)

	image = ""
	return c.Attachment(fileName, image)
}

func (h PlantsHandler) DeleteCustomPlant(c echo.Context) error {
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

	if err := h.plantsUsecase.DeleteCustomPlant(userID, plantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h PlantsHandler) CreateSavedPlant(c echo.Context) error {
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

	if err := h.plantsUsecase.CreateSavedPlant(userID, plantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

func (h PlantsHandler) GetSavedPlants(c echo.Context) error {
	cookie, err := httpAuth.GetCookie(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	userID, err := h.usersUsecase.GetUserIDBySessionID(cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	plants, err := h.plantsUsecase.GetSavedPlants(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, plants)
}

func (h PlantsHandler) DeleteSavedPlant(c echo.Context) error {
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

	if err := h.plantsUsecase.DeleteSavedPlant(userID, plantID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, httpModels.EmptyModel{})
}

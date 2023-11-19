package plantsUsecase

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type PlantsUsecase struct {
	plantsRepository domain.PlantsRepository
}

func NewPlansUsecase(p domain.PlantsRepository, api domain.PlantsAPI) domain.PlantsUsecase {
	return PlantsUsecase{
		plantsRepository: p,
	}
}

func (u PlantsUsecase) GetPlantByID(plantID uint64) (gormModels.XiaomiPlant, error) {
	plant, err := u.plantsRepository.GetPlantByID(plantID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return gormModels.XiaomiPlant{}, errors.Wrapf(
				httpModels.ErrNotFound,
				"plant with id: {%d} not found",
				plantID,
			)
		}
		return gormModels.XiaomiPlant{}, err
	}
	return plant, nil
}

func (u PlantsUsecase) GetPlantByName(plantName string) ([]httpModels.XiaomiPlant, error) {
	plants, err := u.plantsRepository.GetByPlantName(plantName)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.Wrapf(
				httpModels.ErrNotFound,
				"plant with name: {%s} not found",
				plantName,
			)
		}
		return nil, err
	}

	httpPlants := make([]httpModels.XiaomiPlant, 0)
	for _, pl := range plants {
		httpPlants = append(httpPlants, httpModels.XiaomiPlantGormToHttp(pl))
	}
	return httpPlants, nil
}

func (u PlantsUsecase) GetPlantsPage(pageNum uint64) ([]httpModels.XiaomiPlant, error) {
	plants, err := u.plantsRepository.GetPlantsPage(pageNum)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.Wrapf(
				httpModels.ErrNotFound,
				"cannot get page number: {%d}",
				pageNum,
			)
		}
		return nil, err
	}

	httpPlants := make([]httpModels.XiaomiPlant, 0)
	for _, pl := range plants {
		httpPlants = append(httpPlants, httpModels.XiaomiPlantGormToHttp(pl))
	}
	return httpPlants, nil
}

func (u PlantsUsecase) AddPlant(plant httpModels.Plant) error {
	if err := u.plantsRepository.AddUserPlantsRelations(plant.UserID, []int64{int64(plant.ID)}); err != nil {
		return err
	}
	return nil
}

func (u PlantsUsecase) GetPlant(userID uint64, plantID int64) (httpModels.Plant, error) {
	plants, err := u.plantsRepository.GetPlantsByUserID(userID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return httpModels.Plant{}, errors.Wrapf(
				httpModels.ErrNotFound,
				"plants of user with id: {%d} not found",
				userID,
			)
		}
	}

	if !slices.Contains(plants.PlantsID, plantID) {
		return httpModels.Plant{}, errors.Wrapf(
			httpModels.ErrNotFound,
			"Plant with id: {%d} of user: {%d} not found",
			plantID,
			userID,
		)
	}

	return httpModels.Plant{
		ID:     uint64(plantID),
		UserID: plants.UserID,
	}, nil
}

func (u PlantsUsecase) GetPlants(userID uint64) ([]httpModels.XiaomiPlant, error) {
	plantsIDs, err := u.plantsRepository.GetPlantsByUserID(userID)
	if err != nil {
		return nil, err
	}

	pl := plantsIDs.PlantsID
	plants := make([]httpModels.XiaomiPlant, 0)
	for _, p := range pl {
		curPlant, err := u.plantsRepository.GetPlantByID(uint64(p))
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				return nil, errors.Wrapf(httpModels.ErrNotFound, "Plant with id {%d} not found", p)
			}
			return nil, err
		}
		plants = append(plants, httpModels.XiaomiPlantGormToHttp(curPlant))
	}

	return plants, nil
}

func (u PlantsUsecase) DeletePlant(userID, plantID uint64) error {
	user, err := u.plantsRepository.GetPlantsByUserID(userID)
	if err != nil {
		return err
	}

	indexToDel := -1
	for index, plntID := range user.PlantsID {
		if plantID == uint64(plntID) {
			indexToDel = index
			break
		}
	}
	if indexToDel == -1 {
		return errors.Wrapf(
			httpModels.ErrNotFound,
			"plant with id: %d of user with id: %d was not found",
			plantID,
			userID,
		)
	}

	user.PlantsID = append(user.PlantsID[:indexToDel], user.PlantsID[indexToDel+1:]...)

	if err := u.plantsRepository.UpdateUserPlantsRelation(user); err != nil {
		return err
	}
	return nil
}

package plantsUsecase

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	"github.com/cyber_bed/internal/models"
)

type PlantsUsecase struct {
	plantsRepository domain.PlantsRepository
	// plantsAPI        domain.PlantsAPI
}

func NewPlansUsecase(p domain.PlantsRepository, api domain.PlantsAPI) domain.PlantsUsecase {
	return PlantsUsecase{
		plantsRepository: p,
		// plantsAPI:        api,
	}
}

func (u PlantsUsecase) GetPlantByID(plantID uint64) (models.XiaomiPlant, error) {
	plant, err := u.plantsRepository.GetPlantByID(plantID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return models.XiaomiPlant{}, errors.Wrapf(
				models.ErrNotFound,
				"plant with id: {%d} not found",
				plantID,
			)
		}
		return models.XiaomiPlant{}, err
	}
	return plant, nil
}

func (u PlantsUsecase) GetPlantByName(plantName string) ([]models.XiaomiPlant, error) {
	plants, err := u.plantsRepository.GetByPlantName(plantName)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.Wrapf(
				models.ErrNotFound,
				"plant with name: {%s} not found",
				plantName,
			)
		}
		return nil, err
	}
	return plants, nil
}

func (u PlantsUsecase) GetPlantsPage(pageNum uint64) ([]models.XiaomiPlant, error) {
	plants, err := u.plantsRepository.GetPlantsPage(pageNum)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, errors.Wrapf(
				models.ErrNotFound,
				"cannot get page number: {%d}",
				pageNum,
			)
		}
		return nil, err
	}
	return plants, nil
}

func (u PlantsUsecase) AddPlant(plant models.Plant) error {
	if err := u.plantsRepository.AddUserPlantsRelations(plant.UserID, []int64{int64(plant.ID)}); err != nil {
		return err
	}
	return nil
}

func (u PlantsUsecase) GetPlant(userID uint64, plantID int64) (models.Plant, error) {
	plants, err := u.plantsRepository.GetPlantsByUserID(userID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return models.Plant{}, errors.Wrapf(
				models.ErrNotFound,
				"plants of user with id: {%d} not found",
				userID,
			)
		}
	}

	if !slices.Contains(plants.PlantsID, plantID) {
		return models.Plant{}, errors.Wrapf(
			models.ErrNotFound,
			"Plant with id: {%d} of user: {%d} not found",
			plantID,
			userID,
		)
	}

	return models.Plant{
		ID:     uint64(plantID),
		UserID: plants.UserID,
	}, nil
}

func (u PlantsUsecase) GetPlants(userID uint64) ([]models.Plant, error) {
	plantsIDs, err := u.plantsRepository.GetPlantsByUserID(userID)
	if err != nil {
		return nil, err
	}

	pl := plantsIDs.PlantsID
	plants := make([]models.Plant, 0)
	for _, p := range pl {
		plants = append(plants, models.Plant{
			ID: uint64(p),
		})
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
			models.ErrNotFound,
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

package plantsUsecase

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
	coder "github.com/cyber_bed/internal/utils/decoding"
)

type PlantsUsecase struct {
	plantsRepository        domain.PlantsRepository
	notificationsRepository domain.NotificationsRepository
	foldersRepository       domain.FoldersRepository
}

func NewPlansUsecase(
	p domain.PlantsRepository,
	f domain.FoldersRepository,
	n domain.NotificationsRepository,
) domain.PlantsUsecase {
	return PlantsUsecase{
		plantsRepository:        p,
		foldersRepository:       f,
		notificationsRepository: n,
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

func (u PlantsUsecase) GetPlants(userID uint64) (map[uint64]httpModels.XiaomiPlant, error) {
	plantsIDs, err := u.plantsRepository.GetPlantsByUserID(userID)
	if err != nil {
		return nil, err
	}

	pl := plantsIDs.PlantsID
	plants := make(map[uint64]httpModels.XiaomiPlant, 0)
	for _, p := range pl {
		curPlant, err := u.plantsRepository.GetPlantByID(uint64(p))
		if err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				return nil, errors.Wrapf(httpModels.ErrNotFound, "Plant with id {%d} not found", p)
			}
			return nil, err
		}
		plants[curPlant.ID] = httpModels.XiaomiPlantGormToHttp(curPlant)
		// plants = append(plants, httpModels.XiaomiPlantGormToHttp(curPlant))
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

func (u PlantsUsecase) CreateCustomPlant(
	plant httpModels.CustomPlant,
	extension string,
) (uint64, error) {
	updatedPlant := plant
	if len(plant.Image) > 0 {
		encodedImage, err := coder.EncodeToBase64(plant.Image, extension)
		if err != nil {
			return 0, errors.Wrap(
				httpModels.ErrNoImages,
				"cannot create plant due to corrupted image",
			)
		}
		updatedPlant.Image = encodedImage
	}

	cPlantID, err := u.plantsRepository.CreateCustomPlant(updatedPlant)
	if err != nil {
		return 0, err
	}
	return cPlantID, nil
}

func (u PlantsUsecase) UpdateCustomPlant(plant httpModels.CustomPlant, extension string) error {
	updatedPlant := plant
	if len(plant.Image) > 0 {
		encodedImage, err := coder.EncodeToBase64(plant.Image, extension)
		if err != nil {
			return errors.Wrap(httpModels.ErrNoImages, "cannot create plant due to corrupted image")
		}
		updatedPlant.Image = encodedImage
	}
	return u.plantsRepository.UpdateCustomPlant(updatedPlant)
}

func (u PlantsUsecase) GetCustomPlants(userID uint64) ([]httpModels.CustomPlant, error) {
	cPlants := make([]httpModels.CustomPlant, 0)
	cGormPlants, err := u.plantsRepository.GetCustomPlants(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cPlants, nil
		}
		return nil, err
	}

	for _, pl := range cGormPlants {
		cPlants = append(cPlants, httpModels.CustomPlantGormToHttp(pl))
	}
	return cPlants, nil
}

func (u PlantsUsecase) GetCustomPlant(userID, plantID uint64) (httpModels.CustomPlant, error) {
	cGormPlant, err := u.plantsRepository.GetCustomPlant(userID, plantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpModels.CustomPlant{}, errors.Wrapf(
				httpModels.ErrNotFound,
				"custom plant with user_id {%d} and plant_id {%d} not found",
				userID,
				plantID,
			)
		}
		return httpModels.CustomPlant{}, err
	}
	return httpModels.CustomPlantGormToHttp(cGormPlant), nil
}

func (u PlantsUsecase) GetCustomPlantImage(userID, plantID uint64) (string, error) {
	cGormPlant, err := u.plantsRepository.GetCustomPlant(userID, plantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.Wrapf(
				httpModels.ErrNotFound,
				"custom plant with user_id {%d} and plant_id {%d} not found",
				userID,
				plantID,
			)
		}
		return "", err
	}
	return cGormPlant.Image, nil
}

func (u PlantsUsecase) DeleteCustomPlant(userID, plantID uint64) error {
	return u.plantsRepository.DeleteCustomPlant(userID, plantID)
}

func (u PlantsUsecase) CreateSavedPlant(userID, plantID uint64) error {
	_, err := u.plantsRepository.GetSavedPlantByIDs(userID, plantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u.plantsRepository.CreateSavedPlant(userID, plantID)
		}
		return err
	}

	return errors.Wrapf(
		httpModels.ErrRecordExists,
		"plant {%d} was saved earlier by user {%d}",
		plantID,
		userID,
	)
}

func (u PlantsUsecase) GetSavedPlants(userID uint64) ([]httpModels.XiaomiPlant, error) {
	savedPlants, err := u.plantsRepository.GetSavedPlants(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []httpModels.XiaomiPlant{}, nil
		}
		return []httpModels.XiaomiPlant{}, err
	}

	resPlants := make([]httpModels.XiaomiPlant, 0)
	for _, pl := range savedPlants {
		recievedPlant, err := u.plantsRepository.GetPlantByID(pl.PlantID)
		if err != nil {
			return nil, err
		}
		resPlants = append(resPlants, httpModels.XiaomiPlantGormToHttp(recievedPlant))

		nf, err := u.notificationsRepository.GetNotificationsByUserPlantID(userID, pl.PlantID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		resPlants[len(resPlants)-1].WateringPeriod = nf.Period

		ch, err := u.plantsRepository.GetChannelByUserAndPlantID(userID, pl.PlantID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		resPlants[len(resPlants)-1].ChannelID = ch
		resPlants[len(resPlants)-1].IsCustom = false
	}

	customPlants, err := u.plantsRepository.GetCustomPlants(userID)
	if err != nil {
		return []httpModels.XiaomiPlant{}, err
	}
	for _, cp := range customPlants {
		resPlants = append(resPlants, httpModels.XiaomiPlant{
			ID:      cp.ID,
			PlantID: cp.PlantName,
			Basic: struct {
				FloralLanguage string `             json:"floral_language"`
				Origin         string `             json:"origin"`
				Production     string `             json:"production"`
				Category       string `             json:"category"`
				Blooming       string `             json:"blooming"`
				Color          string `             json:"color"`
			}{
				FloralLanguage: cp.About,
			},
			IsCustom: true,
		})
	}

	return resPlants, nil
}

func (u PlantsUsecase) DeleteSavedPlant(userID, plantID uint64) error {
	return u.plantsRepository.DeleteSavedPlant(userID, plantID)
}

func (u PlantsUsecase) GetSavedFieldOfPlant(
	plant httpModels.XiaomiPlant,
	userID uint64,
) (bool, error) {
	likedPlants, err := u.GetSavedPlants(userID)
	if err != nil {
		return false, err
	}
	for _, v := range likedPlants {
		if v.ID == plant.ID {
			return true, nil
		}
	}
	return false, nil
}

func (f PlantsUsecase) CreateChannel(plantID, channelID, userID uint64) (uint64, error) {
	_, err := f.plantsRepository.GetChannelByUserAndPlantID(userID, plantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			id, err := f.plantsRepository.CreateChannel(plantID, channelID, userID)
			if err != nil {
				return 0, err
			}
			return id, nil
		}
		return 0, err
	}

	return 0, errors.Wrapf(
		err,
		"channel with id {%d} and user_id {%d} already exists",
		channelID,
		userID,
	)
}

func (u PlantsUsecase) GetChannelByUserAndPlantID(userID, plantID uint64) (uint64, error) {
	return u.plantsRepository.GetChannelByUserAndPlantID(userID, plantID)
}

func (u PlantsUsecase) UpdateChannel(userID, plantID, channelID uint64) error {
	return u.plantsRepository.UpdateChannelByUserAndPlantID(userID, plantID, channelID)
}

func (u PlantsUsecase) SetUserPlantFields(
	plant httpModels.XiaomiPlant,
	userID uint64,
) (httpModels.XiaomiPlant, error) {
	resPlant := plant
	isSaved, err := u.GetSavedFieldOfPlant(plant, userID)
	if err != nil {
		return httpModels.XiaomiPlant{}, err
	}
	resPlant.IsSaved = isSaved

	// Check if plant was saved
	foldersToCheck, err := u.foldersRepository.GetFolders(userID)
	if err != nil {
		return httpModels.XiaomiPlant{}, err
	}
	for _, f := range foldersToCheck {
		pl, err := u.foldersRepository.GetPlantsID(f.ID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return httpModels.XiaomiPlant{}, err
		}

		var fids []uint64
		for _, v := range pl {
			fids = append(fids, v)
		}

		if slices.Contains(fids, plant.ID) {
			resPlant.IsLiked = true
			resPlant.FolderSaved = append(resPlant.FolderSaved, httpModels.FolderGormToHttp(f))
		}
	}

	savedPlants, err := u.GetSavedPlants(userID)
	if err != nil {
		return httpModels.XiaomiPlant{}, err
	}

	for _, v := range savedPlants {
		if v.ID == plant.ID {
			resPlant.WateringPeriod = v.WateringPeriod
			resPlant.ChannelID = v.ChannelID
			break
		}
	}

	return resPlant, nil
}

func (u PlantsUsecase) SetUserPlantsFields(
	plants []httpModels.XiaomiPlant,
	userID uint64,
) ([]httpModels.XiaomiPlant, error) {
	resPlants := make([]httpModels.XiaomiPlant, 0)

	for _, plant := range plants {
		resPlant, err := u.SetUserPlantFields(plant, userID)
		if err != nil {
			return []httpModels.XiaomiPlant{}, err
		}

		resPlants = append(resPlants, resPlant)
	}

	return resPlants, nil
}

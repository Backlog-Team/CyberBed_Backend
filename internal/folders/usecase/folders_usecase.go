package foldersUsecase

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	httpModels "github.com/cyber_bed/internal/models/http"
)

type FoldersUsecase struct {
	foldersRepository domain.FoldersRepository
	plantsRepository  domain.PlantsRepository
}

func NewFoldersUsecase(
	f domain.FoldersRepository,
	p domain.PlantsRepository,
) domain.FoldersUsecase {
	return FoldersUsecase{
		foldersRepository: f,
		plantsRepository:  p,
	}
}

func (f FoldersUsecase) CreateFolder(folder httpModels.Folder) (uint64, error) {
	folderID, err := f.foldersRepository.CreateFolder(folder)
	if err != nil {
		return 0, err
	}
	return folderID, nil
}

func (f FoldersUsecase) GetFolderByID(id uint64) (httpModels.Folder, error) {
	folder, err := f.foldersRepository.GetFolder(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpModels.Folder{}, errors.Wrapf(
				httpModels.ErrNotFound,
				"Folder with id {%d} not found",
				id,
			)
		}
		return httpModels.Folder{}, err
	}
	return httpModels.FolderGormToHttp(folder), nil
}

func (f FoldersUsecase) GetFoldersByUserID(userID uint64) ([]httpModels.Folder, error) {
	folders, err := f.foldersRepository.GetFolders(userID)
	if err != nil {
		return []httpModels.Folder{}, err
	}

	httpFolders := make([]httpModels.Folder, 0)
	for _, f := range folders {
		httpFolders = append(httpFolders, httpModels.FolderGormToHttp(f))
	}
	return httpFolders, nil
}

func (f FoldersUsecase) GetPlantsFromFolder(folderID uint64) ([]httpModels.XiaomiPlant, error) {
	plantsID, err := f.foldersRepository.GetPlantsID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []httpModels.XiaomiPlant{}, nil
		}
		return nil, err
	}

	var resPlants []httpModels.XiaomiPlant
	for _, id := range plantsID {
		plant, err := f.plantsRepository.GetPlantByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.Wrapf(
					httpModels.ErrNotFound,
					"plant with id {%d} not found",
					id,
				)
			}
			return nil, err
		}
		resPlants = append(resPlants, httpModels.XiaomiPlantGormToHttp(plant))
	}
	return resPlants, nil
}

func (f FoldersUsecase) DeleteFolderByID(id uint64) error {
	return f.foldersRepository.DeleteFolder(id)
}

func (f FoldersUsecase) AddPlantToFolder(folderID, plantID uint64) error {
	if _, err := f.plantsRepository.GetPlantByID(plantID); err != nil {
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrapf(httpModels.ErrNotFound,
					"plant with id {%d} not found",
					plantID,
				)
			}
			return err
		}
	}
	return f.foldersRepository.AddPlantToFolder(folderID, plantID)
}

func (f FoldersUsecase) DeletePlantFromFolder(folderID, plantID uint64) error {
	return f.foldersRepository.UpdateFolderPlant(folderID, plantID)
}

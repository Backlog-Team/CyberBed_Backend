package foldersUsecase

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	"github.com/cyber_bed/internal/models"
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

func (f FoldersUsecase) CreateFolder(folder models.Folder) (uint64, error) {
	folderID, err := f.foldersRepository.CreateFolder(folder)
	if err != nil {
		return 0, err
	}
	return folderID, nil
}

func (f FoldersUsecase) GetFolderByID(id uint64) (models.Folder, error) {
	folder, err := f.foldersRepository.GetFolder(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Folder{}, errors.Wrapf(
				models.ErrNotFound,
				"Folder with id {%d} not found",
				id,
			)
		}
		return models.Folder{}, err
	}
	return folder, nil
}

func (f FoldersUsecase) GetFoldersByUserID(userID uint64) ([]models.Folder, error) {
	folders, err := f.foldersRepository.GetFolders(userID)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (f FoldersUsecase) GetPlantsFromFolder(folderID uint64) ([]models.XiaomiPlant, error) {
	plantsID, err := f.foldersRepository.GetPlantsID(folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(
				models.ErrNotFound,
				"Folder with id {%d} not found",
				folderID,
			)
		}
		return nil, err
	}

	var resPlants []models.XiaomiPlant
	for _, id := range plantsID {
		plant, err := f.plantsRepository.GetPlantByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.Wrapf(
					models.ErrNotFound,
					"plant with id {%d} not found",
					id,
				)
			}
			return nil, err
		}
		resPlants = append(resPlants, plant)
	}
	return resPlants, nil
}

func (f FoldersUsecase) DeleteFolderByID(id uint64) error {
	if err := f.foldersRepository.DeleteFolder(id); err != nil {
		return err
	}
	return nil
}

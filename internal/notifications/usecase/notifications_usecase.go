package notificationsUsecase

import (
	"errors"

	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
	"github.com/cyber_bed/internal/utils/converter"
)

type NotificationsUsecase struct {
	notificationsRepository domain.NotificationsRepository
}

func NewNotificationsUsecase(r domain.NotificationsRepository) domain.NotificationsUsecase {
	return NotificationsUsecase{
		notificationsRepository: r,
	}
}

func (u NotificationsUsecase) CreateNotification(
	notification httpModels.Notification,
) (uint64, error) {
	futureTime, err := converter.StringToTime(notification.ExpirationTime)
	if err != nil {
		return 0, err
	}

	id, err := u.notificationsRepository.CreateNotification(gormModels.Notification{
		UserID:         notification.UserID,
		PlantID:        notification.PlantID,
		ExpirationTime: futureTime,
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u NotificationsUsecase) GetNotificationsByUserID(
	userID uint64,
) ([]gormModels.Notification, error) {
	notifications, err := u.notificationsRepository.GetNotificationsByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []gormModels.Notification{}, nil
		}
		return nil, err
	}

	return notifications, nil
}

func (u NotificationsUsecase) GetNotificationsByUserIDAndStatus(
	userID uint64,
	status gormModels.NotificationStatus,
) ([]gormModels.Notification, error) {
	notifications, err := u.notificationsRepository.GetNotificationsByUserIDAndStatus(
		userID,
		status,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []gormModels.Notification{}, nil
		}
		return nil, err
	}

	return notifications, nil
}

func (u NotificationsUsecase) UpdateNotificationStatus(
	id uint64,
	status gormModels.NotificationStatus,
) error {
	return u.notificationsRepository.UpdateNotificationStatus(id, status)
}

func (u NotificationsUsecase) GetNotificationByID(id uint64) (httpModels.Notification, error) {
	return httpModels.Notification{}, nil
}

func (u NotificationsUsecase) DeleteNotification(id uint64) error {
	return u.notificationsRepository.DeleteNotification(id)
}

func (u NotificationsUsecase) DeleteNotificationByIDAndStatus(
	id uint64,
	status gormModels.NotificationStatus,
) error {
	return u.notificationsRepository.DeleteNotificationByIDAndStatus(id, status)
}

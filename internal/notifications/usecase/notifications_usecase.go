package notificationsUsecase

import (
	"errors"
	"time"

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
) (httpModels.Notification, error) {
	futureTime, err := converter.StringToTime(time.Now(), notification.ExpirationTime)
	if err != nil {
		return httpModels.Notification{}, err
	}

	nf, err := u.notificationsRepository.CreateNotification(gormModels.Notification{
		UserID:         notification.UserID,
		PlantID:        notification.PlantID,
		FolderID:       notification.FolderID,
		TimeStart:      time.Now(),
		ExpirationTime: futureTime,
		Period:         notification.ExpirationTime,
		Status:         gormModels.NotificationStatusWaiting,
	})
	if err != nil {
		return httpModels.Notification{}, err
	}
	return httpModels.NotificationGormToHttp(nf), nil
}

func (u NotificationsUsecase) GetNotificationsByUserID(
	userID uint64,
) ([]httpModels.Notification, error) {
	notifications, err := u.notificationsRepository.GetNotificationsByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []httpModels.Notification{}, nil
		}
		return nil, err
	}

	httpNotifications := make([]httpModels.Notification, 0)
	for _, n := range notifications {
		httpNotifications = append(httpNotifications, httpModels.NotificationGormToHttp(n))
	}

	return httpNotifications, nil
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
	notification, err := u.notificationsRepository.GetNotificationByID(id)
	if err != nil {
		return httpModels.Notification{}, err
	}
	return httpModels.NotificationGormToHttp(notification), nil
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

func (u NotificationsUsecase) UpdatePeriodNotification(
	notification httpModels.Notification,
) (httpModels.Notification, error) {
	nf, err := u.notificationsRepository.GetWaitingNotification(
		notification.UserID,
		notification.FolderID,
		notification.PlantID,
	)
	if err != nil {
		return httpModels.Notification{}, err
	}

	newTargetTime, err := converter.StringToTime(nf.TimeStart, notification.ExpirationTime)
	if err != nil {
		return httpModels.Notification{}, err
	}

	updNt, err := u.notificationsRepository.UpdatePeriodNotification(gormModels.Notification{
		UserID:         notification.UserID,
		PlantID:        notification.PlantID,
		FolderID:       notification.FolderID,
		TimeStart:      time.Now(),
		ExpirationTime: newTargetTime,
		Period:         notification.ExpirationTime,
	})
	if err != nil {
		return httpModels.Notification{}, err
	}

	return httpModels.NotificationGormToHttp(updNt), err
}

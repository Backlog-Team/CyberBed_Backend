package notificationsWS

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
	"gorm.io/gorm"

	"github.com/cyber_bed/internal/domain"
	gormModels "github.com/cyber_bed/internal/models/gorm"
	httpModels "github.com/cyber_bed/internal/models/http"
	logger "github.com/cyber_bed/pkg"
)

type WebSocket struct {
	u         domain.NotificationsUsecase
	wsUpgrade *websocket.Upgrader
}

type status struct {
	CloseReader chan interface{}
	CloseWriter chan interface{}
}

func NewWebSocket(up *websocket.Upgrader, u domain.NotificationsUsecase) WebSocket {
	return WebSocket{
		u:         u,
		wsUpgrade: up,
	}
}

func (h WebSocket) Handler(w http.ResponseWriter, r *http.Request) {
	log := logger.GetInstance().Logrus

	c, err := h.wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer c.Close()

	params := r.URL.Query()
	userID, err := strconv.ParseUint(params.Get("userID"), 10, 64)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Error("Bad Request")
		return
	}

	closeReaderSignal := make(chan interface{})
	closeWriterSignal := make(chan interface{})
	s := &status{
		CloseWriter: closeWriterSignal,
		CloseReader: closeReaderSignal,
	}

	firstLoop := true
	for {
		select {
		case <-s.CloseWriter:
			log.Warn("writer done")
			return
		default:
			if firstLoop {
				notifications, err := h.u.GetNotificationsByUserID(userID)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					handleError(tracerr.Wrap(err), c, log)
				}
				if len(notifications) > 0 {
					err = c.WriteJSON(notifications)
					handleError(tracerr.Wrap(err), c, log)
				}
				firstLoop = false
			}

			notifications, err := h.u.GetNotificationsByUserIDAndStatus(
				userID,
				gormModels.NotificationStatusDone,
			)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				handleError(tracerr.Wrap(err), c, log)
			}

			if len(notifications) > 0 {
				// If we have notifications with expired time
				err = c.WriteJSON(notifications)
				handleError(tracerr.Wrap(err), c, log)
				for _, n := range notifications {
					newNotification, err := h.u.CreateNotification(httpModels.Notification{
						UserID:         n.UserID,
						FolderID:       n.FolderID,
						PlantID:        n.PlantID,
						ExpirationTime: n.Period,
					})
					if err != nil {
						handleError(tracerr.Wrap(err), c, log)
						continue
					}
					err = c.WriteJSON(newNotification)
					handleError(tracerr.Wrap(err), c, log)
				}
			} else {
				// Trigger our hook to check state of notifications
				if _, err = h.u.GetNotificationsByUserIDAndStatus(
					userID,
					gormModels.NotificationStatusWaiting,
				); err != nil {
					handleError(tracerr.Wrap(err), c, log)
				}
			}

			time.Sleep(time.Second)
		}
	}
}

func handleError(err error, c *websocket.Conn, l *logrus.Logger) {
	if err != nil {
		l.WithFields(
			logrus.Fields{
				"error": err,
			},
		).Error(err)
		if err != websocket.ErrCloseSent {
			err = c.Close()
			l.Warn("connection closed: ", err)
		}
	}
}

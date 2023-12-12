package main

import (
	"flag"
	"fmt"
	defaultLogger "log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/cyber_bed/internal/config"
	notificationsWS "github.com/cyber_bed/internal/notifications/delivery/web_socket"
	notificationsRepository "github.com/cyber_bed/internal/notifications/repository"
	notificationsUsecase "github.com/cyber_bed/internal/notifications/usecase"
	logger "github.com/cyber_bed/pkg"
)

func main() {
	var configPath string
	config.ParseFlag(&configPath)
	flag.Parse()

	cfg := config.New()
	if err := cfg.Open(configPath); err != nil {
		defaultLogger.Fatalf("notifications: failed to open config: %s", err)
	}

	log := logger.GetInstance()
	log.SetLevel(logger.ToLevel(cfg.LoggerLvl))
	log.Info("notifications: server started")

	r, err := notificationsRepository.NewPostgres(cfg.FormatDbAddr())
	if err != nil {
		log.Fatalf("notifications: failed to open db: %s", err)
	}
	// defer r.DB.Close()

	u := notificationsUsecase.NewNotificationsUsecase(r)

	h := notificationsWS.NewWebSocket(&websocket.Upgrader{
		HandshakeTimeout: time.Minute,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}, u)

	http.HandleFunc("/ws", h.Handler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
	if err = http.ListenAndServe(addr, nil); err != nil {
		log.Warnf("notifications: stop to listen and serve: %s", err)
	}
}

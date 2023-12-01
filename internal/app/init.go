package app

import (
	httpCORS "net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	httpAuth "github.com/cyber_bed/internal/auth/delivery"
	authMiddlewares "github.com/cyber_bed/internal/auth/delivery/middleware"
	authRepository "github.com/cyber_bed/internal/auth/repository"
	authUsecase "github.com/cyber_bed/internal/auth/usecase"
	"github.com/cyber_bed/internal/config"
	"github.com/cyber_bed/internal/domain"
	foldersHandler "github.com/cyber_bed/internal/folders/delivery"
	foldersRepository "github.com/cyber_bed/internal/folders/repository"
	foldersUsecase "github.com/cyber_bed/internal/folders/usecase"
	httpNotifications "github.com/cyber_bed/internal/notifications/delivery/http"
	notificationsRepository "github.com/cyber_bed/internal/notifications/repository"
	notificationsUsecase "github.com/cyber_bed/internal/notifications/usecase"
	httpPlants "github.com/cyber_bed/internal/plants/delivery"
	plantsRepository "github.com/cyber_bed/internal/plants/repository"
	plantsUsecase "github.com/cyber_bed/internal/plants/usecase"
	domainRecognition "github.com/cyber_bed/internal/recognize-api"
	"github.com/cyber_bed/internal/recognize-api/delivery/http"
	recAPI "github.com/cyber_bed/internal/recognize-api/repository/api"
	recUsecase "github.com/cyber_bed/internal/recognize-api/usecase"
	httpUsers "github.com/cyber_bed/internal/users/delivery"
	usersRepository "github.com/cyber_bed/internal/users/repository"
	usersUsecase "github.com/cyber_bed/internal/users/usecase"
	"github.com/cyber_bed/migrations"
	logger "github.com/cyber_bed/pkg"
)

type Server struct {
	Echo   *echo.Echo
	Config *config.Config

	usersUsecase         domain.UsersUsecase
	authUsecase          domain.AuthUsecase
	plantsUsecase        domain.PlantsUsecase
	recUsecase           domainRecognition.Usecase
	plantsAPI            domain.PlantsAPI
	foldersUsecase       domain.FoldersUsecase
	notificationsUsecase domain.NotificationsUsecase

	usersHandler         httpUsers.UsersHandler
	authHandler          httpAuth.AuthHandler
	recHandler           domainRecognition.Handler
	plantsHandler        httpPlants.PlantsHandler
	foldersHandler       foldersHandler.FoldersHandler
	notificationsHandler httpNotifications.NotificationsHandler

	authMiddleware *authMiddlewares.Middlewares
}

func New(e *echo.Echo, c *config.Config) *Server {
	return &Server{
		Echo:   e,
		Config: c,
	}
}

func (s *Server) init() {
	s.MakeUsecases()
	s.makeMiddlewares()
	s.MakeHandlers()
	s.MakeRouter()

	s.MakeEchoLogger()
	s.makeCORS()

	if s.Config.Database.InitDB.Init {
		if err := s.makeMigrations(s.Config.FormatDbAddr(), s.Config.Database.InitDB.PathToDir); err != nil {
			s.Echo.Logger.Error(err)
		}
	}

	if s.Config.TranslateMode {
		s.Echo.Logger.Info("starting translate mode")
		if err := migrations.TranslatePlants(s.Config.Database.InitDB.PathToDir, ""); err != nil {
			s.Echo.Logger.Error(err)
		}
	}
}

func (s *Server) Start() error {
	s.init()
	return s.Echo.Start(
		s.Config.Server.Address + ":" + strconv.FormatUint(s.Config.Server.Port, 10),
	)
}

func (s *Server) MakeHandlers() {
	s.recHandler = http.NewHandler(s.recUsecase)
	s.authHandler = httpAuth.NewAuthHandler(s.authUsecase, s.usersUsecase, s.Config.CookieSettings)
	s.plantsHandler = httpPlants.NewPlantsHandler(
		s.plantsUsecase,
		s.usersUsecase,
		s.plantsAPI,
		s.foldersUsecase,
	)
	s.foldersHandler = foldersHandler.NewFoldersHandler(
		s.foldersUsecase,
		s.usersUsecase,
		s.notificationsUsecase,
		s.plantsUsecase,
	)
	s.notificationsHandler = httpNotifications.NewNotificationsHandler(
		s.notificationsUsecase,
		s.usersUsecase,
	)
}

func (s *Server) MakeUsecases() {
	pgParams := s.Config.FormatDbAddr()

	authDB, err := authRepository.NewPostgres(pgParams)
	if err != nil {
		s.Echo.Logger.Error(err)
	}

	usersDB, err := usersRepository.NewPostgres(pgParams)
	if err != nil {
		s.Echo.Logger.Error(err)
	}

	plantsDB, err := plantsRepository.NewPostgres(pgParams)
	if err != nil {
		s.Echo.Logger.Error(err)
	}

	foldersDB, err := foldersRepository.NewPostgres(pgParams)
	if err != nil {
		s.Echo.Logger.Error(err)
	}

	u, err := url.Parse(s.Config.RecognizeAPI.BaseURL)
	if err != nil {
		s.Echo.Logger.Error(errors.Wrap(err, "failed to parse base url"))
		return
	}

	recognitionAPI := recAPI.NewRecognitionAPI(
		u,
		s.Config.RecognizeAPI.Token,
		s.Config.RecognizeAPI.ImageField,
		s.Config.RecognizeAPI.MaxImages,
		s.Config.RecognizeAPI.CountResults,
	)

	notificationsDB, err := notificationsRepository.NewPostgres(pgParams)
	if err != nil {
		s.Echo.Logger.Error(err)
	}

	s.authUsecase = authUsecase.NewAuthUsecase(authDB, usersDB, s.Config.CookieSettings)
	s.usersUsecase = usersUsecase.NewUsersUsecase(usersDB)
	s.plantsUsecase = plantsUsecase.NewPlansUsecase(plantsDB, s.plantsAPI)
	s.foldersUsecase = foldersUsecase.NewFoldersUsecase(foldersDB, plantsDB)
	s.recUsecase = recUsecase.New(recognitionAPI, s.plantsAPI, s.plantsUsecase)
	s.notificationsUsecase = notificationsUsecase.NewNotificationsUsecase(notificationsDB)
}

func (s *Server) MakeRouter() {
	v1 := s.Echo.Group("/api")
	v1.Use(logger.Middleware())
	v1.Use(middleware.Secure())

	v1.POST("/signup", s.authHandler.SignUp)
	v1.POST("/login", s.authHandler.Login)
	v1.GET("/auth", s.authHandler.Auth)
	v1.DELETE("/logout", s.authHandler.Logout, s.authMiddleware.LoginRequired)

	v1.POST("/recognize", s.recHandler.Recognize)

	plantsAPI := v1.Group("/search")
	plantsAPI.GET("/plants/:plantID", s.plantsHandler.GetPlantFromAPI)
	plantsAPI.GET("/plants/:plantID/image", s.plantsHandler.GetPlantImage)
	plantsAPI.GET("/plants", s.plantsHandler.GetPlantsFromAPI)

	plants := v1.Group("/plants", s.authMiddleware.LoginRequired)
	plants.GET("/:plantID", s.plantsHandler.GetPlant)
	plants.POST("/:plantID", s.plantsHandler.CreatePlant)
	plants.DELETE("/:plantID", s.plantsHandler.DeletePlant)
	plants.GET("", s.plantsHandler.GetPlants)
	plants.POST("/:plantID/saved", s.plantsHandler.CreateSavedPlant)
	plants.GET("/saved", s.plantsHandler.GetSavedPlants)
	plants.DELETE("/:plantID/saved", s.plantsHandler.DeleteSavedPlant)

	folders := v1.Group("/folders", s.authMiddleware.LoginRequired)
	folders.POST("", s.foldersHandler.CreateFolder)
	folders.GET("", s.foldersHandler.GetFolders)
	folders.GET("/:folderID/plants", s.foldersHandler.GetPlantsFromFolder)
	folders.DELETE("/:folderID", s.foldersHandler.DeleteFolder)
	folders.POST("/:folderID/plants/:plantID", s.foldersHandler.AddPlantToFolder)
	folders.DELETE("/:folderID/plants/:plantID", s.foldersHandler.DeletePlantFromFolder)
	folders.PUT("/:folderID/plants/:plantID", s.foldersHandler.UpdatePeriod)

	folders.POST("/:folderID/plants/:plantID/chan/:channelID", s.foldersHandler.CreateChannel)
	folders.GET("/:folderID/plants/:plantID/chan", s.foldersHandler.GetChannel)

	customPlants := v1.Group("/custom", s.authMiddleware.LoginRequired)
	customPlants.POST("/plants", s.plantsHandler.CreateCustomPlant)
	customPlants.PUT("/plants/:plantID", s.plantsHandler.UpdateCustomPlant)
	customPlants.GET("/plants", s.plantsHandler.GetCustomPlants)
	customPlants.GET("/plants/:plantID", s.plantsHandler.GetCustomPlant)
	customPlants.GET("/plants/:plantID/image", s.plantsHandler.GetCustomPlantImage)
	customPlants.DELETE("/plants/:plantID", s.plantsHandler.DeleteCustomPlant)

	notifications := v1.Group("/notifications", s.authMiddleware.LoginRequired)
	notifications.GET("", s.notificationsHandler.GetNotifications)
	notifications.DELETE("/:notificationID", s.notificationsHandler.DeleteNotification)
	notifications.DELETE("", s.notificationsHandler.DeleteCategoryNotification)
}

func (s *Server) makeMiddlewares() {
	s.authMiddleware = authMiddlewares.New(s.authUsecase, s.usersUsecase)
}

func (s *Server) makeMigrations(url, pathToDir string) error {
	if err := migrations.StartMigration(url, pathToDir); err != nil {
		return err
	}
	return nil
}

func (s *Server) MakeEchoLogger() {
	s.Echo.Logger = logger.GetInstance()
	s.Echo.Logger.SetLevel(logger.ToLevel(s.Config.LoggerLvl))
	s.Echo.Logger.Info("server started")
}

func (s *Server) makeCORS() {
	s.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			httpCORS.MethodGet,
			httpCORS.MethodHead,
			httpCORS.MethodPut,
			httpCORS.MethodPatch,
			httpCORS.MethodPost,
			httpCORS.MethodDelete,
		},
		AllowCredentials: true,
	}))
}

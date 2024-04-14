package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"avito/internal/handler"
	"avito/internal/middleware"
	"avito/internal/pkg"
	"avito/internal/presenter"
	"avito/internal/repository"
	"avito/internal/service"
	"avito/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func setup(r *gin.Engine) {
	loadEnv()

	// Wire dependencies
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	db := pkg.NewDB(os.Getenv("APP_DB_DSN"))
	cache := pkg.NewCache(os.Getenv("APP_CACHE_ADDR"))
	userRepository := repository.NewUser(db)
	bannerRepository := repository.NewBanner(db)
	bannerTagRepository := repository.NewBannerTag(db)
	bannerHistoryRepository := repository.NewBannerHistory(db)
	deleteBannersJobRepository := repository.NewDeleteBannersJob(db)
	bannerPresenter := presenter.NewBanner()
	bannerHistoryPresenter := presenter.NewBannerHistory()
	cacheService := service.NewCache(cache)
	authService := service.NewAuth(cacheService, userRepository, os.Getenv("APP_SECRET_KEY"))
	bannerService := service.NewBanner(
		db,
		cacheService,
		bannerRepository,
		bannerTagRepository,
		bannerHistoryRepository,
		deleteBannersJobRepository,
	)
	bannerHistoryService := service.NewBannerHistory(db, bannerRepository, bannerHistoryRepository)
	validate := validator.New(validator.WithRequiredStructEnabled())
	authHandler := handler.NewAuth(authService)
	userBannerHandler := handler.NewUserBanner(bannerService, bannerPresenter)
	bannerHandler := handler.NewBanner(validate, bannerService, bannerPresenter)
	bannerHistoryHandler := handler.NewBannerHistory(bannerHistoryService, bannerHistoryPresenter)
	authenticateMiddleware := middleware.NewAuthenticate(authService)
	isAdminMiddleware := middleware.NewIsAdmin()
	deleteBannersWorker := worker.NewDeleteBanners(db, logger, deleteBannersJobRepository, bannerRepository)

	// setup router
	r.POST("/auth", authHandler.Auth)
	r.GET("/user_banner", authenticateMiddleware.Handle, userBannerHandler.Get)
	r.GET("/banner", authenticateMiddleware.Handle, isAdminMiddleware.Handle, bannerHandler.Get)
	r.POST("/banner", authenticateMiddleware.Handle, isAdminMiddleware.Handle, bannerHandler.Create)
	r.PATCH("/banner/:id", authenticateMiddleware.Handle, isAdminMiddleware.Handle, bannerHandler.Update)
	r.DELETE("/banner/:id", authenticateMiddleware.Handle, isAdminMiddleware.Handle, bannerHandler.Delete)

	r.GET("/banner_history", authenticateMiddleware.Handle, isAdminMiddleware.Handle, bannerHistoryHandler.Get)
	r.POST("/banner_history", authenticateMiddleware.Handle, isAdminMiddleware.Handle, bannerHistoryHandler.Apply)

	r.DELETE("/banner", authenticateMiddleware.Handle, isAdminMiddleware.Handle, bannerHandler.DeleteMany)

	// start workers
	go deleteBannersWorker.Run(context.Background())
}

func loadEnv() {
	envRoot := os.Getenv("ENV_ROOT")
	if envRoot == "" {
		envRoot = "./"
	}

	if env := os.Getenv("ENV"); env != "" {
		if err := godotenv.Load(envRoot + ".env." + env); err != nil {
			log.Fatal(err)
		}
	}

	if err := godotenv.Load(envRoot + ".env"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	setup(r)

	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		log.Fatal("Addr not provided")
	}

	log.Fatal(r.Run(addr))
}

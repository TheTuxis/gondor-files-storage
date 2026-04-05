package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TheTuxis/gondor-files-storage/internal/config"
	"github.com/TheTuxis/gondor-files-storage/internal/handler"
	"github.com/TheTuxis/gondor-files-storage/internal/middleware"
	"github.com/TheTuxis/gondor-files-storage/internal/model"
	"github.com/TheTuxis/gondor-files-storage/internal/repository"
	"github.com/TheTuxis/gondor-files-storage/internal/service"
)

func main() {
	cfg := config.Load()

	var logger *zap.Logger
	var err error
	if cfg.Environment == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}

	if err := db.AutoMigrate(
		&model.File{},
		&model.Folder{},
		&model.FileVersion{},
		&model.FileShare{},
	); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	storageSvc, err := service.NewStorageService(cfg)
	if err != nil {
		logger.Fatal("failed to initialize storage service", zap.Error(err))
	}

	fileRepo := repository.NewFileRepository(db)
	folderRepo := repository.NewFolderRepository(db)

	fileSvc := service.NewFileService(fileRepo, storageSvc)
	folderSvc := service.NewFolderService(folderRepo)

	healthHandler := handler.NewHealthHandler()
	fileHandler := handler.NewFileHandler(fileSvc, logger)
	folderHandler := handler.NewFolderHandler(folderSvc, logger)

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.AuthMiddleware(cfg.JWTSecret, logger))

	r.MaxMultipartMemory = cfg.MaxUploadSize

	r.GET("/health", healthHandler.Health)
	r.GET("/metrics", healthHandler.Metrics)

	v1 := r.Group("/v1")
	{
		files := v1.Group("/files")
		{
			files.POST("/upload", fileHandler.Upload)
			files.GET("", fileHandler.List)
			files.GET("/:id", fileHandler.GetByID)
			files.GET("/:id/download", fileHandler.Download)
			files.DELETE("/:id", fileHandler.Delete)

			folders := files.Group("/folders")
			{
				folders.POST("", folderHandler.Create)
				folders.GET("", folderHandler.List)
				folders.GET("/:id", folderHandler.GetByID)
				folders.PUT("/:id", folderHandler.Update)
				folders.DELETE("/:id", folderHandler.Delete)
			}
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited")
}

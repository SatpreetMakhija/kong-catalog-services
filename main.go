package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"kong.com/catalog/datastore"
	"kong.com/catalog/handlers"
	"kong.com/catalog/models"
)

var port = "8080"

func main() {

	ctx := context.Background()
	router := gin.Default()

	dbCfg := loadDbConfig()

	db := datastore.PostgresDatastore{}
	err := db.InitConnection(ctx, &dbCfg)

	if err != nil {
		slog.Error("failed to establish connection with db: %w", err)
		os.Exit(1)
	}
	serviceRepo := datastore.NewPostgresServiceRepo(&db)
	serviceModel := models.NewService(serviceRepo)
	serviceHandler := handlers.NewServiceHandler(serviceModel)

	api := router.Group("/api/v1")
	{
		api.GET("services/:id", serviceHandler.GetServiceById)
	}

	router.Run(fmt.Sprintf(":%v", port))
}

func loadDbConfig() datastore.DBConfig {
	return datastore.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Username: os.Getenv("DB_USERNAME"),
	}
}

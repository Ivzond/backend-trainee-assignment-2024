package main

import (
	"log"
	"test_repository/internal/api"
	"test_repository/internal/config"
	"test_repository/internal/database/postgres"
	"test_repository/internal/database/redis"
	"test_repository/internal/repository"
	"test_repository/internal/service"
	"test_repository/pkg/helpers"
)

func main() {
	helpers.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewRedisClient(cfg)

	bannerRepo := repository.NewBannerRepository(db, redisClient)
	bannerService := service.NewBannerService(*bannerRepo, redisClient)

	router := api.SetupRouter(*bannerService)
	log.Println("App is working on port :8080")
	router.Run(":8080")
}

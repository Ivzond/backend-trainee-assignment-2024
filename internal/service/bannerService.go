package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"test_repository/internal/models"
	"test_repository/internal/repository"
	"test_repository/pkg/helpers"
)

type BannerService struct {
	Repo        repository.BannerRepository
	RedisClient *redis.Client
}

func NewBannerService(repo repository.BannerRepository, redisClient *redis.Client) *BannerService {
	return &BannerService{
		Repo:        repo,
		RedisClient: redisClient,
	}
}

func (s *BannerService) GetBannersByTagAndFeature(tagID, featureID int) ([]*models.Banner, error) {
	return s.Repo.GetBannersByTagAndFeature(tagID, featureID)
}

func (s *BannerService) GetBannersFromRedis(tagID, featureID int) ([]*models.Banner, error) {
	cacheKey := fmt.Sprintf("banners:%d:%d", tagID, featureID)
	val, err := s.RedisClient.Get(cacheKey).Bytes()
	if err != nil {
		return nil, err
	}

	var banners []*models.Banner
	err = json.Unmarshal(val, &banners)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (s *BannerService) CreateBanner(banner *models.Banner) (int, error) {
	helpers.InfoLogger.Printf("Banner body in service layer: %v", banner)
	return s.Repo.CreateBanner(banner)
}

func (s *BannerService) UpdateBanner(banner *models.Banner) error {
	return s.Repo.UpdateBanner(banner)
}

func (s *BannerService) DeleteBannerByID(id int) error {
	return s.Repo.DeleteBannerByID(id)
}

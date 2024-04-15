package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"test_repository/internal/config"
	"test_repository/internal/database/postgres"
	"test_repository/internal/database/redis"
	"test_repository/internal/repository"
	"test_repository/internal/service"
	"test_repository/pkg/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
	"test_repository/internal/api"
	"test_repository/internal/models"
)

func TestIntegration_GetUserBanner(t *testing.T) {
	router := setupRouter()
	banner := `{
		"tag_ids": [1, 2],
		"feature_id": 1,
		"content": {
			"title": "New Banner",
			"text": "Description",
			"url": "https://avito.ru"
		},
		"is_active": true
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/banner", strings.NewReader(banner))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "admin_token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	w = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/v1/user_banner?tag_id=1&feature_id=1&use_last_revision=false", nil)
	req.Header.Set("token", "user_token")
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Banner
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestIntegration_GetBanners(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/banner?tag_id=1&feature_id=1&limit=10&offset=0", nil)
	req.Header.Set("token", "admin_token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*models.Banner
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestIntegration_CreateBanner(t *testing.T) {
	router := setupRouter()

	banner := `{
		"tag_ids": [1, 2],
		"feature_id": 1,
		"content": {
			"title": "New Banner",
			"text": "Description",
			"url": "https://avito.ru"
		},
		"is_active": true
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/banner", strings.NewReader(banner))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "admin_token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response["banner_id"])
}

func TestIntegration_UpdateBanner(t *testing.T) {
	router := setupRouter()

	banner := `{
		"tag_ids": [1],
		"feature_id": 2,
		"content": {
			"title": "Updated Banner",
			"text": "Some updated text",
			"url": "https://updated-example.com"
		},
		"is_active": false
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/banner/1", strings.NewReader(banner))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("token", "admin_token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_DeleteBanner(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/v1/banner/1", nil)
	req.Header.Set("token", "admin_token")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func setupRouter() http.Handler {
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
	return router
}

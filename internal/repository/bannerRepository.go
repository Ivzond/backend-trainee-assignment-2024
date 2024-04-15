package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lib/pq"
	"test_repository/internal/models"
	"test_repository/pkg/helpers"
	"time"
)

type BannerRepository struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewBannerRepository(db *sql.DB, redis *redis.Client) *BannerRepository {
	return &BannerRepository{
		DB:    db,
		Redis: redis,
	}
}

func (r *BannerRepository) GetBannersByTagAndFeature(tagID, featureID int) ([]*models.Banner, error) {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("banners:tag:%d:feature:%d", tagID, featureID)
	val, err := r.Redis.Get(cacheKey).Bytes()
	if err == nil {
		var banners []*models.Banner
		err := json.Unmarshal(val, &banners)
		if err != nil {
			return nil, err
		}
		return banners, nil
	}

	query := `SELECT id, tag_ids, feature_id, content, is_active, created_at, updated_at FROM banners 
				WHERE $1 = ANY(tag_ids) AND feature_id = $2 AND is_active = true`

	rows, err := r.DB.QueryContext(ctx, query, tagID, featureID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banners []*models.Banner
	for rows.Next() {
		banner := &models.Banner{}
		var contentJSON []byte
		err := rows.Scan(&banner.ID, (*pq.Int64Array)(&banner.TagIDs), &banner.FeatureID, &contentJSON, &banner.IsActive, &banner.CreatedAt, &banner.UpdatedAt)
		if err != nil {
			return nil, err
		}
		banner.Content = make(map[string]interface{})
		if err := json.Unmarshal(contentJSON, &banner.Content); err != nil {
			return nil, err
		}
		banners = append(banners, banner)
	}

	cacheValue, err := json.Marshal(banners)
	if err != nil {
		return nil, err
	}
	err = r.Redis.Set(cacheKey, cacheValue, 5*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (r *BannerRepository) CreateBanner(banner *models.Banner) (int, error) {
	ctx := context.Background()

	query := `INSERT INTO banners (tag_ids, feature_id, content, is_active, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	createdAt := time.Now()
	updatedAt := time.Now()

	tagIDsArray := pq.Array(banner.TagIDs)
	contentBytes, err := json.Marshal(banner.Content)
	if err != nil {
		return 0, err
	}
	helpers.InfoLogger.Printf("Banner body in repo layer: %v", banner)

	err = r.DB.QueryRowContext(ctx, query, tagIDsArray, banner.FeatureID, contentBytes, banner.IsActive, createdAt, updatedAt).Scan(&banner.ID)
	if err != nil {
		return 0, err
	}

	return banner.ID, nil
}

func (r *BannerRepository) UpdateBanner(banner *models.Banner) error {
	ctx := context.Background()

	query := `UPDATE banners SET tag_ids = $1, feature_id = $2, content = $3, is_active = $4, updated_at = $5 
				WHERE id = $6`

	updatedAt := time.Now()

	tagIDsArray := pq.Array(banner.TagIDs)
	contentBytes, err := json.Marshal(banner.Content)
	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx, query, tagIDsArray, banner.FeatureID, contentBytes, banner.IsActive, updatedAt, banner.ID)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("banners:tag:%d:feature:%d", banner.TagIDs, banner.FeatureID)
	err = r.Redis.Del(cacheKey).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *BannerRepository) DeleteBannerByID(id int) error {
	ctx := context.Background()

	query := "DELETE FROM banners WHERE id = $1"
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("banners:id:%d", id)
	err = r.Redis.Del(cacheKey).Err()
	if err != nil {
		return err
	}

	return nil
}

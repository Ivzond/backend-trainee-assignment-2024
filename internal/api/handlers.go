package api

import (
	"net/http"
	"strconv"
	"test_repository/internal/models"
	"test_repository/internal/service"
	"test_repository/pkg/helpers"

	"github.com/gin-gonic/gin"
)

func GetUserBanner(service service.BannerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID, _ := strconv.Atoi(c.Query("tag_id"))
		featureID, _ := strconv.Atoi(c.Query("feature_id"))
		useLastRevision, _ := strconv.ParseBool(c.Query("use_last_revision"))

		var banners []*models.Banner
		var err error
		if useLastRevision {
			banners, err = service.GetBannersFromRedis(tagID, featureID)
		} else {
			banners, err = service.GetBannersByTagAndFeature(tagID, featureID)
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(banners) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
			return
		}

		c.JSON(http.StatusOK, banners[0])
	}
}

func GetBanners(service service.BannerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tagID, _ := strconv.Atoi(c.Query("tag_id"))
		featureID, _ := strconv.Atoi(c.Query("feature_id"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))  // Default limit to 10 if not provided
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0")) // Default offset to 0 if not provided

		banners, err := service.GetBannersByTagAndFeature(tagID, featureID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if offset > len(banners) {
			banners = []*models.Banner{}
		} else {
			banners = banners[offset:]
			if limit < len(banners) {
				banners = banners[:limit]
			}
		}

		c.JSON(http.StatusOK, banners)
	}
}

func CreateBanner(service service.BannerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var banner models.Banner
		if err := c.BindJSON(&banner); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		helpers.InfoLogger.Printf("Banner body in handlers layer: %v", banner)
		id, err := service.CreateBanner(&banner)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"banner_id": id})
	}
}

func UpdateBanner(service service.BannerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var banner models.Banner
		if err := c.BindJSON(&banner); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		banner.ID = id

		err := service.UpdateBanner(&banner)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func DeleteBanner(service service.BannerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		err := service.DeleteBannerByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNoContent, gin.H{})
	}
}

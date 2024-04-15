package api

import (
	"github.com/gin-gonic/gin"
	"test_repository/internal/service"
)

func SetupRouter(service service.BannerService) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.Use(AdminAuthMiddleware())
		bannerGroup := v1.Group("/banner")
		{
			bannerGroup.GET("", GetBanners(service))
			bannerGroup.POST("", CreateBanner(service))
			bannerGroup.PATCH("/:id", UpdateBanner(service))
			bannerGroup.DELETE("/:id", DeleteBanner(service))
		}
	}
	router.GET("/v1/user_banner", UserAuthMiddleware(), GetUserBanner(service))

	return router
}

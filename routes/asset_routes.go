package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/controllers"
	"gorm.io/gorm"
)

func SetupAssetRoutes(router *gin.Engine, db *gorm.DB) {
	assetController := controllers.NewAssetController(db)

	assetRoutes := router.Group("/v1/assets")
	{
		assetRoutes.GET("/", assetController.ListAssets)
		assetRoutes.POST("/", assetController.CreateAsset)
		assetRoutes.GET("/:id", assetController.GetAsset)
		assetRoutes.PUT("/:id", assetController.UpdateAsset)
		assetRoutes.DELETE("/:id", assetController.DeleteAsset)
		assetRoutes.PATCH("/:id/favorite", assetController.ToggleFavorite)

		assetRoutes.GET("/cameras", assetController.ListCameras)
		assetRoutes.GET("/cameras2", assetController.ListCamerasWithImages)
	}
}

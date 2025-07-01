package routes

import (
	"context"
	"errors"
	"github.com/mahdi-cpp/PhotoKit/repositories"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	repo repositories.AssetRepository
}

func NewAssetHandler(repo repositories.AssetRepository) *AssetHandler {
	return &AssetHandler{repo: repo}
}

func (h *AssetHandler) GetAllAssets(c *gin.Context) {

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	assets, err := h.repo.GetAllAssets(ctx)
	if err != nil {
		log.Printf("Error fetching assets: %v", err)

		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timeout",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve assets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  assets,
		"count": len(assets),
	})
}

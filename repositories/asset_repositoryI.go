package repositories

import (
	"context"
	"github.com/mahdi-cpp/PhotoKit/models"
)

type AssetRepository interface {
	GetAllAssets(ctx context.Context) ([]models.PHAsset, error)
	GetAssetByID(ctx context.Context, id int) (*models.PHAsset, error)
}

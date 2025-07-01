package repositories

import (
	"context"
)

type AssetRepository interface {
	GetAllAssets(ctx context.Context) ([]PHAsset, error)
	GetAssetByID(ctx context.Context, id int) (*PHAsset, error)
}

package repositories

import (
	"context"
	"github.com/lib/pq"
	"github.com/mahdi-cpp/PhotoKit/models"
	"gorm.io/gorm"
	"log"
	"time"
)

var db *gorm.DB

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	// Auto migrate the User models
	err := db.AutoMigrate(&models.PHAsset{})
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{db: db}
}

func (r *Repository) GetAllAssets(ctx context.Context) ([]models.PHAsset, error) {
	var assets []models.PHAsset
	result := r.db.WithContext(ctx).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetAssetByID fetches a single PHAsset by its ID
func (r *Repository) GetAssetByID(ctx context.Context, id int) (*models.PHAsset, error) {
	var asset models.PHAsset
	result := db.WithContext(ctx).First(&asset, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &asset, nil
}

// GetAssetsByMediaType fetches assets filtered by media type
func GetAssetsByMediaType(mediaType string) ([]models.PHAsset, error) {
	var assets []models.PHAsset
	result := db.Where("media_type = ?", mediaType).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetFavoriteAssets fetches all favorite assets
func GetFavoriteAssets() ([]models.PHAsset, error) {
	var assets []models.PHAsset
	result := db.Where("is_favorite = ?", true).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetAssetsInAlbums fetches assets that belong to specific albums
func GetAssetsInAlbums(albumIDs []int32) ([]models.PHAsset, error) {
	var assets []models.PHAsset
	result := db.Where("albums @> ?", pq.Int32Array(albumIDs)).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetRecentAssets fetches assets created in the last N days
func GetRecentAssets(days int) ([]models.PHAsset, error) {
	var assets []models.PHAsset
	since := time.Now().AddDate(0, 0, -days)
	result := db.Where("creation_date > ?", since).Order("creation_date desc").Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

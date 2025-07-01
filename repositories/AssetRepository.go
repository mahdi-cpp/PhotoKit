package repositories

import (
	"context"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"time"
)

type PHAsset struct {
	ID int `gorm:"primaryKey;autoIncrement" json:"id"`

	// Media Characteristics
	Url   string `json:"url"`
	Named string `gorm:"index" json:"named"`

	MediaType string `json:"mediaType"`
	Format    string `json:"format"`

	Orientation int `json:"orientation"`

	PixelWidth  int `json:"pixelWidth"`
	PixelHeight int `json:"pixelHeight"`

	IsFavorite bool `gorm:"default:false" json:"isFavorite"`
	IsHidden   bool `gorm:"default:false" json:"isHidden"`

	// Video Properties
	Duration float64 `gorm:"default:0" json:"duration"`

	// Content Availability
	CanDelete           bool `gorm:"default:true" json:"canDelete"`
	CanEditContent      bool `gorm:"default:true" json:"canEditContent"`
	CanAddToSharedAlbum bool `gorm:"default:true" json:"CanAddToSharedAlbum"`

	// Advanced Properties
	IsUserLibraryAsset bool `gorm:"default:true" json:"IsUserLibraryAsset"`

	Albums  pq.Int32Array `gorm:"type:integer[]" json:"albums"`
	Trips   pq.Int32Array `gorm:"type:integer[]" json:"trips"`
	Persons pq.Int32Array `gorm:"type:integer[]" json:"persons"`
	Cameras pq.Int32Array `gorm:"type:integer[]" json:"cameras"`

	// Dates
	CreationDate     time.Time `gorm:"default:now()" json:"createdAt"`
	ModificationDate time.Time `json:"ModificationDate"`
}

var db *gorm.DB

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	// Auto migrate the User models
	err := db.AutoMigrate(&PHAsset{})
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{db: db}
}

func (r *Repository) GetAllAssets(ctx context.Context) ([]PHAsset, error) {
	var assets []PHAsset
	result := r.db.WithContext(ctx).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetAssetByID fetches a single PHAsset by its ID
func (r *Repository) GetAssetByID(ctx context.Context, id int) (*PHAsset, error) {
	var asset PHAsset
	result := db.WithContext(ctx).First(&asset, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &asset, nil
}

// GetAssetsByMediaType fetches assets filtered by media type
func GetAssetsByMediaType(mediaType string) ([]PHAsset, error) {
	var assets []PHAsset
	result := db.Where("media_type = ?", mediaType).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetFavoriteAssets fetches all favorite assets
func GetFavoriteAssets() ([]PHAsset, error) {
	var assets []PHAsset
	result := db.Where("is_favorite = ?", true).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetAssetsInAlbums fetches assets that belong to specific albums
func GetAssetsInAlbums(albumIDs []int32) ([]PHAsset, error) {
	var assets []PHAsset
	result := db.Where("albums @> ?", pq.Int32Array(albumIDs)).Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

// GetRecentAssets fetches assets created in the last N days
func GetRecentAssets(days int) ([]PHAsset, error) {
	var assets []PHAsset
	since := time.Now().AddDate(0, 0, -days)
	result := db.Where("creation_date > ?", since).Order("creation_date desc").Find(&assets)
	if result.Error != nil {
		return nil, result.Error
	}
	return assets, nil
}

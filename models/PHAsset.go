package models

import (
	"github.com/lib/pq"
	"time"
)

type PHAsset struct {
	ID     int `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId int `gorm:"references:users(id);onDelete:SET NULL" json:"userId"`

	// Media Characteristics
	URL   string `json:"url"`
	Named string `json:"named"`

	MediaType   string `json:"mediaType"`
	Format      string `json:"format"`
	Orientation int    `json:"orientation"`

	PixelWidth  int `json:"pixelWidth"`
	PixelHeight int `json:"pixelHeight"`

	CameraMake  string `gorm:"default:NULL" json:"CameraMake"`
	CameraModel string `gorm:"default:NULL" json:"CameraModel"`

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

	ModificationDate time.Time `gorm:"type:timestamp;default:NULL"`
	CreationDate     time.Time `gorm:"type:timestamp;not null"`
}

package controllers

import (
	"github.com/mahdi-cpp/PhotoKit/models"
	"github.com/mahdi-cpp/PhotoKit/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type AssetController struct {
	db *gorm.DB
}

func NewAssetController(db *gorm.DB) *AssetController {
	return &AssetController{db: db}
}

// CreateAssetRequest defines the payload for creating a new asset
type CreateAssetRequest struct {
	URL        string  `json:"url" binding:"required"`
	Named      string  `json:"named"`
	MediaType  string  `json:"mediaType" binding:"required"`
	Format     string  `json:"format"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Duration   float64 `json:"duration"`
	IsFavorite bool    `json:"isFavorite"`
	Albums     []int32 `json:"albums"`
	Trips      []int32 `json:"trips"`
	Persons    []int32 `json:"persons"`
	Cameras    []int32 `json:"cameras"`
}

// UpdateAssetRequest defines the payload for updating an asset
type UpdateAssetRequest struct {
	Named      *string `json:"named"`
	IsFavorite *bool   `json:"isFavorite"`
	IsHidden   *bool   `json:"isHidden"`
	Albums     []int32 `json:"albums"`
	Trips      []int32 `json:"trips"`
	Persons    []int32 `json:"persons"`
	Cameras    []int32 `json:"cameras"`
}

// CreateAsset godoc
// @Summary Create a new asset
// @Description Create a new photo/video asset
// @Tags assets
// @Accept  json
// @Produce  json
// @Param asset body CreateAssetRequest true "Asset data"
// @Success 201 {object} models.PHAsset
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /assets [post]
func (ac *AssetController) CreateAsset(c *gin.Context) {
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	asset := models.PHAsset{
		URL:          req.URL,
		MediaType:    req.MediaType,
		Format:       req.Format,
		PixelWidth:   req.Width,
		PixelHeight:  req.Height,
		Duration:     req.Duration,
		IsFavorite:   req.IsFavorite,
		Albums:       pq.Int32Array(req.Albums),
		Trips:        pq.Int32Array(req.Trips),
		Persons:      pq.Int32Array(req.Persons),
		CreationDate: time.Now(),
	}

	result := ac.db.Create(&asset)
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create asset")
		return
	}

	utils.SendSuccess(c, http.StatusCreated, asset)
}

// GetAsset godoc
// @Summary Get an asset by ID
// @Description Get a single asset by its ID
// @Tags assets
// @Accept  json
// @Produce  json
// @Param id path int true "Asset ID"
// @Success 200 {object} models.PHAsset
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /assets/{id} [get]
func (ac *AssetController) GetAsset(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var asset models.PHAsset
	result := ac.db.First(&asset, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Asset not found")
		} else {
			utils.SendError(c, http.StatusInternalServerError, "Failed to fetch asset")
		}
		return
	}

	utils.SendSuccess(c, http.StatusOK, asset)
}

// UpdateAsset godoc
// @Summary Update an asset
// @Description Update an existing asset
// @Tags assets
// @Accept  json
// @Produce  json
// @Param id path int true "Asset ID"
// @Param asset body UpdateAssetRequest true "Asset update data"
// @Success 200 {object} models.PHAsset
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /assets/{id} [put]
func (ac *AssetController) UpdateAsset(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var req UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var asset models.PHAsset
	result := ac.db.First(&asset, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Asset not found")
		} else {
			utils.SendError(c, http.StatusInternalServerError, "Failed to fetch asset")
		}
		return
	}

	// Update fields if they are provided in the request
	//if req.Named != nil {
	//	asset.Named = *req.Named
	//}
	if req.IsFavorite != nil {
		asset.IsFavorite = *req.IsFavorite
	}
	if req.IsHidden != nil {
		asset.IsHidden = *req.IsHidden
	}
	if req.Albums != nil {
		asset.Albums = pq.Int32Array(req.Albums)
	}
	if req.Trips != nil {
		asset.Trips = pq.Int32Array(req.Trips)
	}
	if req.Persons != nil {
		asset.Persons = pq.Int32Array(req.Persons)
	}

	asset.ModificationDate = time.Now()

	result = ac.db.Save(&asset)
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	utils.SendSuccess(c, http.StatusOK, asset)
}

// DeleteAsset godoc
// @Summary Delete an asset
// @Description Delete an asset by ID
// @Tags assets
// @Accept  json
// @Produce  json
// @Param id path int true "Asset ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /assets/{id} [delete]
func (ac *AssetController) DeleteAsset(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	result := ac.db.Delete(&models.PHAsset{}, id)
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to delete asset")
		return
	}
	if result.RowsAffected == 0 {
		utils.SendError(c, http.StatusNotFound, "Asset not found")
		return
	}

	c.Status(http.StatusNoContent)
}

// ListAssets godoc
// @Summary List all assets
// @Description Get a list of all assets with optional filtering
// @Tags assets
// @Accept  json
// @Produce  json
// @Param mediaType query string false "Filter by media type"
// @Param favorite query bool false "Filter by favorite status"
// @Param recentDays query int false "Filter by recent days"
// @Param album query int false "Filter by album ID"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset results"
// @Success 200 {array} models.PHAsset
// @Failure 500 {object} utils.ErrorResponse
// @Router /assets [get]
func (ac *AssetController) ListAssets(c *gin.Context) {
	query := ac.db.Model(&models.PHAsset{})

	// Apply filters
	if mediaType := c.Query("mediaType"); mediaType != "" {
		query = query.Where("media_type = ?", mediaType)
	}
	if favorite := c.Query("favorite"); favorite != "" {
		isFavorite, _ := strconv.ParseBool(favorite)
		query = query.Where("is_favorite = ?", isFavorite)
	}
	if recentDays := c.Query("recentDays"); recentDays != "" {
		days, _ := strconv.Atoi(recentDays)
		since := time.Now().AddDate(0, 0, -days)
		query = query.Where("creation_date > ?", since).Order("creation_date desc")
	}
	if albumID := c.Query("album"); albumID != "" {
		albumIDInt, _ := strconv.Atoi(albumID)
		query = query.Where("albums @> ?", pq.Int32Array{int32(albumIDInt)})
	}

	if cameraId := c.Query("cameras"); cameraId != "" {
		cameraIdInt, _ := strconv.Atoi(cameraId)
		query = query.Where("cameras @> ?", pq.Int32Array{int32(cameraIdInt)})
	}

	// Apply pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	query = query.Limit(limit).Offset(offset)

	var assets []models.PHAsset
	result := query.Find(&assets)
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch assets")
		return
	}

	utils.SendSuccess(c, http.StatusOK, assets)
}

// ToggleFavorite godoc
// @Summary Toggle favorite status
// @Description Toggle an asset's favorite status
// @Tags assets
// @Accept  json
// @Produce  json
// @Param id path int true "Asset ID"
// @Success 200 {object} models.PHAsset
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /assets/{id}/favorite [patch]
func (ac *AssetController) ToggleFavorite(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var asset models.PHAsset
	result := ac.db.First(&asset, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			utils.SendError(c, http.StatusNotFound, "Asset not found")
		} else {
			utils.SendError(c, http.StatusInternalServerError, "Failed to fetch asset")
		}
		return
	}

	asset.IsFavorite = !asset.IsFavorite
	asset.ModificationDate = time.Now()

	result = ac.db.Save(&asset)
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to update asset")
		return
	}

	utils.SendSuccess(c, http.StatusOK, asset)
}

// CameraResponse defines the JSON response structure
type CameraResponse struct {
	Make       string `json:"make"`
	Model      string `json:"model"`
	AssetCount int    `json:"assetCount"`
}

// ListCameras godoc
// curl -X GET "http://localhost:8095/assets/cameras" -H "Accept: application/json"
// curl -X GET "http://localhost:8095/assets/cameras?make=samsung" -H "Accept: application/json"
// @Summary List all camera models and manufacturers
// @Description Get a list of all camera models with optional filtering by manufacturer, model, or usage count
// @Tags cameras
// @Accept  json
// @Produce  json
// @Param make query string false "Filter by manufacturer (e.g., 'Sony', 'Canon')"
// @Param model query string false "Filter by model name (e.g., 'A7III', 'iPhone 14 Pro')"
// @Param minAssets query int false "Minimum number of assets for a camera to be included"
// @Param limit query int false "Limit results (default: 20)"
// @Param offset query int false "Offset results (default: 0)"
// @Success 200 {array} CameraResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /cameras [get]
func (ac *AssetController) ListCameras(c *gin.Context) {

	query := ac.db.Model(&models.PHAsset{}).
		Select("camera_make as make, camera_model as model, COUNT(*) as asset_count").
		Where("camera_make IS NOT NULL AND camera_model IS NOT NULL").
		Group("camera_make, camera_model")

	// Apply filters
	if make := c.Query("make"); make != "" {
		query = query.Where("camera_make = ?", make)
	}
	if model := c.Query("model"); model != "" {
		query = query.Where("camera_model = ?", model)
	}
	if minAssets := c.Query("minAssets"); minAssets != "" {
		min, _ := strconv.Atoi(minAssets)
		query = query.Having("COUNT(*) >= ?", min)
	}

	// Apply sorting (default: most-used first)
	query = query.Order("asset_count DESC")

	// Pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	query = query.Limit(limit).Offset(offset)

	var cameras []CameraResponse
	result := query.Find(&cameras)
	if result.Error != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to fetch cameras")
		return
	}

	utils.SendSuccess(c, http.StatusOK, cameras)
}

// CameraResponse2 Updated response struct
type CameraResponse2 struct {
	CameraMake  string   `json:"cameraMake"`
	CameraModel string   `json:"cameraModel"`
	Urls        []string `json:"urls"` // URLs of 3 sample photos
}

// ListCamerasWithImages godoc
// @Summary List all camera models with sample assets
// @Description Get camera models with 3 sample photo URLs for each
// @Tags cameras
// @Accept json
// @Produce json
// @Param make query string false "Filter by manufacturer"
// @Success 200 {array} CameraResponse
// @Router /cameras [get]
func (ac *AssetController) ListCamerasWithImages(c *gin.Context) {
	// First, get all camera models
	var cameras []struct {
		Make  string `json:"make"`
		Model string `json:"model"`
	}

	ac.db.Model(&models.PHAsset{}).
		Select("DISTINCT camera_make as make, camera_model as model").
		Where("camera_make IS NOT NULL AND camera_model IS NOT NULL").
		Find(&cameras)

	// For each camera, get 3 sample assets
	var response []CameraResponse2
	for _, cam := range cameras {
		var sampleAssets []models.PHAsset
		ac.db.Model(&models.PHAsset{}).
			Where("camera_make = ? AND camera_model = ?", cam.Make, cam.Model).
			Limit(3).
			Find(&sampleAssets)

		// Extract URLs
		var assetUrls []string
		for _, asset := range sampleAssets {
			assetUrls = append(assetUrls, asset.URL)
		}

		response = append(response, CameraResponse2{
			CameraMake:  cam.Make,
			CameraModel: cam.Model,
			Urls:        assetUrls, // Add sample photo URLs
		})
	}

	utils.SendSuccess(c, http.StatusOK, gin.H{"cameras": response})
}

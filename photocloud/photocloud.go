package photocloud

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	AssetsDir     = "/var/cloud/applications/PhotoKit/photocloud/assets/" // Use external storage
	MetadataDir   = "/var/cloud/applications/PhotoKit/photocloud/metadata/"
	ThumbnailsDir = "/var/cloud/applications/PhotoKit/photocloud/thumbnails/"
	IndexFile     = "/var/cloud/applications/PhotoKit/photocloud/index.dat"
)

// PHAsset represents a photo/video asset in the system
type PHAsset struct {
	ID               int       `json:"id"`
	UserId           int       `json:"userId"`
	Filename         string    `json:"filename"`
	Named            string    `json:"named"`
	Format           string    `json:"format"`
	MediaType        string    `json:"mediaType"`
	CreationDate     time.Time `json:"creationDate"`
	ModificationDate time.Time `json:"modificationDate"`
	Albums           []int     `json:"albums"`
	Persons          []int     `json:"persons"`
	IsFavorite       bool      `json:"isFavorite"`
	IsHidden         bool      `json:"isHidden"`

	PixelWidth  int `json:"pixelWidth"`
	PixelHeight int `json:"pixelHeight"`
	// Add other fields as needed...
}

// AssetUpdate defines the fields that can be updated
type AssetUpdate struct {
	Named      *string `json:"named,omitempty"`
	Albums     *[]int  `json:"albums,omitempty"`
	Persons    *[]int  `json:"persons,omitempty"`
	IsFavorite *bool   `json:"isFavorite,omitempty"`
	IsHidden   *bool   `json:"isHidden,omitempty"`
}

type StorageSystem struct {
	mu sync.RWMutex

	// Indexes
	assetIndex    map[int]string   // ID -> filename
	userIndex     map[int][]int    // UserID -> []AssetID
	dateIndex     map[string][]int // "YYYY-MM-DD" -> []AssetID
	textIndex     map[string][]int // Lowercase words -> []AssetID
	favoriteIndex map[int]bool     // AssetID -> bool (for fast favorite queries)
	hiddenIndex   map[int]bool     // AssetID -> bool (for fast hidden queries)

	// Cache
	assetCache   map[int]*PHAsset // LRU cache
	cacheMutex   sync.Mutex
	cacheQueue   []int // LRU tracking
	maxCacheSize int

	// Dirty flags for index persistence
	indexDirty bool
}

func NewStorageSystem() (*StorageSystem, error) {
	s := &StorageSystem{
		assetIndex:    make(map[int]string),
		userIndex:     make(map[int][]int),
		dateIndex:     make(map[string][]int),
		textIndex:     make(map[string][]int),
		favoriteIndex: make(map[int]bool),
		assetCache:    make(map[int]*PHAsset),
		maxCacheSize:  1000,
	}

	// Ensure directories exist
	dirs := []string{AssetsDir, MetadataDir, ThumbnailsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Load index
	if err := s.loadIndex(); err != nil {
		log.Printf("Could not load index: %v. Rebuilding...", err)
		if err := s.rebuildIndex(); err != nil {
			return nil, fmt.Errorf("failed to rebuild index: %v", err)
		}
	}

	// Start periodic maintenance
	go s.periodicMaintenance()

	return s, nil
}

func (s *StorageSystem) UploadAsset(userID int, file io.Reader, filename string, size int64) (*PHAsset, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate ID based on timestamp to avoid collisions
	id := int(time.Now().UnixNano())

	// Create file paths
	ext := filepath.Ext(filename)
	assetFilename := fmt.Sprintf("%d%s", id, ext)
	assetPath := filepath.Join(AssetsDir, assetFilename)
	metaPath := filepath.Join(MetadataDir, fmt.Sprintf("%d.json", id))

	// Save asset file
	assetFile, err := os.Create(assetPath)
	if err != nil {
		return nil, err
	}
	defer assetFile.Close()

	if _, err := io.Copy(assetFile, file); err != nil {
		return nil, err
	}

	// Create asset metadata
	now := time.Now()
	asset := &PHAsset{
		ID:           id,
		UserId:       userID,
		Filename:     assetFilename,
		Named:        strings.TrimSuffix(filename, ext),
		CreationDate: now,
		Format:       ext,
	}

	// Save metadata to file
	metaData, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		return nil, err
	}

	// Update indexes
	s.assetIndex[asset.ID] = asset.Filename
	s.userIndex[asset.UserId] = append(s.userIndex[asset.UserId], asset.ID)

	dateKey := asset.CreationDate.Format("2006-01-02")
	s.dateIndex[dateKey] = append(s.dateIndex[dateKey], asset.ID)

	// Update text index
	words := strings.Fields(strings.ToLower(asset.Named))
	for _, word := range words {
		if len(word) > 2 {
			s.textIndex[word] = append(s.textIndex[word], asset.ID)
		}
	}

	// Periodically save index instead of on every upload
	if len(s.assetIndex)%100 == 0 {
		go s.saveIndex()
	}

	return asset, nil
}

func (s *StorageSystem) SearchAssets(query string) ([]*PHAsset, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	words := strings.Fields(query)

	// Find matching IDs for each word
	idSets := make([]map[int]bool, len(words))
	for i, word := range words {
		ids := s.textIndex[word]
		idSet := make(map[int]bool)
		for _, id := range ids {
			idSet[id] = true
		}
		idSets[i] = idSet
	}

	// Intersection of all word matches
	resultIDs := make(map[int]bool)
	for id := range idSets[0] {
		inAll := true
		for i := 1; i < len(idSets); i++ {
			if !idSets[i][id] {
				inAll = false
				break
			}
		}
		if inAll {
			resultIDs[id] = true
		}
	}

	// Get assets
	var assets []*PHAsset
	for id := range resultIDs {
		asset, err := s.getAsset(id)
		if err != nil {
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

// Other optimized methods (GetByUser, GetByDate, etc.) similar to before
// but using the getAsset method that utilizes the cache

// ... [Additional methods] ...

// UpdateAsset handles partial updates to PHAsset metadata
func (s *StorageSystem) UpdateAsset(id int, update AssetUpdate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Load the current asset
	asset, err := s.getAssetFromDisk(id)
	if err != nil {
		return fmt.Errorf("asset not found: %w", err)
	}

	// Track which indexes need updating
	textIndexUpdate := false
	favoriteIndexUpdate := false

	// 2. Apply updates
	if update.Named != nil {
		// Update text index if name changed
		if *update.Named != asset.Named {
			textIndexUpdate = true
		}
		asset.Named = *update.Named
	}

	if update.Albums != nil {
		asset.Albums = *update.Albums
	}

	if update.Persons != nil {
		asset.Persons = *update.Persons
	}

	if update.IsFavorite != nil {
		if *update.IsFavorite != asset.IsFavorite {
			favoriteIndexUpdate = true
		}
		asset.IsFavorite = *update.IsFavorite
	}

	if update.IsHidden != nil {
		asset.IsHidden = *update.IsHidden
	}

	// 3. Update modification timestamp
	asset.ModificationDate = time.Now()

	// 4. Save updated asset to disk
	if err := s.saveAssetToDisk(asset); err != nil {
		return fmt.Errorf("failed to save asset: %w", err)
	}

	// 5. Update in-memory cache
	s.cacheMutex.Lock()
	if cached, exists := s.assetCache[id]; exists {
		// Update cached version
		*cached = *asset
		// Move to end of LRU queue
		for i, cachedID := range s.cacheQueue {
			if cachedID == id {
				s.cacheQueue = append(s.cacheQueue[:i], s.cacheQueue[i+1:]...)
				break
			}
		}
		s.cacheQueue = append(s.cacheQueue, id)
	}
	s.cacheMutex.Unlock()

	// 6. Update indexes if needed
	if textIndexUpdate {
		s.updateTextIndex(id, asset.Named)
		s.indexDirty = true
	}

	if favoriteIndexUpdate {
		s.favoriteIndex[id] = asset.IsFavorite
		s.indexDirty = true
	}

	return nil
}

// updateTextIndex updates the text index for an asset
func (s *StorageSystem) updateTextIndex(id int, newName string) {
	// Remove old words
	if oldAsset, err := s.getAssetFromDisk(id); err == nil {
		oldWords := strings.Fields(strings.ToLower(oldAsset.Named))
		for _, word := range oldWords {
			if len(word) > 2 {
				s.removeFromTextIndex(id, word)
			}
		}
	}

	// Add new words
	newWords := strings.Fields(strings.ToLower(newName))
	for _, word := range newWords {
		if len(word) > 2 {
			s.textIndex[word] = append(s.textIndex[word], id)
		}
	}
}

// removeFromTextIndex removes an asset from the text index for a specific word
func (s *StorageSystem) removeFromTextIndex(id int, word string) {
	ids, exists := s.textIndex[word]
	if !exists {
		return
	}

	for i, assetID := range ids {
		if assetID == id {
			// Remove element without preserving order (faster)
			s.textIndex[word] = append(ids[:i], ids[i+1:]...)
			break
		}
	}

	// Clean up empty entries
	if len(s.textIndex[word]) == 0 {
		delete(s.textIndex, word)
	}
}

// getAssetFromDisk loads asset directly from disk (bypasses cache)
func (s *StorageSystem) getAssetFromDisk(id int) (*PHAsset, error) {
	// Get metadata file path
	metaPath := filepath.Join(MetadataDir, fmt.Sprintf("%d.json", id))

	// Read file
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var asset PHAsset
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, err
	}

	return &asset, nil
}

// saveAssetToDisk saves asset to disk
func (s *StorageSystem) saveAssetToDisk(asset *PHAsset) error {
	// Get metadata file path
	metaPath := filepath.Join(MetadataDir, fmt.Sprintf("%d.json", asset.ID))

	// Marshal to JSON
	data, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(metaPath, data, 0644)
}

// getAsset retrieves asset (uses cache when possible)
func (s *StorageSystem) getAsset(id int) (*PHAsset, error) {
	// Try cache first
	s.cacheMutex.Lock()
	if asset, exists := s.assetCache[id]; exists {
		// Move to end of LRU queue
		for i, cachedID := range s.cacheQueue {
			if cachedID == id {
				s.cacheQueue = append(s.cacheQueue[:i], s.cacheQueue[i+1:]...)
				break
			}
		}
		s.cacheQueue = append(s.cacheQueue, id)
		s.cacheMutex.Unlock()
		return asset, nil
	}
	s.cacheMutex.Unlock()

	// Not in cache, load from disk
	asset, err := s.getAssetFromDisk(id)
	if err != nil {
		return nil, err
	}

	// Add to cache
	s.cacheMutex.Lock()
	if len(s.assetCache) >= s.maxCacheSize {
		// Evict least recently used
		evictID := s.cacheQueue[0]
		delete(s.assetCache, evictID)
		s.cacheQueue = s.cacheQueue[1:]
	}
	s.assetCache[id] = asset
	s.cacheQueue = append(s.cacheQueue, id)
	s.cacheMutex.Unlock()

	return asset, nil
}

// Additional helper methods for index management
func (s *StorageSystem) periodicMaintenance() {
	ticker := time.NewTicker(30 * time.Minute)
	saveTicker := time.NewTicker(5 * time.Minute)

	for {
		select {
		case <-ticker.C:
			if err := s.rebuildIndex(); err != nil {
				log.Printf("Periodic maintenance failed: %v", err)
			}

		case <-saveTicker.C:
			if s.indexDirty {
				if err := s.saveIndex(); err != nil {
					log.Printf("Index save failed: %v", err)
				} else {
					s.indexDirty = false
					log.Println("Index saved successfully")
				}
			}
		}
	}
}

// rebuildIndex reconstructs all indexes from disk metadata
func (s *StorageSystem) rebuildIndex() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Starting index rebuild...")
	start := time.Now()

	// Clear existing indexes
	s.assetIndex = make(map[int]string)
	s.userIndex = make(map[int][]int)
	s.dateIndex = make(map[string][]int)
	s.textIndex = make(map[string][]int)
	s.favoriteIndex = make(map[int]bool)

	// Get list of metadata files
	files, err := os.ReadDir(MetadataDir)
	if err != nil {
		return err
	}

	// Process files in batches
	batchSize := 5000
	for i := 0; i < len(files); i += batchSize {
		end := i + batchSize
		if end > len(files) {
			end = len(files)
		}

		batch := files[i:end]
		for _, file := range batch {
			if file.IsDir() {
				continue
			}

			// Only process JSON files
			if filepath.Ext(file.Name()) != ".json" {
				continue
			}

			// Parse asset ID from filename
			var id int
			_, err := fmt.Sscanf(file.Name(), "%d.json", &id)
			if err != nil {
				continue
			}

			// Load asset
			asset, err := s.getAssetFromDisk(id)
			if err != nil {
				continue
			}

			// Add to indexes
			s.assetIndex[asset.ID] = asset.Filename
			s.userIndex[asset.UserId] = append(s.userIndex[asset.UserId], asset.ID)

			dateKey := asset.CreationDate.Format("2006-01-02")
			s.dateIndex[dateKey] = append(s.dateIndex[dateKey], asset.ID)

			s.favoriteIndex[asset.ID] = asset.IsFavorite

			words := strings.Fields(strings.ToLower(asset.Named))
			for _, word := range words {
				if len(word) > 2 {
					s.textIndex[word] = append(s.textIndex[word], asset.ID)
				}
			}
		}

		log.Printf("Processed %d/%d files", end, len(files))
	}

	// Mark index as dirty to trigger save
	s.indexDirty = true

	log.Printf("Index rebuild completed in %v", time.Since(start))
	return nil
}

// saveIndex persists indexes to disk
func (s *StorageSystem) saveIndex() error {
	// Create index data structure
	indexData := struct {
		AssetIndex    map[int]string   `json:"assetIndex"`
		UserIndex     map[int][]int    `json:"userIndex"`
		DateIndex     map[string][]int `json:"dateIndex"`
		TextIndex     map[string][]int `json:"textIndex"`
		FavoriteIndex map[int]bool     `json:"favoriteIndex"`
	}{
		AssetIndex:    s.assetIndex,
		UserIndex:     s.userIndex,
		DateIndex:     s.dateIndex,
		TextIndex:     s.textIndex,
		FavoriteIndex: s.favoriteIndex,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(indexData, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(IndexFile, data, 0644)
}

// loadIndex loads indexes from disk
func (s *StorageSystem) loadIndex() error {
	data, err := os.ReadFile(IndexFile)
	if err != nil {
		return err
	}

	// Parse index data
	var indexData struct {
		AssetIndex    map[int]string   `json:"assetIndex"`
		UserIndex     map[int][]int    `json:"userIndex"`
		DateIndex     map[string][]int `json:"dateIndex"`
		TextIndex     map[string][]int `json:"textIndex"`
		FavoriteIndex map[int]bool     `json:"favoriteIndex"`
	}

	if err := json.Unmarshal(data, &indexData); err != nil {
		return err
	}

	// Apply to system
	s.assetIndex = indexData.AssetIndex
	s.userIndex = indexData.UserIndex
	s.dateIndex = indexData.DateIndex
	s.textIndex = indexData.TextIndex
	s.favoriteIndex = indexData.FavoriteIndex

	return nil
}

// UpdateAssetHandler API Handler for Updates
func UpdateAssetHandler(s *StorageSystem) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// Extract asset ID from URL
		vars := mux.Vars(r)
		assetID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid asset ID", http.StatusBadRequest)
			return
		}

		// Parse update payload
		var update AssetUpdate
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Apply update
		if err := s.UpdateAsset(assetID, update); err != nil {
			log.Printf("Update failed: %v", err)
			http.Error(w, "Update failed", http.StatusInternalServerError)
			return
		}

		// Return updated asset
		asset, err := s.getAsset(assetID)
		if err != nil {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(asset)
	}
}

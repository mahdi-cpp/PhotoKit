package photocloud

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

// Query represents a complex query with multiple conditions
type Query struct {
	UserID      *int       `json:"userId,omitempty"`
	IsFavorite  *bool      `json:"isFavorite,omitempty"`
	IsHidden    *bool      `json:"isHidden,omitempty"`
	MediaType   *string    `json:"mediaType,omitempty"`
	CameraMake  *string    `json:"cameraMake,omitempty"`
	CameraModel *string    `json:"cameraModel,omitempty"`
	StartDate   *time.Time `json:"startDate,omitempty"`
	EndDate     *time.Time `json:"endDate,omitempty"`
	TextSearch  *string    `json:"textSearch,omitempty"`
	AlbumID     *int       `json:"albumId,omitempty"`
	PersonID    *int       `json:"personId,omitempty"`
	MinWidth    *int       `json:"minWidth,omitempty"`
	MaxWidth    *int       `json:"maxWidth,omitempty"`
	MinHeight   *int       `json:"minHeight,omitempty"`
	MaxHeight   *int       `json:"maxHeight,omitempty"`
	Limit       int        `json:"limit,omitempty"`
	Offset      int        `json:"offset,omitempty"`
	OrderBy     string     `json:"orderBy,omitempty"` // "date", "name", "size"
	OrderDesc   bool       `json:"orderDesc,omitempty"`
}

// QueryResult contains query results with pagination info
type QueryResult struct {
	Assets     []*PHAsset `json:"assets"`
	TotalCount int        `json:"totalCount"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
}

// ExecuteQuery processes complex queries efficiently
func (s *StorageSystem) ExecuteQuery(query Query) (*QueryResult, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	startTime := time.Now()
	defer func() {
		log.Printf("Query executed in %v", time.Since(startTime))
	}()

	// Step 1: Identify which indexes we can use
	useDateIndex := query.StartDate != nil || query.EndDate != nil
	useUserIndex := query.UserID != nil
	useTextIndex := query.TextSearch != nil && len(*query.TextSearch) > 2
	useFavoriteIndex := query.IsFavorite != nil
	useHiddenIndex := query.IsHidden != nil

	// Step 2: Get candidate IDs from the most selective index
	var candidateIDs map[int]bool
	var err error

	switch {
	case useUserIndex:
		candidateIDs, err = s.getUserCandidateIDs(*query.UserID)
	case useDateIndex:
		candidateIDs, err = s.getDateCandidateIDs(query.StartDate, query.EndDate)
	case useTextIndex:
		candidateIDs, err = s.getTextCandidateIDs(*query.TextSearch)
	case useFavoriteIndex:
		candidateIDs, err = s.getFavoriteCandidateIDs(*query.IsFavorite)
	case useHiddenIndex:
		candidateIDs, err = s.getHiddenCandidateIDs(*query.IsHidden)
	default:
		// If no indexable conditions, start with all IDs
		candidateIDs = make(map[int]bool)
		for id := range s.assetIndex {
			candidateIDs[id] = true
		}
	}

	if err != nil {
		return nil, err
	}

	// Step 3: Apply additional filters in memory
	filteredIDs := s.applyFilters(candidateIDs, query)

	// Step 4: Apply sorting and pagination
	result := &QueryResult{
		TotalCount: len(filteredIDs),
		PageSize:   query.Limit,
		Page:       query.Offset/query.Limit + 1,
	}

	// Apply sorting
	sortedIDs := s.applySorting(filteredIDs, query)

	// Apply pagination
	paginatedIDs := s.applyPagination(sortedIDs, query.Offset, query.Limit)

	// Step 5: Load asset data
	assets := make([]*PHAsset, 0, len(paginatedIDs))
	for _, id := range paginatedIDs {
		asset, err := s.getAsset(id)
		if err != nil {
			log.Printf("Error loading asset %d: %v", id, err)
			continue
		}
		assets = append(assets, asset)
	}

	result.Assets = assets
	return result, nil
}

// getCandidateIDs retrieves initial candidate IDs from indexes
func (s *StorageSystem) getUserCandidateIDs(userID int) (map[int]bool, error) {
	ids, exists := s.userIndex[userID]
	if !exists {
		return make(map[int]bool), nil
	}

	result := make(map[int]bool, len(ids))
	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}

func (s *StorageSystem) getDateCandidateIDs(start, end *time.Time) (map[int]bool, error) {
	result := make(map[int]bool)

	// If no date range specified, use all dates
	if start == nil && end == nil {
		for _, ids := range s.dateIndex {
			for _, id := range ids {
				result[id] = true
			}
		}
		return result, nil
	}

	// Determine date range to query
	var startDate, endDate time.Time
	if start != nil {
		startDate = *start
	} else {
		startDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	if end != nil {
		endDate = *end
	} else {
		endDate = time.Now().Add(24 * time.Hour) // Include today
	}

	// Iterate through each day in the range
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateKey := d.Format("2006-01-02")
		if ids, exists := s.dateIndex[dateKey]; exists {
			for _, id := range ids {
				result[id] = true
			}
		}
	}

	return result, nil
}

func (s *StorageSystem) getTextCandidateIDs(query string) (map[int]bool, error) {
	words := strings.Fields(strings.ToLower(query))
	if len(words) == 0 {
		return make(map[int]bool), nil
	}

	// Find matching IDs for each word
	idSets := make([]map[int]bool, len(words))
	for i, word := range words {
		ids := s.textIndex[word]
		idSet := make(map[int]bool, len(ids))
		for _, id := range ids {
			idSet[id] = true
		}
		idSets[i] = idSet
	}

	// Start with the first word's results
	result := idSets[0]

	// Intersect with subsequent words (AND logic)
	for i := 1; i < len(idSets); i++ {
		for id := range result {
			if !idSets[i][id] {
				delete(result, id)
			}
		}
	}

	return result, nil
}

func (s *StorageSystem) getFavoriteCandidateIDs(favorite bool) (map[int]bool, error) {
	result := make(map[int]bool)
	for id, isFavorite := range s.favoriteIndex {
		if isFavorite == favorite {
			result[id] = true
		}
	}
	return result, nil
}

func (s *StorageSystem) getHiddenCandidateIDs(hidden bool) (map[int]bool, error) {
	result := make(map[int]bool)
	for id, isHidden := range s.hiddenIndex {
		if isHidden == hidden {
			result[id] = true
		}
	}
	return result, nil
}

// applyFilters applies non-indexed filters to candidate IDs
func (s *StorageSystem) applyFilters(candidates map[int]bool, query Query) []int {
	result := make([]int, 0, len(candidates))

	for id := range candidates {
		asset, err := s.getAsset(id)
		if err != nil {
			continue
		}

		// Apply all filters
		if query.IsFavorite != nil && asset.IsFavorite != *query.IsFavorite {
			continue
		}

		if query.IsHidden != nil && asset.IsHidden != *query.IsHidden {
			continue
		}

		if query.MediaType != nil && asset.MediaType != *query.MediaType {
			continue
		}

		if query.AlbumID != nil {
			found := false
			for _, album := range asset.Albums {
				if album == *query.AlbumID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if query.PersonID != nil {
			found := false
			for _, person := range asset.Persons {
				if person == *query.PersonID {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if query.MinWidth != nil && asset.PixelWidth < *query.MinWidth {
			continue
		}

		if query.MaxWidth != nil && asset.PixelWidth > *query.MaxWidth {
			continue
		}

		if query.MinHeight != nil && asset.PixelHeight < *query.MinHeight {
			continue
		}

		if query.MaxHeight != nil && asset.PixelHeight > *query.MaxHeight {
			continue
		}

		// Passed all filters
		result = append(result, id)
	}

	return result
}

// applySorting sorts asset IDs based on query parameters
func (s *StorageSystem) applySorting(ids []int, query Query) []int {
	if len(ids) == 0 || query.OrderBy == "" {
		return ids
	}

	// Load necessary attributes for sorting
	type sortableAsset struct {
		ID   int
		Date time.Time
		Name string
		Size int64 // Using pixel width as size proxy
	}

	assets := make([]sortableAsset, len(ids))
	for i, id := range ids {
		asset, err := s.getAsset(id)
		if err != nil {
			continue
		}

		assets[i] = sortableAsset{
			ID:   id,
			Date: asset.CreationDate,
			Name: asset.Named,
			Size: int64(asset.PixelWidth),
		}
	}

	// Sort based on requested field
	switch query.OrderBy {
	case "date":
		sort.Slice(assets, func(i, j int) bool {
			if query.OrderDesc {
				return assets[i].Date.After(assets[j].Date)
			}
			return assets[i].Date.Before(assets[j].Date)
		})
	case "name":
		sort.Slice(assets, func(i, j int) bool {
			if query.OrderDesc {
				return assets[i].Name > assets[j].Name
			}
			return assets[i].Name < assets[j].Name
		})
	case "size":
		sort.Slice(assets, func(i, j int) bool {
			if query.OrderDesc {
				return assets[i].Size > assets[j].Size
			}
			return assets[i].Size < assets[j].Size
		})
	}

	// Extract sorted IDs
	sortedIDs := make([]int, len(assets))
	for i, asset := range assets {
		sortedIDs[i] = asset.ID
	}

	return sortedIDs
}

// applyPagination applies pagination to the sorted IDs
func (s *StorageSystem) applyPagination(ids []int, offset, limit int) []int {
	if limit <= 0 {
		limit = 50 // Default page size
	}

	if offset < 0 {
		offset = 0
	}

	if offset >= len(ids) {
		return []int{}
	}

	end := offset + limit
	if end > len(ids) {
		end = len(ids)
	}

	return ids[offset:end]
}

// QueryAssetsByFavoriteAndVisibility is a specialized function for your use case
func (s *StorageSystem) QueryAssetsByFavoriteAndVisibility(favorite, hidden bool, limit int) ([]*PHAsset, error) {
	query := Query{
		IsFavorite: &favorite,
		IsHidden:   &hidden,
		Limit:      limit,
		OrderBy:    "date",
		OrderDesc:  true,
	}

	result, err := s.ExecuteQuery(query)
	if err != nil {
		return nil, err
	}

	return result.Assets, nil
}

// QueryAssetsByCamera  is a specialized function for your use case
func (s *StorageSystem) QueryAssetsByCamera(cameraModel string, hidden bool, limit int) ([]*PHAsset, error) {
	query := Query{
		CameraModel: &cameraModel,
		IsHidden:    &hidden,
		Limit:       limit,
		OrderBy:     "date",
		OrderDesc:   true,
	}

	result, err := s.ExecuteQuery(query)
	if err != nil {
		return nil, err
	}

	return result.Assets, nil
}

// QueryAssetsHandler API Handler for Query
func QueryAssetsHandler(s *StorageSystem) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters from request
		var query Query
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			http.Error(w, "Invalid query format", http.StatusBadRequest)
			return
		}

		// Apply default limit
		if query.Limit <= 0 || query.Limit > 1000 {
			query.Limit = 1000
		}

		result, err := s.ExecuteQuery(query)
		if err != nil {
			http.Error(w, "Query execution failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}

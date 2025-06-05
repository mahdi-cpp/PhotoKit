package repository

import (
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/model"
	"github.com/mahdi-cpp/PhotoKit/utils"
)

var pinnedCollectionDTO PinnedCollectionDTO
var pinnedCollectionDTO2 PinnedCollectionDTO

type PinnedCollectionDTO struct {
	PinnedCollections []PinnedCollection `json:"pinnedCollections"`
}

type PinnedCollection struct {
	Name  string        `json:"name"`
	Type  string        `json:"type"`
	Icon  string        `json:"icon"`
	Image model.UIImage `json:"image"`
}

func GetPinned(folder string) {

	var file = "data.txt"
	var uiImages = cache.ReadOfFile(folder, file)
	var count = len(uiImages)

	if count > 10 {
		count = 10
	}

	var index = 0
	var nameIndex = 0

	for i := 0; i < count; i++ {
		if nameIndex >= len(utils.FackNames) {
			nameIndex = 0
		}

		var pinned = PinnedCollection{}
		pinned.Name = utils.FackNames[nameIndex]
		pinned.Image = uiImages[index]

		if index == 0 {
			pinned.Name = "Favourite"
			pinned.Type = "favourite"
			pinned.Icon = "icons8-favourite-60.png"
			pinned.Image.Name = "chat_19"
		} else if index == 1 {
			pinned.Name = "Map"
			pinned.Type = "map"
			pinned.Icon = "icons8-albums-50.png"
			pinned.Image.Name = "Screenshot from 2024-08-08 01-04-57"
		} else if index == 2 {
			pinned.Name = "Trips"
			pinned.Type = "trips"
			pinned.Icon = "icons8-trip-50.png"
		} else if index == 3 {
			pinned.Name = "Videos"
			pinned.Type = "videos"
			pinned.Icon = "icons8-video-60.png"
		} else if index == 4 {
			pinned.Name = "Album"
			pinned.Type = "album"
			pinned.Icon = "icons8-albums-50.png"
		} else if index == 5 {
			pinned.Name = "Electronic"
			pinned.Type = "album"
			pinned.Icon = "icons8-albums-50.png"
		}

		pinnedCollectionDTO.PinnedCollections = append(pinnedCollectionDTO.PinnedCollections, pinned)

		nameIndex++
		index++
	}

	index = 0
}

func GetPinnedGallery(folder string) {

	var file = "data.txt"
	var photos = cache.ReadOfFile(folder, file)
	var count = len(photos)

	var index = 0
	var nameIndex = 0

	for i := 0; i < count; i++ {
		if nameIndex >= len(utils.FackNames) {
			nameIndex = 0
		}

		var pinned = PinnedCollection{}
		pinned.Name = utils.FackNames[nameIndex]
		pinned.Image = photos[index]

		pinnedCollectionDTO2.PinnedCollections = append(pinnedCollectionDTO2.PinnedCollections, pinned)

		nameIndex++
		index++
	}

	index = 0
}

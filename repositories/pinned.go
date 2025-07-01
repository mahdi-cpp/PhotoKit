package repositories

import (
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/models"
	"github.com/mahdi-cpp/PhotoKit/utils"
)

var pinnedCollectionDTO PinnedCollectionDTO
var pinnedCollectionDTO2 PinnedCollectionDTO

type PinnedCollectionDTO struct {
	PinnedCollections []PinnedCollection `json:"pinnedCollections"`
}

type PinnedCollection struct {
	Name  string         `json:"name"`
	Type  string         `json:"type"`
	Icon  string         `json:"icon"`
	Image models.UIImage `json:"image"`
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

		if index == 2 {
			pinned.Name = "Favourite"
			pinned.Type = "favourite"
			pinned.Icon = "icons8-favourite-60"
			pinned.Image.Named = "chat_19"
		} else if index == 3 {
			pinned.Name = "Map"
			pinned.Type = "map"
			pinned.Icon = "icons8-albums-50"
			pinned.Image.Named = "Screenshot from 2024-08-08 01-04-57"
		} else if index == 0 {
			pinned.Name = "Camera"
			pinned.Type = "trips"
			pinned.Icon = "camera_photo_51"
			pinned.Image.Named = "IMG_20141015_185832"
		} else if index == 1 {
			pinned.Name = "Screenshots"
			pinned.Type = "trips"
			pinned.Icon = "screenshot_60"
			pinned.Image.Named = "all_84"
		} else if index == 4 {
			pinned.Name = "Videos"
			pinned.Type = "videos"
			pinned.Icon = "icons8-video-60"
			pinned.Image.Named = "IMG_20141015_185832"
		} else if index == 5 {
			pinned.Name = "الکترونیک"
			pinned.Type = "album"
			pinned.Icon = "all_84"
			pinned.Image.Named = "b44b4f5b11b6d88022746825379a323f0badc1c2_1697894819 (1)"
		} else if index == 6 {
			pinned.Name = "Telegram"
			pinned.Type = "album"
			pinned.Icon = ""
			pinned.Image.Named = "021"
		} else if index == 7 {
			pinned.Name = "Trips"
			pinned.Type = "icons8-trip-50"
			pinned.Icon = "icons8-albums-50"
			pinned.Image.Named = "IMG_20141015_185832"
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

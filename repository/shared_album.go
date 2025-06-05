package repository

import (
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/model"
	"github.com/mahdi-cpp/PhotoKit/utils"
)

var shareAlbumDTO ShareAlbumDTO

type ShareAlbumDTO struct {
	Albums []ShareAlbum `json:"albums"`
}

type ShareAlbum struct {
	Avatar    model.UIImage   `json:"avatar"`
	Username  string          `json:"username"`
	AlbumName string          `json:"albumName"`
	Images    []model.UIImage `json:"images"`
}

func GetShareAlbums(folder string) ShareAlbumDTO {

	var file = "data.txt"
	var uiImages = cache.ReadOfFile(folder, file)
	var shareAlbumDTO ShareAlbumDTO

	var count = len(uiImages) / 6
	var index = 0
	var nameIndex = 0

	if count > 25 {
		count = 25
	}

	for i := 0; i < count; i++ {
		var shareAlbum = ShareAlbum{}
		if nameIndex+1 >= len(utils.FackNames) {
			nameIndex = 0
		}

		var avatar = model.UIImage{}
		avatar.Name = "chat_25"
		avatar.Size.Width = 500
		avatar.Size.Height = 500
		shareAlbum.Avatar = avatar
		shareAlbum.Username = "username5"
		shareAlbum.AlbumName = utils.FackTrips[nameIndex]

		for j := 0; j < 5; j++ {
			var image model.UIImage
			image = uiImages[index+2+j]
			shareAlbum.Images = append(shareAlbum.Images, image)
		}

		shareAlbumDTO.Albums = append(shareAlbumDTO.Albums, shareAlbum)
		index += 6
		nameIndex += 1
	}

	return shareAlbumDTO
}

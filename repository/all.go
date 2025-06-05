package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/model"
	"github.com/mahdi-cpp/PhotoKit/utils"
	"sync"
)

var root = "/"

func InitPhotos() {

	//GetGalleries("/var/cloud/9/", false)
	//GetGalleries("/var/cloud/id/mahan/", false)
	GetGalleries("/var/cloud/00-all/", false)

	//GetGalleries("/var/cloud/00-instagram/razzle-photo/", false)
	//GetGalleries("/var/cloud/00-instagram/razzle/", true)
	//GetGalleries("/var/cloud/00-instagram/ashtonhall/", true)
	//GetGalleries("/var/cloud/00-instagram/video/", true)

	var a = "/var/cloud/00-instagram/razzle-photo/"
	var m = "/var/cloud/id/mahan/"

	GetProfile()
	GetRecently(m)
	GetPeoples("/var/cloud/00-instagram/ashtonhall/")
	GetTrips(m)

	GetPinned("/var/cloud/id/mahan/")
	GetPinnedGallery(a)

	GetYears("/var/cloud/00-all/")

	albumDTO = GetAlbums(a)

	//shareAlbumDTO = GetShareAlbums("/var/cloud/id/mahan/")
	shareAlbumDTO = GetShareAlbums("/var/cloud/id/mahan/")

	cameraDTO = GetCameras(a)

	utils.GetCities()
	utils.GetNames()
}

var newSubTitle *SubtitleDTO

func RestSubtitle() map[string]any {
	return gin.H{
		"name":          newSubTitle.Name,
		"subtitleItems": newSubTitle.Subtitles,
	}
}

func ReloadSubtitle() {
	newSubTitle, _ = GetSubtitle()
}

func RestFeed() map[string]any {
	return gin.H{
		"profileDTO":          profileDTO,
		"galleryDTO":          galleryDTO,
		"recentDaysDTO":       recentDaysDTO,
		"peopleDTO":           peopleDTO,
		"tripDTO":             tripDTO,
		"pinnedCollectionDTO": pinnedCollectionDTO,
		"albumDTO":            albumDTO,
		"shareAlbumDTO":       shareAlbumDTO,
		"cameraDTO":           cameraDTO,
	}
}

func RestGallery() map[string]any {
	return gin.H{
		"galleryDTO": galleryDTO,
	}
}

func RestRecentDays() map[string]any {
	return gin.H{
		"recentDaysDTO": recentDaysDTO,
	}
}

func RestCamera() map[string]any {
	return gin.H{
		"cameraDTO": cameraDTO,
	}
}

func RestAlbums() map[string]any {
	return gin.H{
		"albumDTO": albumDTO,
	}
}

func RestTrips() map[string]any {
	return gin.H{
		"tripDTO": tripDTO,
	}
}
func RestPinnedCollections() map[string]any {
	return gin.H{
		"pinnedCollectionDTO": pinnedCollectionDTO2,
	}
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func RestUser() map[string]any {
	return gin.H{
		"Id":    "12",
		"name":  "Mahdi",
		"email": "mahdi.cpp@gmail.com",
	}
}

type UIImageCache struct {
	sync.RWMutex
	Cache map[int]model.UIImage
}

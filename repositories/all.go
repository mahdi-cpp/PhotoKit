package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/mahdi-cpp/PhotoKit/models"
	"github.com/mahdi-cpp/PhotoKit/utils"
	"sync"
)

var root = "/"

func InitPhotos() {

	//FetchLibraries("/var/cloud/4/", false)
	//FetchLibraries("/var/cloud/00-all/", false)
	//FetchLibraries("/var/cloud/family/", false)
	//FetchLibraries("/var/cloud/00-instagram/razzle-photo/", false)

	//FetchLibraries("/var/cloud/00-instagram/ashtonhall/", true)
	//FetchLibraries("/var/cloud/00-instagram/razzle/", true)
	//FetchLibraries("/var/cloud/00-instagram/video/", true)

	var a = "/var/cloud/family/"
	var m = "/var/cloud/family/"

	GetRecently("/var/cloud/00-instagram/ashtonhall/")
	GetPeoples("/var/cloud/people/")
	GetTrips(m)
	GetPinned("/var/cloud/00-instagram/razzle-photo/")
	GetPinnedGallery(a)

	albumDTO = GetAlbums(a)

	shareAlbumDTO = GetShareAlbums("/var/cloud/00-instagram/razzle-photo/")
	cameraDTO = GetCameras("/var/cloud/00-instagram/video/")

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

func RestCollections() map[string]any {
	return gin.H{
		"recentDaysDTO":       recentDaysDTO,
		"peopleDTO":           peopleDTO,
		"tripDTO":             tripDTO,
		"pinnedCollectionDTO": pinnedCollectionDTO,
		"albumDTO":            albumDTO,
		"shareAlbumDTO":       shareAlbumDTO,
		"cameraDTO":           cameraDTO,
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

//type User struct {
//	Id    int    `json:"id"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//}

func RestUser() map[string]any {
	return gin.H{
		"Id":    "12",
		"name":  "Mahdi",
		"email": "mahdi.cpp@gmail.com",
	}
}

type UIImageCache struct {
	sync.RWMutex
	Cache map[int]models.UIImage
}

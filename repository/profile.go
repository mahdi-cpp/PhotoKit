package repository

import (
	"github.com/mahdi-cpp/PhotoKit/model"
)

var profileDTO ProfileDTO

type ProfileDTO struct {
	User   User          `json:"user"`
	Avatar model.UIImage `json:"avatar"`
	Name   string        `json:"name"`
}

func GetProfile() {

	var image = model.UIImage{}
	image.Name = "chat_31"
	image.Size.Width = 100
	image.Size.Height = 100

	profileDTO.Avatar = image
	profileDTO.Name = "Mahdi Abdolmaleki"

}

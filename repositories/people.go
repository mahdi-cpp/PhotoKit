package repositories

import (
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/models"
	"github.com/mahdi-cpp/PhotoKit/utils"
)

var peopleDTO PeopleDTO

type PeopleDTO struct {
	PeopleGroup []PeopleGroup `json:"peopleGroupArray"`
}

type PeopleGroup struct {
	Name  string         `json:"name"`
	Photo models.UIImage `json:"photo"`
}

func GetPeoples(folder string) {

	var file = "data.txt"
	var uiImages = cache.ReadOfFile(folder, file)

	var count = len(uiImages) - 1

	if count > 15 {
		count = 15
	}

	var index = 0
	var nameIndex = 0

	for i := 0; i < count; i++ {
		var personGroup = PeopleGroup{}

		if nameIndex+1 >= len(utils.FackNames) {
			nameIndex = 0
		}

		personGroup.Name = utils.FackNames[nameIndex]
		personGroup.Photo = uiImages[index+1]
		peopleDTO.PeopleGroup = append(peopleDTO.PeopleGroup, personGroup)

		nameIndex++
		index += 1
	}

	index = 0
}

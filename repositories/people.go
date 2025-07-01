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
	Names  []string       `json:"names"`
	Photo1 models.UIImage `json:"photo1"`
	Photo2 models.UIImage `json:"photo2"`
}

func GetPeoples(folder string) {

	var file = "data.txt"
	var uiImages = cache.ReadOfFile(folder, file)

	var count = (len(uiImages) / 2)

	if count > 15 {
		count = 15
	}

	var index = 0
	var nameIndex = 0

	for i := 0; i < count; i++ {
		var personGroup = PeopleGroup{}

		if nameIndex+2 >= len(utils.FackNames) {
			nameIndex = 0
		}

		personGroup.Names = append(personGroup.Names, utils.FackNames[nameIndex])
		personGroup.Names = append(personGroup.Names, utils.FackNames[nameIndex+1])

		personGroup.Photo1 = uiImages[index+1]
		personGroup.Photo2 = uiImages[index+2]

		peopleDTO.PeopleGroup = append(peopleDTO.PeopleGroup, personGroup)

		nameIndex++
		index += 2
	}

	index = 0
}

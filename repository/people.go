package repository

import (
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/model"
	"github.com/mahdi-cpp/PhotoKit/utils"
)

var peopleDTO PeopleDTO

type PeopleDTO struct {
	PeopleGroup []PeopleGroup `json:"peopleGroupArray"`
}

type PeopleGroup struct {
	Names  []string      `json:"names"`
	Photo1 model.UIImage `json:"photo1"`
	Photo2 model.UIImage `json:"photo2"`
	Photo3 model.UIImage `json:"photo3"`
	Photo4 model.UIImage `json:"photo4"`
}

func GetPeoples(folder string) {

	var file = "data.txt"
	var uiImages = cache.ReadOfFile(folder, file)

	var count = (len(uiImages) / 4) - 4

	if count > 15 {
		count = 15
	}

	var index = 0
	var nameIndex = 0

	for i := 0; i < count; i++ {
		var personGroup = PeopleGroup{}

		if nameIndex+4 >= len(utils.FackNames) {
			nameIndex = 0
		}

		personGroup.Names = append(personGroup.Names, utils.FackNames[nameIndex])
		personGroup.Names = append(personGroup.Names, utils.FackNames[nameIndex+1])
		personGroup.Names = append(personGroup.Names, utils.FackNames[nameIndex+2])
		personGroup.Names = append(personGroup.Names, utils.FackNames[nameIndex+3])

		personGroup.Photo1 = uiImages[index+1]
		personGroup.Photo2 = uiImages[index+2]
		personGroup.Photo3 = uiImages[index+3]
		personGroup.Photo4 = uiImages[index+4]

		peopleDTO.PeopleGroup = append(peopleDTO.PeopleGroup, personGroup)

		nameIndex++
		index += 4
	}

	index = 0
}

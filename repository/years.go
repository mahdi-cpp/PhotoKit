package repository

import (
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/model"
	"strconv"
)

var yearsDTO YearsDTO

type YearsDTO struct {
	Years []Year `json:"years"`
}

type Year struct {
	Title string        `json:"title"`
	Photo model.UIImage `json:"photo"`
}

func GetYears(folder string) {

	var file = "data.txt"
	uiImages := cache.ReadOfFile(folder, file)
	var count = len(uiImages)

	var index = 0

	if count > 20 {
		count = 20
	}

	var yearsCount = 1403

	for i := 0; i < count; i++ {
		var year = Year{}
		year.Photo = uiImages[index]
		year.Title = "" + strconv.Itoa(yearsCount)

		yearsDTO.Years = append(yearsDTO.Years, year)
		index++
		yearsCount--
	}
}

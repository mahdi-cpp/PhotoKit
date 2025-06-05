package repository

import (
	"fmt"
	"github.com/asticode/go-astisub"
)

var subtitleDTO SubtitleDTO

type SubtitleDTO struct {
	Name      string `json:"name"`
	Subtitles []Item `json:"subtitles"`
}

type Item struct {
	Index   int      `json:"index"`
	EndAt   float64  `json:"endAt"`
	Lines   []string `json:"lines"`
	StartAt float64  `json:"startAt"`
}

func GetSubtitle() (*SubtitleDTO, error) {

	s, err := astisub.OpenFile("/home/mahdi/files/srt/Hacksaw.Ridge.2016.720p.1080p.BluRay.UTF.Fa.srt")

	fmt.Println("aliali45")

	if err != nil {
		fmt.Printf("Can not read File", err)
		return nil, err
	}

	subtitleDTO = SubtitleDTO{}

	for _, item := range s.Items {
		var newItem Item

		for _, line := range item.Lines {
			newItem.Lines = append(newItem.Lines, line.String())
		}
		newItem.Index = item.Index
		newItem.StartAt = item.StartAt.Seconds()
		newItem.EndAt = item.EndAt.Seconds()

		subtitleDTO.Subtitles = append(subtitleDTO.Subtitles, newItem)
	}

	return &subtitleDTO, nil

}

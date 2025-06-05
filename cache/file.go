package cache

import (
	"encoding/json"
	"fmt"
	"github.com/mahdi-cpp/PhotoKit/model"
	"os"
)

// InputJSON represents the structure of the input JSON
type InputJSON struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Orientation int    `json:"orientation"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	FileType    string `json:"fileType"`
}

func ReadOfFile(folder string, file string) []model.UIImage {

	var inputImages []InputJSON

	// Open the file for reading
	f, err := os.Open(folder + file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer f.Close() // Ensure the file is closed when we're done

	// Create a JSON decoder and decode the data into the slice
	decoder := json.NewDecoder(f)
	if err := decoder.Decode(&inputImages); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	// Convert to UIImage structs
	var outputImages []model.UIImage

	for _, img := range inputImages {

		aspectRatio := float32(img.Height) / float32(img.Width)

		output := model.UIImage{
			Name:        img.Name,
			FileType:    img.FileType,
			Orientation: img.Orientation,
			AspectRatio: aspectRatio,
			Size: model.CGSize{
				Width:  float32(img.Width),
				Height: float32(img.Height),
			},
			VideoInfo: model.VideoInfo{
				IsVideo:         false, // Default value
				VideoDuration:   0,     // Default value
				HasSubtitle:     false, // Default value
				VideoFormat:     "",    // Default value
				HasVideoControl: false, // Default value
			},
		}

		outputImages = append(outputImages, output)
	}

	return outputImages
}

//func ReadOfFile(folder string, file string) []model.UIImage {
//
//	var photos []model.UIImage
//
//	// Open the file for reading
//	f, err := os.Open(folder + file)
//	if err != nil {
//		fmt.Println("Error opening file:", err)
//		return nil
//	}
//	defer f.Close() // Ensure the file is closed when we're done
//
//	// Create a JSON decoder and decode the data into the slice
//	decoder := json.NewDecoder(f)
//	if err := decoder.Decode(&photos); err != nil {
//		fmt.Println("Error decoding JSON:", err)
//		return nil
//	}
//
//	return photos
//}

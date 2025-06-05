package repository

import (
	"fmt"
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/model"
	"os"
	"path/filepath"
	"strings"
)

type GalleryDTO struct {
	Images []model.UIImage `json:"images"`
}

var galleryDTO GalleryDTO
var imagesArray []model.UIImage

func GetGalleries(folder string, IsVideo bool) {

	var file = "data.txt"
	uiImages := cache.ReadOfFile(folder, file)
	var count = len(uiImages)
	var index = 0

	dir := "/home/mahdi/files/"
	videoFormats, err := ListVideoFormatsInDirectory(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for i := 0; i < count; i++ {
		var image = model.UIImage{}
		image = uiImages[index]
		image.VideoInfo = model.VideoInfo{
			HasSubtitle:     true,
			IsVideo:         IsVideo,
			HasVideoControl: false,
			VideoFormat:     videoFormats[image.Name],
		}

		galleryDTO.Images = append(galleryDTO.Images, image)
		index++
	}
}

// ListVideoFormatsInDirectory reads a directory and returns a map of video formats.
func ListVideoFormatsInDirectory(dir string) (map[string]string, error) {
	videoInfoArray := make(map[string]string)

	// Read the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Check each file
	for _, file := range files {
		if !file.IsDir() { // Skip directories
			filename := file.Name()
			format, err := GetVideoFormat(filename)
			if err == nil {
				//// Retrieve video duration
				//duration, err := getVideoDuration(filepath.Join(dir, filename))
				//if err != nil {
				//	fmt.Println("Error getting duration:", err)
				//	continue
				//}

				filename = strings.ReplaceAll(filename, ".mp4", "")
				filename = strings.ReplaceAll(filename, ".mkv", "")
				videoInfoArray[filename] = format
				//VideoInfo{Format: format}
				//videoInfoArray. = format // Store the filename and format
				//VideoInfo[filename] = map[string]interface{}{
				//	"format":   format,
				//	"duration": duration,
				//}
				fmt.Println(filename)
			} else {
				// Optionally log the unsupported format
				fmt.Println(err)
			}
		}
	}

	return videoInfoArray, nil
}

// GetVideoFormat returns the video format based on the file extension.
func GetVideoFormat(filename string) (string, error) {
	// Get the file extension
	ext := strings.ToLower(filepath.Ext(filename))

	// Check for video formats
	switch ext {
	case ".mp4":
		return "mp4", nil
	case ".mkv":
		return "mkv", nil
	default:
		return "", fmt.Errorf("unsupported video format: %s", ext)
	}
}

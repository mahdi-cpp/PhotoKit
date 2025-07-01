package repositories

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/PhotoKit/utils"
	"image"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var PHAssetsPath = "/var/cloud/applications/PhotoKit/Assets/"
var uploadPath = "/var/cloud/applications/PhotoKit/upload/"

func checkAssetExists(named string) (bool, error) {
	var exists bool
	err := db.Model(&PHAsset{}).
		Select("count(*) > 0").
		Where("named = ?", named).
		Find(&exists).
		Error
	return exists, err
}

func CreateAssetOfUploadDirectory(id int) {

	files, err := os.ReadDir(uploadPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {

		var named = file.Name()

		// Check if asset exists by name
		exists, err := checkAssetExists(named)
		if err != nil {
			fmt.Printf("Error checking product: %v\n", err)
			continue
		}

		if exists {
			fmt.Printf("Asset '%s' exists\n", named)
			continue
		}

		if strings.HasSuffix(file.Name(), ".jpg") ||
			strings.HasSuffix(file.Name(), ".JPG") ||
			strings.HasSuffix(file.Name(), ".jpeg") ||
			strings.HasSuffix(file.Name(), ".JPEG") {

			var assetUrl = uuid.New().String()
			var assetFormat = ".jpg"

			err = utils.CopyFile(uploadPath+file.Name(), PHAssetsPath+assetUrl+assetFormat)
			if err != nil {
				panic(err)
			}

			var portrait = false
			var Orientation = 0

			var a = PHAssetsPath + assetUrl + ".jpg"
			if utils.PhotoHasExifData(a) {
				has, orientation := utils.ReadExifData(a)
				if has {
					fmt.Println("Orientation: ", orientation)
					if strings.Compare(orientation, "6") == 0 {
						portrait = true
					}

					i, err := strconv.Atoi(orientation)
					if err != nil {
						fmt.Println("Orientation: ", err)
					} else {
						Orientation = i
					}
				}
			} else {
				fmt.Println("not exif data")
			}

			w, h := getImageDimension(PHAssetsPath + assetUrl + ".jpg")
			var width = 0
			var height = 0
			if Orientation == 6 {
				width = h
				height = w
			} else {
				width = w
				height = h
			}

			newPHAsset := PHAsset{
				Url:          assetUrl,
				Named:        file.Name(),
				MediaType:    "image",
				Format:       "jpg",
				Orientation:  Orientation,
				PixelWidth:   width,
				PixelHeight:  height,
				CreationDate: time.Now(),
			}

			// Save the chat to the database
			if err := db.Create(&newPHAsset).Error; err != nil {
				log.Printf("Failed to create PHAsset: %v", err)
			} else {
				fmt.Printf("Created PHAsset: %+v\n", newPHAsset)
				CreateTinyAsset(file.Name(), assetUrl, 540, portrait)
				CreateTinyAsset(file.Name(), assetUrl, 270, portrait)
				CreateTinyAsset(file.Name(), assetUrl, 135, portrait)
				CreateTinyAsset(file.Name(), assetUrl, 70, portrait)
			}
		}
	}
}

func CreateTinyAsset(sourceName string, assetNewName string, createSize int, portrait bool) {

	file := uploadPath + sourceName
	fmt.Println("CreateTinyAsset: ", sourceName, createSize)

	srcImage, err := imaging.Open(file)
	if err != nil {
		panic(err)
	}

	var dstImage *image.NRGBA

	if portrait {
		// Resize the cropped image to width = 200px preserving the aspect ratio.
		dstImage = imaging.Resize(srcImage, 0, createSize, imaging.Lanczos)
		dstImage = imaging.Rotate270(dstImage)

	} else {
		// Resize the cropped image to width = 200px preserving the aspect ratio.
		dstImage = imaging.Resize(srcImage, createSize, 0, imaging.Lanczos)
	}

	var name2 = PHAssetsPath + "thumbnail/" + assetNewName + "_" + strconv.Itoa(createSize) + ".jpg"

	err = imaging.Save(dstImage, name2)
	if err != nil {
		panic(err)
	}
}

func getImageDimension(imagePath string) (int, int) {
	img, err := imaging.Open(imagePath) // Replace "image.jpg" with the path to your image file
	if err != nil {
		fmt.Println("Error opening image:", err)
		return 0, 0
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	fmt.Printf("Image width: %d\n", width)
	fmt.Printf("Image height: %d\n", height)
	return width, height
}

func CreateUser() {

}

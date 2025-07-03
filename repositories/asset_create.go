package repositories

import (
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/mahdi-cpp/PhotoKit/models"
	"github.com/mahdi-cpp/PhotoKit/storage"
	"github.com/mahdi-cpp/PhotoKit/utils"
	"gorm.io/gorm"
	"image"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var PHAssetsPath = "/var/cloud/applications/PhotoKit/Assets/"
var uploadPath = "/var/cloud/applications/PhotoKit/upload/"

func checkUrlExists(named string) (bool, error) {
	var exists bool
	err := db.Model(&models.PHAsset{}).
		Select("count(*) > 0").
		Where("url = ?", named).
		Find(&exists).
		Error
	return exists, err
}

func checkAssetExists(named string) (bool, error) {
	var exists bool
	err := db.Model(&models.PHAsset{}).
		Select("count(*) > 0").
		Where("named = ?", named).
		Find(&exists).
		Error
	return exists, err
}

func CreateAssetOfUploadDirectory(db1 *gorm.DB, id int) {

	db = db1

	var userIdPath = strconv.FormatInt(int64(id), 10) + "/"

	files, err := os.ReadDir(uploadPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {

		var named = file.Name()

		found, textFile, line, err := storage.SearchTextInFiles(PHAssetsPath+userIdPath, ".txt", named, true)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			//continue
		}
		if found {
			fmt.Printf("Found text in %s (line %d)\n", textFile, line)
			continue
		}

		if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".JPG") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".JPEG") {

			var assetUrl = uuid.New().String()
			var assetFormat = ".jpg"

			err = storage.CopyFile(uploadPath+file.Name(), PHAssetsPath+userIdPath+assetUrl+assetFormat)
			if err != nil {
				panic(err)
			}

			//textFile, err := os.OpenFile(PHAssetsPath+userIdPath+assetUrl+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			//if err != nil {
			//	panic(err)
			//}
			//defer textFile.Close()
			//
			//if _, err := textFile.WriteString(file.Name() + "\n"); err != nil {
			//	panic(err)
			//}

			var portrait = false
			var Orientation = 0

			var a = PHAssetsPath + userIdPath + assetUrl + ".jpg"

			var cameraMake = ""
			var cameraModel = ""

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

				cMake, cModel, err := utils.GetCameraModel(a)
				if err != nil {
					log.Printf("Warning: error getting camera info: %v", err)
					cameraMake = ""
					cameraModel = ""
				} else {
					cameraMake = cMake
					cameraModel = cModel

					// Convert to NULL if empty after sanitization
					if cameraMake == "" {
						cameraMake = "NULL" // For raw SQL, or use sql.NullString
					}
					if cameraModel == "" {
						cameraModel = "NULL"
					}
				}

			} else {
				fmt.Println("not exif data")
			}

			w, h := getImageDimension(PHAssetsPath + userIdPath + assetUrl + ".jpg")
			var width = 0
			var height = 0
			if Orientation == 6 {
				width = h
				height = w
			} else {
				width = w
				height = h
			}

			asset := models.PHAsset{
				URL:         assetUrl,
				Named:       named,
				MediaType:   "image",
				Format:      "jpg",
				Orientation: Orientation,

				CameraMake:  cameraMake,
				CameraModel: cameraModel,

				PixelWidth:  width,
				PixelHeight: height,

				CreationDate: time.Now(),
			}

			// Save the asset
			err = storage.SaveAsset(&asset, userIdPath)
			if err != nil {
				fmt.Println("Error saving asset:", err)
				return
			}

			// Save the chat to the database
			if err := db.Debug().Create(&asset).Error; err != nil {
				log.Printf("Failed to create PHAsset: %v", err)
			} else {
				fmt.Printf("Created PHAsset: %+v\n", asset)
				CreateTinyAsset(file.Name(), assetUrl, userIdPath, 540, portrait)
				CreateTinyAsset(file.Name(), assetUrl, userIdPath, 270, portrait)
				CreateTinyAsset(file.Name(), assetUrl, userIdPath, 135, portrait)
				CreateTinyAsset(file.Name(), assetUrl, userIdPath, 70, portrait)
			}
		}
	}
}

func CreateOnlyDatabase(db1 *gorm.DB, userId int) {

	db = db1

	var userIdPath = strconv.FormatInt(int64(userId), 10) + "/"

	files, err := os.ReadDir(PHAssetsPath + userIdPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {

		var named = strings.Replace(file.Name(), ".jpg", "", 1)

		// Check if asset exists by name
		exists, err := checkUrlExists(named)
		if err != nil {
			fmt.Printf("Error checking product: %v\n", err)
			continue
		}

		if exists {
			fmt.Printf("Asset '%s' exists\n", named)
			continue
		}

		if strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".JPG") || strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".JPEG") {

			var Orientation = 0
			var a = PHAssetsPath + userIdPath + file.Name()

			var cameraMake = ""
			var cameraModel = ""

			if utils.PhotoHasExifData(a) {
				has, orientation := utils.ReadExifData(a)

				if has {
					fmt.Println("Orientation: ", orientation)
					if strings.Compare(orientation, "6") == 0 {
						//portrait = true
					}

					i, err := strconv.Atoi(orientation)
					if err != nil {
						fmt.Println("Orientation: ", err)
					} else {
						Orientation = i
					}
				}

				cMake, cModel, err := utils.GetCameraModel(a)
				if err != nil {
					log.Printf("Warning: error getting camera info: %v", err)
					cameraMake = ""
					cameraModel = ""
				} else {
					cameraMake = cMake
					cameraModel = cModel

					// Convert to NULL if empty after sanitization
					if cameraMake == "" {
						cameraMake = "NULL" // For raw SQL, or use sql.NullString
					}
					if cameraModel == "" {
						cameraModel = "NULL"
					}
				}

			} else {
				fmt.Println("not exif data")
			}

			w, h := getImageDimension(PHAssetsPath + userIdPath + file.Name())
			var width = 0
			var height = 0
			if Orientation == 6 {
				width = h
				height = w
			} else {
				width = w
				height = h
			}

			newPHAsset := models.PHAsset{
				UserId:      userId,
				URL:         named,
				MediaType:   "image",
				Format:      "jpg",
				Orientation: Orientation,

				CameraMake:  cameraMake,
				CameraModel: cameraModel,

				PixelWidth:  width,
				PixelHeight: height,

				CreationDate: time.Now(),
			}

			// Save the chat to the database
			if err := db.Debug().Create(&newPHAsset).Error; err != nil {
				log.Printf("Failed to create PHAsset: %v", err)
			}
		}
	}
}

func CreateTinyAsset(sourceName string, assetNewName string, userIdPath string, createSize int, portrait bool) {

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

	var name2 = PHAssetsPath + userIdPath + "thumbnail/" + assetNewName + "_" + strconv.Itoa(createSize) + ".jpg"

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

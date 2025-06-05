package cache

import (
	"bytes"
	"fmt"
	"github.com/mahdi-cpp/PhotoKit/model"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var folders = []string{

	"/var/cloud/00-all/",
	"/var/cloud/00-all/thumbnail/",

	"/var/cloud/4/",
	"/var/cloud/4/thumbnail/",

	"/var/cloud/all/",
	"/var/cloud/all/thumbnail/",

	"/var/cloud/music/albums/",
	"/var/cloud/music/albums/thumbnail/",

	"/var/cloud/paris/",
	"/var/cloud/paris/thumbnail/",

	"/var/cloud/reynardlowell/",
	"/var/cloud/reynardlowell/thumbnail/",

	"/var/cloud/ikea/",
	"/var/cloud/ikea/thumbnail/",

	"/var/cloud/wallpaper/",

	"/var/cloud/id/me/",
	"/var/cloud/id/me/thumbnail",
	"/var/cloud/id/mahan/",
	"/var/cloud/id/mahan/thumbnail",

	"/var/cloud/ik/",
	"/var/cloud/ik/thumbnail",

	"/var/cloud/id/ali/",
	"/var/cloud/id/ali2/",
	"/var/cloud/id/trip/",
	"/var/cloud/id/go/",

	"/var/cloud/chat/movie/movie/",

	"/var/cloud/chat/movie/animation/",

	"/var/cloud/chat/movie/movie/thumbnail/",
	"/var/cloud/chat/movie/animation/thumbnail/",

	"/var/cloud/camera-sequrity/",
	"/var/cloud/camera-sequrity/thumbnail/",

	"/var/cloud/chat/pdf/",
	"/var/cloud/chat/electronic/",
	"/var/cloud/chat/map/",
	"/var/cloud/chat/voice/",

	"/var/cloud/chat/pdf/thumbnail/",
	"/var/cloud/chat/electronic/thumbnail/",
	"/var/cloud/chat/map/thumbnail/",
	"/var/cloud/chat/voice/thumbnail/",

	"/var/cloud/call2/",
	"/var/cloud/call2/thumbnail/",

	"/var/cloud/id/ali/thumbnail/",
	"/var/cloud/id/ali2/thumbnail/",
	"/var/cloud/id/face/thumbnail/",
	"/var/cloud/id/trip/thumbnail/",
	"/var/cloud/id/go/thumbnail/",

	"/var/cloud/chat_users/",
	"/var/cloud/chat_users/thumbnail/",

	"/var/cloud/tinyhome/new/",
	"/var/cloud/tinyhome/new/thumbnail/",

	"/var/cloud/tinyhome/06/",
	"/var/cloud/tinyhome/06/thumbnail/",

	"/var/cloud/video/",
	"/var/cloud/video/thumbnail/",

	"/var/cloud/elec/",
	"/var/cloud/elec//thumbnail/",
	//----------------------------------------------------
	"/var/cloud/00-instagram/razzle-photo/",
	"/var/cloud/00-instagram/razzle-photo/thumbnail/",

	"/var/cloud/00-instagram/ashtonhall/",
	"/var/cloud/00-instagram/ashtonhall/thumbnail/",

	"/var/cloud/00-instagram/lucaspopan/",
	"/var/cloud/00-instagram/lucaspopan/thumbnail/",

	"/var/cloud/00-instagram/razzle/",
	"/var/cloud/00-instagram/razzle/thumbnail/",

	"/var/cloud/00-instagram/nickloveswildlife/",
	"/var/cloud/00-instagram/nickloveswildlife/thumbnail/",

	"/var/cloud/00-instagram/video/",
	"/var/cloud/00-instagram/video/thumbnail/",
}

var iconFolder = "/var/cloud/icons/"

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

func ReadIcons() {
	// Specify the directory you want to read
	dir := "/var/cloud/icons" // Change this to your target directory

	// Read the directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("failed to read directory: %v", err)
	}

	// Iterate over the entries
	for _, entry := range entries {
		if !entry.IsDir() { // Check if it is not a directory

			if strings.Contains(entry.Name(), ".png") {
				addIconCash(entry.Name())
				//fmt.Printf("Reading file: %s\n", entry.Name())
			}
		}
	}

	fmt.Println(len(iconCache.cache))
}

type ImageCache struct {
	sync.RWMutex
	cache map[string][]byte
}

var thumbCache = ImageCache{cache: make(map[string][]byte)}

type IconCache struct {
	sync.RWMutex
	cache map[string][]byte
}

var iconCache = IconCache{cache: make(map[string][]byte)}

func GetThumbCash(filename string) ([]byte, bool) {
	thumbCache.RLock()
	imgData, exists := thumbCache.cache[filename]
	thumbCache.RUnlock()
	return imgData, exists
}

func GetIconCash(filename string) ([]byte, bool) {
	iconCache.RLock()
	imgData, exists := iconCache.cache[filename]
	iconCache.RUnlock()
	return imgData, exists
}

// ConvertImageToBytes converts an image.Image to a byte slice.
func ConvertImageToBytes(img image.Image, format string) ([]byte, error) {
	var buf bytes.Buffer
	var err error

	// Encode the image based on the specified format
	switch format {
	case "jpg":
		err = jpeg.Encode(&buf, img, nil)
	case "png":
		err = png.Encode(&buf, img)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func AddThumbCash(filepath string, filename string) {

	originalImage, err := LoadImage(filepath)
	if err != nil {
		fmt.Println("addToCash Error loading image:", err)
		return
	}

	imgBytes, err := ConvertImageToBytes(originalImage, "jpg") // Change to "png" for PNG format
	if err != nil {
		fmt.Println("Error ConvertImageToBytes: ", err)
		return
	}

	thumbCache.Lock()
	thumbCache.cache[filename] = imgBytes
	thumbCache.Unlock()
}

func addIconCash(iconName string) {
	icon, err := LoadImage("/var/cloud/icons/" + iconName) // Change path accordingly
	if err != nil {
		fmt.Println("Error loading image:", err)
		return
	}

	imgBytes, err := ConvertImageToBytes(icon, "png") // Change to "png" for PNG format
	if err != nil {
		fmt.Println("Error ConvertImageToBytes: ", err)
		return
	}

	iconCache.Lock()
	iconCache.cache[iconName] = imgBytes
	iconCache.Unlock()
}

func SearchFile(filename string) (string, error) {
	for _, folder := range folders {
		// Construct the full path to the file
		fullPath := filepath.Join(folder, filename)

		// Check if the file exists
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil // File found
		} else if os.IsNotExist(err) {
			// File does not exist in this directory
			continue
		} else {
			// Other error (e.g., permission issues)
			return "", err
		}
	}
	return "", fmt.Errorf("file %s not found in any of the specified folders", filename)
}

// LoadImage loads an image from a file.
func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

//--------------------------------

// IDGenerator is a struct that holds the current ID and a mutex for thread safety
type IDGenerator struct {
	currentID int
	mu        sync.Mutex
}

// NewIDGenerator creates a new IDGenerator instance
func NewIDGenerator() *IDGenerator {
	return &IDGenerator{
		currentID: 0,
	}
}

// NextID generates the next unique ID
func (g *IDGenerator) NextID() int {
	g.mu.Lock()         // Lock to ensure thread safety
	defer g.mu.Unlock() // Unlock after generating the ID
	g.currentID++       // Increment the current ID
	return g.currentID  // Return the new ID
}

var IdGen = NewIDGenerator()

type UIImageCache struct {
	sync.RWMutex
	Cache map[int]model.UIImage
}

var UIImageMemory = UIImageCache{Cache: make(map[int]model.UIImage)}

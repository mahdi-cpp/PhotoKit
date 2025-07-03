package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mahdi-cpp/PhotoKit/models"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	AssetsBaseDir = "/var/cloud/applications/PhotoKit/Assets/" // Directory to store JSON files
)

// Ensure directory exists
func init() {
	if _, err := os.Stat(AssetsBaseDir); os.IsNotExist(err) {
		os.MkdirAll(AssetsBaseDir, 0755)
	}
}

// SaveAsset saves a PHAsset to a JSON file
func SaveAsset(asset *models.PHAsset, userId string) error {

	// Create filename based on ID and creation date
	filename := filepath.Join(AssetsBaseDir+userId, asset.URL+".json")

	// Convert to JSON
	data, err := json.MarshalIndent(asset, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filename, data, 0644)
}

// LoadAsset loads a PHAsset from JSON file by ID
func LoadAsset(id int) (*models.PHAsset, error) {

	files, err := os.ReadDir(AssetsBaseDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if file.IsDir() {
			continue
		}

		// Read file
		data, err := os.ReadFile(filepath.Join(AssetsBaseDir, file.Name()))
		if err != nil {
			continue // Skip files we can't read
		}

		var asset models.PHAsset
		if err := json.Unmarshal(data, &asset); err != nil {
			continue // Skip invalid JSON
		}

		if asset.ID == id {
			return &asset, nil
		}
	}

	return nil, os.ErrNotExist
}

// LoadAllAssets loads all PHAssets from JSON files
func LoadAllAssets() ([]models.PHAsset, error) {
	files, err := os.ReadDir(AssetsBaseDir)
	if err != nil {
		return nil, err
	}

	var assets []models.PHAsset
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(AssetsBaseDir, file.Name()))
		if err != nil {
			continue
		}

		var asset models.PHAsset
		if err := json.Unmarshal(data, &asset); err != nil {
			continue
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

// DeleteAsset removes a PHAsset JSON file
func DeleteAsset(id int) error {
	files, err := os.ReadDir(AssetsBaseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data, err := os.ReadFile(filepath.Join(AssetsBaseDir, file.Name()))
		if err != nil {
			continue
		}

		var asset models.PHAsset
		if err := json.Unmarshal(data, &asset); err != nil {
			continue
		}

		if asset.ID == id {
			return os.Remove(filepath.Join(AssetsBaseDir, file.Name()))
		}
	}

	return os.ErrNotExist
}

func FileInDirectoryExists(dir, filename string) bool {
	path := filepath.Join(dir, filename)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CopyFile(src, dst string) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// ErrTextFound Create a custom error to signal when we've found our text
var ErrTextFound = errors.New("text found in file")

// SearchTextInFiles searches for text in all files with given extension in directory
// Returns (found, filePath, lineNumber, error)
func SearchTextInFiles(dir, ext, text string, stopOnFirst bool) (bool, string, int, error) {

	var foundFile string
	var foundLine int

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.EqualFold(filepath.Ext(path), ext) {

			found, line, err := searchTextInFile(path, text)
			if err != nil {
				return nil // Continue with next file
			}
			if found {
				foundFile = path
				foundLine = line
				if stopOnFirst {
					return ErrTextFound
				}
			}
		}
		return nil
	})

	if err == ErrTextFound {
		return true, foundFile, foundLine, nil
	}
	if err != nil {
		return false, "", 0, fmt.Errorf("directory walk error: %w", err)
	}
	return foundFile != "", foundFile, foundLine, nil
}

// SearchTextInFile searches for text in a single file
// Returns (found, lineNumber, error)
func searchTextInFile(filename, text string) (bool, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), text) {
			return true, lineNumber, nil
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return false, 0, fmt.Errorf("scan error: %w", err)
	}

	return false, 0, nil
}

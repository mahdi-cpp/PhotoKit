package main

import (
	"github.com/mahdi-cpp/PhotoKit/cache"
	"github.com/mahdi-cpp/PhotoKit/repository"
)

func main() {

	repository.InitPhotos()
	cache.ReadIcons()

	Run()
}

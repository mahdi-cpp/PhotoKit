package main

import (
	"fmt"
	"github.com/mahdi-cpp/PhotoKit/config"
	"github.com/mahdi-cpp/PhotoKit/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
)

var db *gorm.DB
var err error

func main() {

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "api_v1.", // schema name
			SingularTable: false,
		}})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err != nil {
		println("Failed to connect database gallery\"")
		os.Exit(1)
	}

	//// Initialize repositories
	//repo := repositories.NewRepository(db)
	//
	//// Initialize handler
	//assetHandler := routes.NewAssetHandler(repo)
	//
	//// Set up routes
	//router.GET("/v1/assets", assetHandler.GetAllAssets)
	//// Setup routes
	//routes.SetupUserRoutes(router, db)
	//routes.SetupAssetRoutes(router, db)

	// Create repositories
	//userRepo := repositories.NewUserRepository(db)

	// Create a new user
	//newUser := &models.User{
	//	Username:    "mahdi",
	//	PhoneNumber: "+989355512619",
	//	Email:       "mahdi.cpp@gmail.com",
	//	FirstName:   "mahdi",
	//	LastName:    "abdolmaleki",
	//	Bio:         "Software developer",
	//	IsOnline:    true,
	//	LastSeen:    time.Now(),
	//}

	//err := userRepo.CreateUser(newUser)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Created user with ID: %d\n", newUser.ID)

	repositories.CreateAssetOfUploadDirectory(db, 1)
	//repositories.CreateOnlyDatabase(db, 1)

	//repositories.InitPhotos()
	//cache.ReadIcons()
	//Run()
}

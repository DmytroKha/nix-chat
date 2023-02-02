package main

import (
	"fmt"
	"github.com/DmytroKha/nix-chat/config"
	_ "github.com/DmytroKha/nix-chat/docs"
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/filesystem"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/router"
	_ "github.com/go-sql-driver/mysql"
	mysqlG "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// @title       nix_chat API
// @version     1.0
// @description API Server for nix_chat application.

// @host     localhost:8090
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
func main() {
	var conf = config.GetConfiguration()

	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v",
		conf.DatabaseUser,
		conf.DatabasePassword,
		conf.DatabaseHost,
		conf.DatabaseName)
	db, err := gorm.Open(mysqlG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	userRepository := database.NewUserRepository(db)
	userService := app.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	authService := app.NewAuthService(userService, conf)
	authController := controllers.NewAuthController(authService, userService)

	imageRepository := database.NewImageRepository(db)
	imageStorageService := filesystem.NewImageStorageService(conf.FileStorageLocation)
	imageService := app.NewImageService(imageRepository, imageStorageService)
	imageController := controllers.NewImageController(imageService)

	e := router.New(
		userController,
		authController,
		imageController,
	)

	// service start at port :8090
	err = e.Start(":8090")
	if err != nil {
		log.Fatalln(err)
	}
}

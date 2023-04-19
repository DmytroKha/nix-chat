package main

import (
	"context"
	"fmt"
	"github.com/DmytroKha/nix-chat/config"
	_ "github.com/DmytroKha/nix-chat/docs"
	"github.com/DmytroKha/nix-chat/internal/app"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/DmytroKha/nix-chat/internal/infra/filesystem"
	"github.com/DmytroKha/nix-chat/internal/infra/http/controllers"
	"github.com/DmytroKha/nix-chat/internal/infra/http/router"
	"github.com/DmytroKha/nix-chat/internal/infra/http/websocket"
	_ "github.com/go-sql-driver/mysql"
	mysqlG "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// @title       nix_chat
// @version     1.0
// @description Server for nix_chat application.

// @host     localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()
	var conf = config.GetConfiguration()

	config.CreateRedisClient()

	err := database.Migrate(conf)
	if err != nil {
		log.Fatalf("Unable to apply migrations: %q\n", err)
	}

	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?%v",
		conf.DatabaseUser,
		conf.DatabasePassword,
		conf.DatabaseHost,
		conf.DatabaseName,
		"parseTime=true")
	db, err := gorm.Open(mysqlG.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	imageRepository := database.NewImageRepository(db)
	imageStorageService := filesystem.NewImageStorageService(conf.FileStorageLocation)
	imageService := app.NewImageService(imageRepository, imageStorageService)
	imageController := controllers.NewImageController(imageService)

	userRepository := database.NewUserRepository(db)
	userService := app.NewUserService(userRepository, imageService)
	userController := controllers.NewUserController(userService)

	blacklistRepository := database.NewBlacklistRepository(db)
	blacklistService := app.NewBlacklistService(blacklistRepository)

	friendlistRepository := database.NewFriendlistRepository(db)
	friendlistService := app.NewFriendlistService(friendlistRepository)

	roomRepository := database.NewRoomRepository(db)

	wsServer := websocket.NewWebsocketServer(roomRepository, userRepository, blacklistService, friendlistService)
	go wsServer.Run(ctx)

	authService := app.NewAuthService(userService, conf)
	authController := controllers.NewAuthController(authService, userService, wsServer)

	e := router.New(
		userController,
		authController,
		imageController,
		wsServer,
	)

	// service start at port :8080
	err = e.Start(":8080")
	if err != nil {
		log.Fatalln(err)
	}
}

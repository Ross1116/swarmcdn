package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/config"
	"github.com/ross1116/swarmcdn/handlers"
	"github.com/ross1116/swarmcdn/utils"
)

func main() {
	config.InitConfig()

	app := utils.NewApp(*config.AppConfig)

	router := gin.Default()
	router.POST("/upload", handlers.UploadHandler(app))
	router.GET("/chunks/:hash", handlers.GetChunkHandler)
	router.GET("/manifest/:fileID", handlers.GetLatestManifestHandler)

	router.GET("/peers", handlers.GetKnownPeers)
	router.POST("/peers/register", handlers.AddKnownPeer)

	router.Run(":8080")
}

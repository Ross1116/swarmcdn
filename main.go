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
	router.POST("/upload", handlers.MakeUploadHandler(app))

	router.Run(":8080")
}

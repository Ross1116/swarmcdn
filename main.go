package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/handlers"
)

func main() {
	router := gin.Default()

	router.POST("/upload", handlers.UploadHandler)

	router.Run(":8080")
}

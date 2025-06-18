package main

import (
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/chunks/:hash", GetChunkHandler)

	port := choosePort("9000", "9001")
	router.Run(":" + port)
}

func choosePort(primary, fallback string) string {
	ln, err := net.Listen("tcp", ":"+primary)
	if err == nil {
		_ = ln.Close()
		return primary
	}
	return fallback
}

func GetChunkHandler(c *gin.Context) {
	hash := c.Param("hash")
	path := filepath.Join("chunks", hash+".blob")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chunk not found"})
		return
	}

	c.File(path)
}

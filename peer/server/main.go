package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const ChunksDir = "client/chunks"

func ServeChunks(port string) {
	router := gin.Default()
	router.GET("/chunks/:hash", GetChunkHandler)
	router.GET("/health", CheckHealthHandler)

	// port := choosePort("9000", "9001")
	router.Run(":" + port)
}

// func choosePort(primary, fallback string) string {
// 	ln, err := net.Listen("tcp", ":"+primary)
// 	if err == nil {
// 		_ = ln.Close()
// 		return primary
// 	}
// 	return fallback
// }

func GetChunkHandler(c *gin.Context) {
	hash := c.Param("hash")
	path := filepath.Join(ChunksDir, hash+".blob")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chunk not found"})
		return
	}

	c.File(path)
}

func CheckHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

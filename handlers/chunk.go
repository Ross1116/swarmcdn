package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetChunkHandler(c *gin.Context) {
	hash := c.Param("hash")
	path := filepath.Join("storage/chunks", hash+".blob")

	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chunk not found"})
		return
	}

	c.File(path)
}

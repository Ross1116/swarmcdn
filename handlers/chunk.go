package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

func GetChunkHandler(c *gin.Context) {
	hash := c.Param("hash")
	path := utils.GetChunkPath(hash)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chunk not found"})
		return
	}

	c.File(path)
}

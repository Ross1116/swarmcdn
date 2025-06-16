package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

func GetLatestManifestHandler(c *gin.Context) {
	fileID := c.Param("fileID")

	index, err := utils.LoadIndex("storage/index.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read index"})
		return
	}

	var foundIndex *utils.FileIndex
	for _, entry := range index {
		if entry.FileID == fileID {
			foundIndex = &entry
			break
		}
	}

	if foundIndex == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File ID not found in index"})
		return
	}

	manifestPath := filepath.Join("storage", "manifests", fileID, fmt.Sprintf("v%d.json", foundIndex.LatestVersion))
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Manifest file not found"})
		return
	}

	c.File(manifestPath)
}

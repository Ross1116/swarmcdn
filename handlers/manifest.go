package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

func GetLatestManifestHandler(c *gin.Context) {
	fileID := c.Param("fileID")

	indexPath := utils.GetIndexFilePath()
	index, err := utils.LoadIndex(indexPath)
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

	manifestPath := utils.GetManifestPath(fileID, foundIndex.LatestVersion)
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Manifest file not found"})
		return
	}

	c.File(manifestPath)
}

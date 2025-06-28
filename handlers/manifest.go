package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

func GetLatestManifestHandler(c *gin.Context) {
	username := c.Param("username")
	filename := c.Param("filename")

	version, err := utils.GetNextManifestVersion(username, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to resolve version"})
		return
	}

	latestVersion := version - 1
	if latestVersion <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No versions found for file"})
		return
	}

	manifestPath := utils.GetManifestPath(username, filename, latestVersion)
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Manifest file not found"})
		return
	}

	c.File(manifestPath)
}

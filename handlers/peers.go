package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

func GetKnownPeers(c *gin.Context) {
	peerFilePath := utils.GetPeersFilePath()

	if _, err := os.Stat(peerFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Peers file not found"})
		return
	}

	c.File(peerFilePath)
}

func AddKnownPeer(c *gin.Context) {
	var request struct {
		URL string `json:"url"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	peerFilePath := utils.GetPeersFilePath()

	if _, err := os.Stat(peerFilePath); os.IsNotExist(err) {
		empty, _ := json.Marshal([]string{})
		if err := os.WriteFile(peerFilePath, empty, 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create peer file"})
			return
		}
	}

	content, err := os.ReadFile(peerFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read peer file"})
	}

	var peers []string
	err = json.Unmarshal(content, &peers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed umarshalling peers file"})
		return
	}

	for _, p := range peers {
		if p == request.URL {
			c.JSON(http.StatusOK, gin.H{"message": "Peer already registered"})
			return
		}
	}

	peers = append(peers, request.URL)
	updatedContent, err := json.MarshalIndent(peers, "", " ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to marshall the url data"})
		return
	}
	if err = os.WriteFile(peerFilePath, updatedContent, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update peer file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Peer registered successfully"})

}

package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

const ChunksDir = "peer/client/chunks"

func ServeChunks(port string) {
	router := gin.Default()
	router.GET("/chunks/:hash", GetChunkHandler)
	router.GET("/health", CheckHealthHandler)
	router.POST("/upload_chunk", UploadChunkHandler)

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

func UploadChunkHandler(c *gin.Context) {
	expectedHash := c.PostForm("hash")
	if expectedHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing hash"})
		return
	}

	fileHeader, err := c.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing chunk file"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open uploaded file"})
		return
	}
	defer file.Close()

	chunkData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded chunk"})
		return
	}

	_, err = utils.SaveChunkIfValid(chunkData, expectedHash, ChunksDir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chunk uploaded and verified"})
}

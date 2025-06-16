package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ross1116/swarmcdn/utils"
)

func UploadHandler(app *utils.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			log.Println("Failed to get form file:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received or file is invalid"})
			return
		}

		tempPath := "storage/original/" + file.Filename

		if err := c.SaveUploadedFile(file, tempPath); err != nil {
			log.Println("Failed to save file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
			return
		}
		defer os.Remove(tempPath)

		chunks, err := app.Chunker.ChunkFile(tempPath, "storage/chunks")
		if err != nil {
			log.Println("Failed to chunk file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Chunking failed"})
			return
		}

		chunkHashes := make([]string, len(chunks))
		for i, chunk := range chunks {
			chunkHashes[i] = chunk.SHA256Hash
		}

		fileID := uuid.New().String()
		version := 1
		manifest := utils.Manifest{
			FileID:     fileID,
			Version:    version,
			Filename:   file.Filename,
			Chunks:     chunkHashes,
			UploadedAt: time.Now(),
		}

		manifestDir := filepath.Join("storage", "manifests", fileID)
		_ = os.MkdirAll(manifestDir, os.ModePerm)
		manifestPath := filepath.Join(manifestDir, fmt.Sprintf("v%d.json", version))
		err = app.Manifest.SaveManifest(manifest, manifestPath)
		if err != nil {
			log.Println("Failed to save manifest", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write manifest"})
			return
		}

		indexPath := "storage/index.json"
		index, err := utils.LoadIndex(indexPath)
		if err != nil {
			log.Println("Failed to load index", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read index"})
			return
		}

		index = utils.UpdateIndexEntry(index, manifest)

		if err := utils.SaveIndex(indexPath, index); err != nil {
			log.Println("Failed to write index:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update index"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      fmt.Sprintf("'%s' uploaded and chunked!", file.Filename),
			"chunks":       len(chunks),
			"fileID":       fileID,
			"version":      version,
			"manifestPath": manifestPath,
			"indexPath":    indexPath,
		})
	}
}

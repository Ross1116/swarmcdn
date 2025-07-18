package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

func UploadHandler(app *utils.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			log.Println("Failed to get form file:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received or file is invalid"})
			return
		}

		tempPath := utils.GetOriginalPath(file.Filename)
		if err := c.SaveUploadedFile(file, tempPath); err != nil {
			log.Println("Failed to save file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
			return
		}
		defer os.Remove(tempPath)

		chunks, err := app.Chunker.ChunkFile(tempPath, utils.ChunksDir)
		if err != nil {
			log.Println("Failed to chunk file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Chunking failed"})
			return
		}

		chunkHashes := make([]string, 0, len(chunks))
		seen := make(map[string]bool)

		for _, chunk := range chunks {
			if !seen[chunk.SHA256Hash] {
				chunkHashes = append(chunkHashes, chunk.SHA256Hash)
				seen[chunk.SHA256Hash] = true
			}
		}
		basePath := filepath.Join(utils.ManifestsDir, username, file.Filename)

		version, err := utils.GetNextManifestVersion(username, file.Filename)
		if err != nil {
			log.Println("Failed to determine manifest version:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Version resolution failed"})
			return
		}

		fileHash := sha256.Sum256([]byte(username + "/" + file.Filename))
		fileID := hex.EncodeToString(fileHash[:])
		manifest := utils.Manifest{
			FileID:     fileID,
			Version:    version,
			Filename:   file.Filename,
			Chunks:     chunkHashes,
			UploadedAt: time.Now(),
		}

		manifestPath := filepath.Join(basePath, fmt.Sprintf("v%d.json", version))
		if err := os.MkdirAll(filepath.Dir(manifestPath), os.ModePerm); err != nil {
			log.Println("Failed to create manifest directory:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare manifest directory"})
			return
		}

		if err := app.Manifest.SaveManifest(manifest, manifestPath); err != nil {
			log.Println("Failed to save manifest", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write manifest"})
			return
		}

		indexPath := utils.GetIndexFilePath()
		index, err := utils.LoadIndex(indexPath)
		if err != nil {
			log.Println("Failed to load index", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read index"})
			return
		}

		index = utils.UpdateIndexEntry(index, username, manifest)

		if err := utils.SaveIndex(indexPath, index); err != nil {
			log.Println("Failed to write index:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update index"})
			return
		}

		peerList, err := utils.LoadPeerList()
		if err != nil {
			log.Println("Unable to load the peer list", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read peers file"})
			return
		}
		go func() {
			if err = RedistributeChunks(chunkHashes, peerList); err != nil {
				log.Println("Unable to redistribute the chunks to peers: ", err)
			}
		}()

		c.JSON(http.StatusOK, gin.H{
			"message":      fmt.Sprintf("'%s' uploaded and chunked!", file.Filename),
			"chunks":       len(chunkHashes),
			"fileID":       fileID,
			"version":      version,
			"manifestPath": manifestPath,
			"indexPath":    indexPath,
		})
	}
}

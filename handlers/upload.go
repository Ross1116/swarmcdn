package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ross1116/swarmcdn/utils"
)

func MakeUploadHandler(app *utils.App) gin.HandlerFunc {
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

		chunks, err := app.Chunker.ChunkFile(tempPath, "storage/chunks")
		if err != nil {
			log.Println("Failed to chunk file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Chunking failed"})
			return
		}

		_ = os.Remove(tempPath)

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("'%s' uploaded and chunked!", file.Filename),
			"chunks":  len(chunks),
		})
	}
}

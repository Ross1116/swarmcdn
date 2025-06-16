package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Failed to get form file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received or file is invalid"})
		return
	}
	savePath := "storage/original/" + file.Filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		log.Println("Failed to save file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the file"})
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

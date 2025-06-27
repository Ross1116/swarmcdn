package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
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

func RedistributeChunks(chunkHashes []string, peerURLs []string) error {
	if len(peerURLs) == 0 {
		return fmt.Errorf("no peers available")
	}

	numPeers := len(peerURLs)

	for i, hash := range chunkHashes {
		chunkPath := utils.GetChunkPath(hash)

		startIndex := i % numPeers
		success := false

		for attempt := range numPeers {
			peerIndex := (startIndex + attempt) % numPeers
			peerURL := peerURLs[peerIndex]

			err := UploadChunkToPeer(chunkPath, peerURL)
			if err != nil {
				log.Printf("Failed to upload chunk %s to peer %s: %v", hash, peerURL, err)
				continue
			}

			log.Printf("Successfully uploaded chunk %s to peer %s", hash, peerURL)
			success = true
			break
		}

		if !success {
			log.Printf("Failed to upload chunk %s to any peer", hash)
			return fmt.Errorf("failed to upload chunk %s to any peer", hash)
		}
	}

	return nil
}

func UploadChunkToPeer(chunkPath string, peerURL string) error {
	file, err := os.Open(chunkPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file data: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	uploadEndpoint := peerURL + "/upload_chunk"
	resp, err := http.Post(
		uploadEndpoint,
		writer.FormDataContentType(),
		&requestBody,
	)
	if err != nil {
		return fmt.Errorf("Unable to post the file to the peer %s: %v", peerURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s\n%s", resp.Status, string(body))
	}

	fmt.Println("File uploaded successfully")
	return nil
}

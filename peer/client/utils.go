package main

import (
	"log"
	"os"
)

type Manifest struct {
	FileID     string   `json:"file_id"`
	Filename   string   `json:"filename"`
	Version    int      `json:"version"`
	Chunks     []string `json:"chunks"`
	UploadedAt string   `json:"uploaded_at"`
}

const serverURL = "http://localhost:8080"
const maxConcurrentDownloads = 5

func InitDirectories() {
	dirs := []string{
		ChunksDir,
		ManifestsDir,
		DownloadsDir,
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func fetchManifest(reader *bufio.Reader) error {
	fmt.Println("Enter username:")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input:", err)
		return err
	}
	username = strings.TrimSpace(username)

	fmt.Println("Enter filename:")
	filename, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input:", err)
		return err
	}
	filename = strings.TrimSpace(filename)

	manifestFilePath := GetManifestPath(filename)
	url := fmt.Sprintf("%s/manifest/%s/%s", serverURL, username, filename)

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch manifest. Status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		return fmt.Errorf("failed to fetch manifest: %s", resp.Status)
	}

	manifestData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return err
	}

	if err := os.MkdirAll(ManifestsDir, 0755); err != nil {
		return fmt.Errorf("failed to create manifest directory: %v", err)
	}

	if err := os.WriteFile(manifestFilePath, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %v", err)
	}
	log.Printf("Manifest saved to %s\n", manifestFilePath)

	var manifest Manifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		log.Println("Error decoding manifest:", err)
		return err
	}

	if err := downloadChunksParallel(manifest.Chunks); err != nil {
		log.Println("One or more chunks failed to download.")
		return fmt.Errorf("could not download all chunks: %w", err)
	}
	log.Println("All chunks downloaded successfully.")

	if err := os.MkdirAll(DownloadsDir, 0755); err != nil {
		return fmt.Errorf("failed to create downloads directory: %v", err)
	}
	if err := reconstructFile(manifest, DownloadsDir); err != nil {
		return err
	}

	return nil
}

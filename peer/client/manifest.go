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
	fmt.Println("Enter File ID:")
	fileID, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading input:", err)
		return err
	}
	fileID = strings.TrimSpace(fileID)
	manifestFilePath := fmt.Sprintf("chunks/%s.json", fileID)

	resp, err := http.Get(fmt.Sprintf("%s/manifest/%s", serverURL, fileID))
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch manifest. Status: %s\n", resp.Status)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading the response: ", err)
			return err
		}

		fmt.Println(string(body))
		return err
	}

	manifestData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading the response: ", err)
		return err
	}

	manifestFile, err := os.Create(manifestFilePath)
	if err != nil {
		return fmt.Errorf("failed to create manifest file: %v", err)
	}
	defer manifestFile.Close()

	_, err = manifestFile.Write(manifestData)
	if err != nil {
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
		return err
	}

	if err := os.MkdirAll("downloads", 0755); err != nil {
		return fmt.Errorf("failed to create downloads directory: %v", err)
	}
	err = reconstructFile(manifest, "downloads")
	if err != nil {
		return err
	}

	return nil
}

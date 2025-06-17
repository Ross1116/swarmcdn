package peer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func fetchChunk(hash string) error {
	chunkFilePath := fmt.Sprintf("chunks/%s.blob", hash)

	if data, err := os.ReadFile(chunkFilePath); err == nil {
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) == hash {
			log.Printf("Chunk already exists and is valid: %s\n", chunkFilePath)
			return nil
		}
		log.Printf("Corrupt chunk detected (%s), deleting and redownloading.", hash)
		_ = os.Remove(chunkFilePath)
	}

	const maxRetries = 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Downloading chunk %s (attempt %d/%d)...", hash, attempt, maxRetries)

		resp, err := http.Get(fmt.Sprintf("%s/chunks/%s", serverURL, hash))
		if err != nil {
			log.Printf("Error fetching chunk %s: %v", hash, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Failed to fetch chunk %s: %s\n%s", hash, resp.Status, string(body))
			continue
		}

		chunkData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading chunk %s data: %v", hash, err)
			continue
		}

		// check if hashsum is same as server
		sum := sha256.Sum256(chunkData)
		if hex.EncodeToString(sum[:]) != hash {
			log.Printf("Chunk %s failed hash check. Retrying...", hash)
			continue
		}

		chunkFile, err := os.Create(chunkFilePath)
		if err != nil {
			return fmt.Errorf("error creating chunk file: %v", err)
		}
		defer chunkFile.Close()

		if _, err := chunkFile.Write(chunkData); err != nil {
			return fmt.Errorf("error writing chunk file: %v", err)
		}

		log.Printf("Chunk %s downloaded and verified successfully.", hash)
		return nil
	}

	return fmt.Errorf("failed to download valid chunk %s after %d attempts", hash, maxRetries)
}

func downloadChunksParallel(chunkHashes []string) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrentDownloads)
	errChan := make(chan error, len(chunkHashes))
	for _, hash := range chunkHashes {
		wg.Add(1)

		go func(h string) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := fetchChunk(h); err != nil {
				errChan <- fmt.Errorf("chunk %s: %v", h, err)
			}
		}(hash)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		for err := range errChan {
			log.Println("Error during chunk download:", err)
		}
		return fmt.Errorf("some chunks failed to download")
	}

	return nil
}

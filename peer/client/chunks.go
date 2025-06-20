package main

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
	chunkFilePath := GetChunkPath(hash)

	if data, err := os.ReadFile(chunkFilePath); err == nil {
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) == hash {
			log.Printf("Chunk already exists and is valid: %s\n", chunkFilePath)
			return nil
		}
		log.Printf("Corrupt chunk detected (%s), deleting and redownloading.", hash)
		_ = os.Remove(chunkFilePath)
	}

	peers := loadPeerList()
	filteredPeers := make([]string, 0, len(peers))
	for _, peer := range peers {
		if peerURL == peer {
			continue
		}
		filteredPeers = append(filteredPeers, peer)
	}
	filteredPeers = append(filteredPeers, serverURL)

	const maxRetries = 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Printf("Attempt %d to download chunk %s...", attempt, hash)

		for _, peer := range filteredPeers {
			url := fmt.Sprintf("%s/chunks/%s", peer, hash)
			log.Printf("Trying %s", url)

			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error contacting %s: %v", peer, err)
				continue
			}

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				log.Printf("Failed from %s: %s\n%s", peer, resp.Status, string(body))
				resp.Body.Close()
				continue
			}

			chunkData, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf("Error reading response from %s: %v", peer, err)
				continue
			}

			sum := sha256.Sum256(chunkData)
			if hex.EncodeToString(sum[:]) != hash {
				log.Printf("Hash mismatch from %s. Retrying...", peer)
				continue
			}

			err = os.WriteFile(chunkFilePath, chunkData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write chunk %s: %v", hash, err)
			}

			log.Printf("Chunk %s downloaded and verified from %s", hash, peer)
			return nil
		}
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

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		for _, err := range errors {
			log.Println("Error during chunk download:", err)
		}
		return fmt.Errorf("some chunks failed to download")
	}

	return nil
}

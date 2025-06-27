package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type DefaultChunker struct {
	ChunkSize int
}

func (c *DefaultChunker) ChunkFile(inputPath string, outputDir string) ([]ChunkMeta, error) {
	var chunks []ChunkMeta

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, c.ChunkSize)
	index := 0

	for {
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead == 0 {
			break
		}

		chunkData := make([]byte, bytesRead)
		copy(chunkData, buffer[:bytesRead])

		hash := sha256.Sum256(chunkData)
		hashString := hex.EncodeToString(hash[:])
		chunkFileName := fmt.Sprintf("%s.blob", hashString)

		log.Printf("Chunk %d: size=%d, hash=%s", index, bytesRead, hashString)

		if _, err := SaveChunkIfValid(chunkData, hashString, outputDir); err != nil {
			return nil, err
		}

		chunks = append(chunks, ChunkMeta{
			Filename:   chunkFileName,
			SHA256Hash: hashString,
			Index:      index,
		})
		index++
	}

	return chunks, nil
}

func SaveChunkIfValid(data []byte, expectedHash string, outputDir string) (string, error) {
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])

	if hashString != expectedHash {
		return "", fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, hashString)
	}

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	chunkPath := filepath.Join(outputDir, hashString+".blob")

	if _, err := os.Stat(chunkPath); os.IsNotExist(err) {
		if err := os.WriteFile(chunkPath, data, 0644); err != nil {
			return "", fmt.Errorf("failed to write chunk file: %w", err)
		}
	}

	return chunkPath, nil
}

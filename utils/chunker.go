package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
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

		hash := sha256.Sum256(buffer[:bytesRead])
		hashString := hex.EncodeToString(hash[:])
		chunkFileName := fmt.Sprintf("%s.blob", hashString)
		chunkFilePath := filepath.Join(outputDir, chunkFileName)

		if _, err := os.Stat(chunkFilePath); os.IsNotExist(err) {
			chunkFile, err := os.Create(chunkFilePath)
			if err != nil {
				return nil, err
			}
			_, err = chunkFile.Write(buffer[:bytesRead])
			chunkFile.Close()
			if err != nil {
				return nil, err
			}
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

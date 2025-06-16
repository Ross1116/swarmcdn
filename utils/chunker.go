package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type ChunkMeta struct {
	FileName   string `json:"file_name"`
	SHA256Hash string `json:"sha256_hash"`
	Index      int    `json:"index"`
}

type Config struct {
	ChunkSize int
	ServerURL string
}

type DefaultChunker struct {
	chunkSize int
}

type DefaultUploader struct {
	serverURL string
}

type DefaultManifestManager struct{}

func (c *DefaultChunker) ChunkFile(filepath string) ([]ChunkMeta, error) {
	var chunks []ChunkMeta

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buffer := make([]byte, c.chunkSize)
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

		chunkFileName := fmt.Sprintf("%s.chunk.%d", filepath, index)
		chunkFile, err := os.Create(chunkFileName)
		if err != nil {
			return nil, err
		}

		_, err = chunkFile.Write(buffer[:bytesRead])
		if err != nil {
			return nil, err
		}

		chunks = append(chunks, ChunkMeta{FileName: chunkFileName, SHA256Hash: hashString, Index: index})

		chunkFile.Close()

		index += 1
	}

	return chunks, nil
}

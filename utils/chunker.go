package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ross1116/swarmcdn/config"
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
	ChunkSize int
}

type DefaultUploader struct {
	ServerURL string
}

type DefaultManifestManager struct{}

type App struct {
	Config   config.Config
	Chunker  DefaultChunker
	Uploader DefaultUploader
	Manifest DefaultManifestManager
}

func NewApp(cfg config.Config) *App {
	return &App{
		Config:   cfg,
		Chunker:  DefaultChunker{ChunkSize: cfg.ChunkSize},
		Uploader: DefaultUploader{ServerURL: cfg.ServerURL},
		Manifest: DefaultManifestManager{},
	}
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
			FileName:   chunkFileName,
			SHA256Hash: hashString,
			Index:      index,
		})
		index++
	}

	return chunks, nil
}

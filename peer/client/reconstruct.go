package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func reconstructFile(manifest Manifest, outputFilePath string) error {
	outputFile, err := os.Create(
		filepath.Join(outputFilePath, filepath.Base(manifest.Filename)),
	)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	for _, hash := range manifest.Chunks {
		chunkPath := GetChunkPath(hash)

		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return fmt.Errorf("failed to read chunk %s: %v", hash, err)
		}

		_, err = outputFile.Write(chunkData)
		if err != nil {
			return fmt.Errorf("failed to write chunk %s to output: %v", hash, err)
		}
	}

	log.Printf("File reconstructed successfully at %s/%s\n", outputFilePath, manifest.Filename)
	return nil
}

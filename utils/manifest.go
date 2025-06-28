package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type DefaultManifestManager struct{}

type Manifest struct {
	FileID     string    `json:"file_id"`
	Filename   string    `json:"filename"`
	Version    int       `json:"version"`
	Chunks     []string  `json:"chunks"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (c *DefaultManifestManager) SaveManifest(manifest Manifest, outputPath string) error {
	err := os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return err
	}

	content, err := json.MarshalIndent(manifest, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *DefaultManifestManager) LoadManifest(inputPath string) (Manifest, error) {
	var manifest Manifest

	content, err := os.ReadFile(inputPath)
	if err != nil {
		return manifest, err
	}

	err = json.Unmarshal(content, &manifest)
	if err != nil {
		return manifest, err
	}

	return manifest, nil
}

func GetNextManifestVersion(username, filename string) (int, error) {
	basePath := filepath.Join(ManifestsDir, username, filename)
	files, err := os.ReadDir(basePath)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}

	version := 1
	for _, f := range files {
		var v int
		if _, err := fmt.Sscanf(f.Name(), "v%d.json", &v); err == nil && v >= version {
			version = v + 1
		}
	}
	return version, nil
}

package utils

import (
	"encoding/json"
	"os"
)

type DefaultManifestManager struct{}

type Manifest struct {
	Filename string   `json:"filename"`
	Chunks   []string `json:"chunks"`
}

func (c *DefaultManifestManager) SaveManifest(manifest Manifest, outputPath string) error {
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

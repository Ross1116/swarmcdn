package utils

import (
	"encoding/json"
	"os"
	"slices"
	"sort"
	"time"
)

type FileIndex struct {
	FileID        string    `json:"file_id"`
	Username      string    `json:"username"`
	Filename      string    `json:"filename"`
	LatestVersion int       `json:"latest_ver"`
	AllVersions   []int     `json:"all_versions"`
	UploadedAt    time.Time `json:"uploaded_at"`
	Tags          []string  `json:"tags"`
}

func SaveIndex(outputPath string, index []FileIndex) error {
	content, err := json.MarshalIndent(index, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, content, 0644)
}

func LoadIndex(inputPath string) ([]FileIndex, error) {
	var index []FileIndex

	content, err := os.ReadFile(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return index, nil
		}
		return nil, err
	}

	err = json.Unmarshal(content, &index)
	return index, err
}

func UpdateIndexEntry(index []FileIndex, username string, manifest Manifest) []FileIndex {
	for i, entry := range index {
		if entry.FileID == manifest.FileID {
			entry.LatestVersion = manifest.Version
			entry.AllVersions = appendIfMissing(entry.AllVersions, manifest.Version)
			sort.Ints(entry.AllVersions)
			entry.UploadedAt = manifest.UploadedAt
			index[i] = entry
			return index
		}
	}

	// for new file
	index = append(index, FileIndex{
		FileID:        manifest.FileID,
		Username:      username,
		Filename:      manifest.Filename,
		LatestVersion: manifest.Version,
		AllVersions:   []int{manifest.Version},
		UploadedAt:    manifest.UploadedAt,
	})
	return index
}

func appendIfMissing(list []int, v int) []int {
	if !slices.Contains(list, v) {
		list = append(list, v)
	}
	return list
}

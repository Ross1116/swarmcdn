package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type Tracker struct {
	Chunk string   `json:"chunk"`
	Peers []string `json:"peers"`
}

func SaveTrackers(outputPath string, trackers []Tracker) error {
	for _, v := range trackers {
		filename := filepath.Join(outputPath, v.Chunk+".json")

		content, err := json.MarshalIndent(v, "", " ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadTrackerFile(inputPath string) (Tracker, error) {
	var tracker Tracker

	content, err := os.ReadFile(inputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return tracker, nil
		} else {
			return tracker, err
		}
	}

	err = json.Unmarshal(content, &tracker)
	return tracker, err
}

func UpdateTrackerEntry(inputPath string, chunkHash string, newPeer string) error {
	filePath := filepath.Join(inputPath, chunkHash+".json")
	tracker, err := LoadTrackerFile(filePath)
	if err != nil {
		return err
	}

	tracker.Chunk = chunkHash

	if !slices.Contains(tracker.Peers, newPeer) {
		tracker.Peers = append(tracker.Peers, newPeer)
	}

	content, err := json.MarshalIndent(tracker, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, content, 0644)
}

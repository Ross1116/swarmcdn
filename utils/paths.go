package utils

import (
	"fmt"
	"path/filepath"
)

const (
	StorageDir   = "storage"
	OriginalDir  = StorageDir + "/original"
	ChunksDir    = StorageDir + "/chunks"
	ManifestsDir = StorageDir + "/manifests"
	TrackersDir  = StorageDir + "/trackers"
	indexFile    = StorageDir + "/index.json"
	peersFile    = StorageDir + "/peers.json"
)

// func GetManifestPath(fileID string, version int) string {
// 	return filepath.Join(ManifestsDir, fileID, fmt.Sprintf("v%d.json", version))
// }

func GetManifestPath(username, filename string, version int) string {
	return filepath.Join(ManifestsDir, username, filename, fmt.Sprintf("v%d.json", version))
}

func GetChunkPath(hash string) string {
	return filepath.Join(ChunksDir, hash+".blob")
}

func GetOriginalPath(filename string) string {
	return filepath.Join(OriginalDir, filename)
}

func GetIndexFilePath() string {
	return indexFile
}

func GetPeersFilePath() string {
	return peersFile
}

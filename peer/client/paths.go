package main

import (
	"path/filepath"
)

const (
	BaseDir      = "peer/client"
	ChunksDir    = BaseDir + "/chunks"
	DownloadsDir = BaseDir + "/downloads"
	ManifestsDir = BaseDir + "/manifests"
	peersFile    = BaseDir + "/peers.json"
)

func GetChunkPath(hash string) string {
	return filepath.Join(ChunksDir, hash+".blob")
}

func GetManifestPath(fileID string) string {
	return filepath.Join(ManifestsDir, fileID+".json")
}

func GetDownloadPath(filename string) string {
	return filepath.Join(DownloadsDir, filename)
}

func GetPeersFilePath() string {
	return peersFile
}

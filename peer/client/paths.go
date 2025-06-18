package main

import (
	"path/filepath"
)

const (
	BaseDir      = "peer/client"
	ChunksDir    = BaseDir + "/chunks"
	DownloadsDir = "peer/client/downloads"
	ManifestsDir = BaseDir + "/manifests"
	PeersFile    = BaseDir + "/peers.json"
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

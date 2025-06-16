package utils

import "github.com/ross1116/swarmcdn/config"

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

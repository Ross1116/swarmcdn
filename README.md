# SwarmCDN

SwarmCDN is a lightweight peer-to-peer content delivery network written in Go. It enables distributed file storage and retrieval across registered peers with chunk-level deduplication, replication, and manifest-based versioning.

## Features

- **Chunk-Based Storage:** Files are split into 512KB chunks, hashed via SHA256 for content-addressed storage.
- **Manifest-Based Versioning:** Every upload generates a manifest under `manifests/<username>/<filename>/vN.json`.
- **Peer Replication:** Chunks are uploaded to multiple healthy peers. The tracker maintains who holds what.
- **P2P CLI Peer Client:** Users can upload, download, reconstruct files, and serve chunks via a peer HTTP server.
- **Health Monitoring:** Dead peers are periodically detected and removed from the network automatically.
- **Central Coordinator (Server):** Handles uploads, manifest tracking, peer management, and redistribution.

## Folder Structure

    .
    ├── storage/
    │   ├── manifests/         # Versioned manifests: /username/filename/vN.json
    │   ├── chunks/            # Server-side backup of chunk blobs
    │   ├── trackers/          # tracker/<chunk_hash>.json (replication map)
    │   ├── index.json         # File index for easy lookup
    │   └── peers.json         # Peer registry
    ├── peer/
    │   ├── client/            # CLI + embedded HTTP server
    │   └── server/            # HTTP chunk receiver
    ├── utils/                 # Tracker logic, health checks, peer mgmt
    └── example_files/         # Demo files for upload

## Server (Control Plane)

    go run main.go

Handles:

- Chunking + deduplication
- Manifest generation and indexing
- Peer registration
- Upload redistribution
- Peer health monitoring

## Peer Client

    make run-client

Capabilities:

- Upload from local file system
- Automatically splits, hashes, and sends chunks
- Fetches manifest and downloads missing chunks
- Runs an HTTP server for sharing chunks
- Deduplicates chunks already present

## Manifest Format

    {
      "file_id": "uuid",
      "filename": "example.txt",
      "version": 2,
      "chunks": ["<sha256-1>", "<sha256-2>", "..."],
      "uploaded_at": "2025-06-27T14:05:10+05:30"
    }

## Roadmap

- [ ] Index manifest entries per user (for file listing/history)
- [ ] Optimize chunk downloads with retries + goroutines
- [ ] Add peer auth/token for uploads
- [ ] Enable streamable playback (MP4 chunks)
- [ ] Static file serving support (e.g., React apps)
- [ ] Add a web dashboard
- [ ] Add in-browser (WASM) peer support

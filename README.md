# SwarmCDN

SwarmCDN is a lightweight peer-to-peer content delivery network designed to enable distributed file uploads and downloads across multiple nodes. Built in Go, it supports chunked uploads, manifest tracking, peer registration, and fault-tolerant retrieval.

## Features

- 🔹 Chunk-based file upload and download
- 🔹 Manifest-based file versioning and integrity
- 🔹 Peer-to-peer chunk sharing across registered clients
- 🔹 Periodic peer health checks and auto-removal of dead nodes
- 🔹 Redistributed chunk uploading to multiple peers for balance
- 🔹 SHA-256 chunk validation to ensure data integrity

## Usage

### Peer Client

```bash
make run-client
```

Features:
- Upload files from `example_files/`
- Download files by File ID
- Auto-registers with main server
- Starts a chunk server (default port: 9000/9001)
- Chunks stored in `peer/client/chunks/`
- Downloads reconstructed to `peer/client/downloads/`

### Server (Control Plane)

Handles:
- File uploads + chunking
- Manifest generation and storage
- Peer registration (`peers.json`)
- Redistributing chunks to multiple peers

### Health Monitoring

- Periodically checks each peer via `/health`
- Automatically deletes unresponsive peers from `peers.json`

## Manifest Format

Each uploaded file generates a manifest:

```json
{
  "file_id": "uuid",
  "filename": "example.txt",
  "version": 1,
  "chunks": ["<hash1>", "<hash2>", ...],
  "uploaded_at": "2025-06-27T14:05:10+05:30"
}
```

Future: Versioning will be updated on overwrite (e.g. `v2.json`). Chunk ownership tracking will be offloaded to `tracker.json`.

## Folder Structure

```
.
├── storage/
│   ├── manifests/         # Stores manifest JSONs
│   └── chunks/            # Main server backup chunk storage
├── peer/
│   ├── client/            # CLI + chunk server
│   └── server/            # UploadChunkHandler, etc.
├── utils/                 # Peer management, health checks
└── example_files/         # Test upload files
```

## Future Work

- [ ] Chunk replication and tracker.json per chunk
- [ ] Peer authentication
- [ ] Streamable video playback
- [ ] React/HTML page serving via CDN
- [ ] Dashboard for peer monitoring
- [ ] WASM/browser-based peers

---

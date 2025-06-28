# SwarmCDN

SwarmCDN is a lightweight peer-to-peer content delivery network (CDN) written in Go. It enables distributed file storage, chunk-level deduplication, and versioned retrieval across registered peers.

---

## 🚀 Features

- **Chunk-Based Storage:** Files are split into 512KB SHA256-addressed chunks.
- **Versioned Manifests:** Each upload creates a manifest at `storage/manifests/<username>/<filename>/vN.json`.
- **Peer Replication:** Chunks are replicated to multiple healthy peers, tracked via `storage/trackers/<chunk_hash>.json`.
- **Global File Index:** `index.json` maintains metadata for all uploads: user, filename, versions, timestamps.
- **Integrity + Deduplication:** Chunks are not re-uploaded if they already exist; all chunks are hash-verified.
- **CLI Peer Client:** Lightweight CLI for upload, download, chunk serving, and file reconstruction.
- **Health Monitoring:** Automatically removes unreachable peers from the registry.
- **Parallel Downloads:** Chunk downloads happen concurrently with retries for robustness.
- **Planned: Per-User Indexing:** Enables user-specific file listings, exports, and auditability.

---

## 🗂️ Folder Structure

```

.
├── storage/
│   ├── chunks/               # Server-side backup of chunk blobs
│   ├── manifests/            # Versioned manifests per user and file
│   ├── trackers/             # Chunk replication map
│   ├── index.json            # Global file metadata index
│   └── peers.json            # Peer registry
│
├── peer/
│   ├── client/               # CLI + embedded HTTP server
│   └── server/               # Peer-side chunk receiver
│
├── utils/                    # Manifest/index helpers, health checks
└── example\_files/            # Sample files for testing

````

---

## 📄 Manifest Format

Each manifest captures a full version of a file upload.

```json
{
  "file_id": "uuid-hash",
  "filename": "example.txt",
  "version": 2,
  "chunks": [
    "sha256-hash-1",
    "sha256-hash-2"
  ],
  "uploaded_at": "2025-06-28T12:17:32+05:30"
}
````

---

## 🗃️ File Index (`index.json`)

`index.json` tracks file metadata for all uploads.

```json
[
  {
    "file_id": "sha256(username/filename)",
    "username": "ross",
    "filename": "example.txt",
    "latest_ver": 2,
    "all_versions": [1, 2],
    "uploaded_at": "2025-06-28T12:17:32+05:30",
    "tags": []
  }
]
```

Planned: `manifests/<username>/index.json` for faster per-user queries.

---

## 🧪 Running the Server

```bash
go run main.go
```

Responsible for:

* File chunking and deduplication
* Manifest and index management
* Peer registration and health monitoring
* Chunk redistribution

---

## 🧑‍🤝‍🧑 Running the Peer Client

```bash
make run-client
```

Provides:

* File upload with automatic chunking
* Manifest fetching and reconstruction
* Embedded HTTP server to serve chunks
* Deduplication for existing chunks

---

## 🔄 Upload + Download Flow

1. Upload → server chunks and hashes the file.
2. Manifest saved under user/filename/version.
3. Chunks are distributed to registered peers.
4. Index updated for metadata tracking.
5. Download reconstructs file using manifest + peer chunks (with hash validation).

---

## 🛣️ Roadmap

* [x] Manifest versioning
* [x] Global file metadata index
* [x] Peer-based chunk deduplication
* [x] Parallel downloads with retries
* [ ] Per-user index: `manifests/<username>/index.json`
* [ ] CLI: List/download file history
* [ ] Peer authentication/token support
* [ ] Static file + React app serving
* [ ] Streaming playback (e.g. MP4)
* [ ] Web dashboard
* [ ] WASM-based browser peers

package peer

type Manifest struct {
	FileID     string   `json:"file_id"`
	Filename   string   `json:"filename"`
	Version    int      `json:"version"`
	Chunks     []string `json:"chunks"`
	UploadedAt string   `json:"uploaded_at"`
}

const serverURL = "http://localhost:8080"
const maxConcurrentDownloads = 5

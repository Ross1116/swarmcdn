package utils

type ChunkMeta struct {
	FileName   string `json:"file_name"`
	SHA256Hash string `json:"sha256_hash"`
	Index      int    `json:"index"`
}

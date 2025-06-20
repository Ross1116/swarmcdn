package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/ross1116/swarmcdn/peer/server"
)

type Manifest struct {
	FileID     string   `json:"file_id"`
	Filename   string   `json:"filename"`
	Version    int      `json:"version"`
	Chunks     []string `json:"chunks"`
	UploadedAt string   `json:"uploaded_at"`
}

const serverURL = "http://localhost:8080"
const maxConcurrentDownloads = 5

var peerURL string

func InitDirectories() {
	dirs := []string{
		ChunksDir,
		ManifestsDir,
		DownloadsDir,
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
}

func choosePort(primary, fallback string) string {
	ln, err := net.Listen("tcp", ":"+primary)
	if err == nil {
		_ = ln.Close()
		return primary
	}
	return fallback
}

func runBackgroundTask() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalln("Failed to dial udp with: ", err)
		return
	}
	defer conn.Close()

	chunkPort := choosePort("9000", "9001")
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	peerURL := fmt.Sprintf("http://%s:%s", localAddr.IP.String(), chunkPort)
	err = registerPeer(peerURL)
	if err != nil {
		log.Fatalln("Failed to register peer with : ", err)
		return
	}
	fmt.Println("Client registered with peer URL:", peerURL)

	go server.ServeChunks(chunkPort)
	fmt.Println("Chunk server running in the background")
	time.Sleep(250 * time.Millisecond)
}

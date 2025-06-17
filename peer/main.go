package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter File ID: ")
	fileID, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read input:", err)
	}
	fileID = strings.TrimSpace(fileID)

	manifest, err := FetchManifest(fileID)
	if err != nil {
		log.Fatal("Manifest fetch failed:", err)
	}

	if err := DownloadChunksParallel(manifest.Chunks); err != nil {
		log.Fatal("Chunk download failed:", err)
	}

	if err := os.MkdirAll("downloads", 0755); err != nil {
		log.Fatal("Failed to create downloads directory:", err)
	}

	if err := ReconstructFile(manifest, "downloads"); err != nil {
		log.Fatal("File reconstruction failed:", err)
	}

	log.Println("Peer download completed successfully.")
}

package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	err := FetchManifest(reader)
	if err != nil {
		log.Println("Manifest fetching failed with: ", err)
		return
	}
}

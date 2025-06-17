package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	err := fetchManifest(reader)
	if err != nil {
		log.Println("Manifest fetching failed with: ", err)
		return
	}
}

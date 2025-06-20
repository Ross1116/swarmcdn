package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	InitDirectories()
	reader := bufio.NewReader(os.Stdin)

	runBackgroundTask()

	for {
		fmt.Println("\nChoose an action:")
		fmt.Println("1. Fetch Manifest")
		fmt.Println("2. Upload File")
		fmt.Println("3. Exit")
		fmt.Print("Enter choice [1-3]: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input:", err)
			continue
		}

		choice := strings.TrimSpace(input)

		switch choice {
		case "1":
			err := fetchManifest(reader)
			if err != nil {
				log.Println("Manifest fetching failed:", err)
			}
		case "2":
			err := uploadFile(reader)
			if err != nil {
				log.Println("File upload failed:", err)
			}
		case "3":
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid choice. Please enter 1, 2, or 3.")
		}
	}

}

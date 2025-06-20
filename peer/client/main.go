package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalln("Failed to dial udp with: ", err)
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	peerURL := localAddr.IP.String()

	err = registerPeer(serverURL, peerURL)
	if err != nil {
		log.Fatalln("Failed to register peer with : ", err)
		return
	}
	fmt.Println("Client registered with peer URL:", peerURL)

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

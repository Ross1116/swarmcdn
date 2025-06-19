package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	err := fetchManifest(reader)
	if err != nil {
		log.Println("Manifest fetching failed with: ", err)
		return
	}

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Println("Failed to dial udp with: ", err)
		return
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.IPAddr)
	peerURL := localAddr.IP.String()

	err = registerPeer(serverURL, peerURL)
	if err != nil {
		log.Println("Failed to register peer with : ", err)
		return
	}
}

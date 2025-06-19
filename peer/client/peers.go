package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func loadPeerList() []string {
	file := PeersFile
	data, err := os.ReadFile(file)
	if err != nil {
		log.Printf("Could not read peers.json: %v", err)
		return nil
	}

	var peers []string
	if err := json.Unmarshal(data, &peers); err != nil {
		log.Printf("Invalid peers.json format: %v", err)
		return nil
	}

	unique := make(map[string]bool)
	var uniquePeers []string
	for _, peer := range peers {
		peer = strings.TrimSpace(peer)
		peer = strings.TrimRight(peer, "/")
		if peer == "" || unique[peer] {
			continue
		}
		unique[peer] = true
		uniquePeers = append(uniquePeers, peer)
	}

	return uniquePeers
}

func updatePeerList(serverURL, peersFilePath string) error {
	resp, err := http.Get(fmt.Sprintf("%s/peers", serverURL))
	if err != nil {
		return fmt.Errorf("failed to fetch peers: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error fetching peers: %s - %s", resp.Status, string(body))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read peers response: %v", err)
	}

	if err := os.WriteFile(peersFilePath, content, 0644); err != nil {
		return fmt.Errorf("failed to update peers file: %v", err)
	}

	return err
}

func registerPeer(serverURL, myPeerURL string) error {
	payload := map[string]string{"url": myPeerURL}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal url JSON %v:", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/peers/register", serverURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response data: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error: %s - %s", resp.Status, string(body))
	}

	log.Printf("Successfully register peer: %s", string(body))
	err = updatePeerList(serverURL, PeersFile)
	if err != nil {
		return fmt.Errorf("Warning: Peer registered but failed to update peer list: %v", err")
	}
	return err
}

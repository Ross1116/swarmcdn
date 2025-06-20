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

// loadPeerList loads peers from peers.json and ensures the central server is included as fallback.
func loadPeerList() []string {
	data, err := os.ReadFile(PeersFile)
	if err != nil {
		log.Printf("Could not read peers.json: %v", err)
		return []string{strings.TrimRight(serverURL, "/")}
	}

	var peers []string
	if err := json.Unmarshal(data, &peers); err != nil {
		log.Printf("Invalid peers.json format: %v", err)
		return []string{strings.TrimRight(serverURL, "/")}
	}

	server := strings.TrimRight(serverURL, "/")
	seen := make(map[string]bool)
	var uniquePeers []string

	for _, peer := range peers {
		peer = strings.TrimSpace(strings.TrimRight(peer, "/"))
		if peer == "" || seen[peer] {
			continue
		}
		seen[peer] = true
		uniquePeers = append(uniquePeers, peer)
	}

	if !seen[server] {
		uniquePeers = append(uniquePeers, server)
	}

	return uniquePeers
}

// updatePeerList fetches updated peer list from the central server and writes it to peers.json.
func updatePeerList() error {
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

	if err := os.WriteFile(PeersFile, content, 0644); err != nil {
		return fmt.Errorf("failed to write peers.json: %v", err)
	}

	return nil
}

func registerPeer(myPeerURL string) error {
	payload := map[string]string{"url": myPeerURL}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal peer registration payload: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/peers/register", serverURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send registration request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read registration response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed: %s - %s", resp.Status, string(body))
	}

	log.Printf("Successfully registered peer: %s", string(body))

	if err := updatePeerList(); err != nil {
		return fmt.Errorf("peer registered but failed to update peer list: %v", err)
	}

	return nil
}

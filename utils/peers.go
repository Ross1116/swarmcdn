package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func LoadPeerList() ([]string, error) {
	file := GetPeersFilePath()
	data, err := os.ReadFile(file)
	if err != nil {
		log.Printf("Could not read peers.json: %v\n", err)
		return nil, err
	}

	var peers []string
	if err := json.Unmarshal(data, &peers); err != nil {
		log.Printf("Invalid peers.json format: %v\n", err)
		return nil, err
	}

	unique := make(map[string]bool)
	var uniquePeers []string
	for _, peer := range peers {
		peer = strings.TrimSpace(strings.TrimRight(peer, "/"))
		if peer == "" || unique[peer] {
			continue
		}
		unique[peer] = true
		uniquePeers = append(uniquePeers, peer)
	}

	return uniquePeers, nil
}

func SavePeers(peerList []string) error {
	updatedContent, err := json.MarshalIndent(peerList, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal the peer list: %w", err)
	}

	if err := os.WriteFile(GetPeersFilePath(), updatedContent, 0644); err != nil {
		return fmt.Errorf("failed to write peer file: %w", err)
	}

	return nil
}

func DeletePeer(peerURL string) error {
	peerURL = strings.TrimSpace(strings.TrimRight(peerURL, "/"))

	data, err := LoadPeerList()
	if err != nil {
		log.Printf("Unable to load peer list: %v\n", err)
		return err
	}

	var updatedPeers []string
	for _, peer := range data {
		if peer != peerURL {
			updatedPeers = append(updatedPeers, peer)
		}
	}

	if err := SavePeers(updatedPeers); err != nil {
		log.Printf("Error saving updated peer list: %v\n", err)
		return err
	}

	log.Printf("Successfully deleted peer: %s\n", peerURL)
	return nil
}

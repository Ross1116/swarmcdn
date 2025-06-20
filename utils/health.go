package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	HealthPath      = "/health"
	CheckInterval   = 100 * time.Second
	TimeoutInterval = 3 * time.Second
)

func CheckHealth(baseURL string) {
	ctx, cancel := context.WithTimeout(context.Background(), TimeoutInterval)
	defer cancel()

	url := baseURL + HealthPath

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request for %s with err %s\n", baseURL, err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error checking health for %s: %v\n", baseURL, err)
		DeletePeer(baseURL)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Health check failed for %s: status=%s, body=%s\n", baseURL, resp.Status, string(body))
		DeletePeer(baseURL)
		return
	}

	fmt.Printf("Health check performed successfully for %s\n", baseURL)
}

func CheckHealthPeriodic() {
	ticker := time.NewTicker(CheckInterval)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				peers, err := LoadPeerList()
				if err != nil {
					fmt.Println("Error loading peers:", err)
					continue
				}
				for _, peer := range peers {
					go CheckHealth(peer)
				}
			}
		}
	}()

	select {}
}

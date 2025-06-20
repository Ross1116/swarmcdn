package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func uploadFile(reader *bufio.Reader) error {
	fmt.Print("Enter file path to upload: ")
	pathInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read file path: %v", err)
	}
	path := strings.TrimSpace(pathInput)

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file data: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/upload", serverURL),
		writer.FormDataContentType(),
		&requestBody,
	)
	if err != nil {
		return fmt.Errorf("Unable to post the file to the server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch manifest. Status: %s\n", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch manifest: %s\n%s", resp.Status, string(body))
	}

	fmt.Println("File uploaded successfully")
	return nil
}

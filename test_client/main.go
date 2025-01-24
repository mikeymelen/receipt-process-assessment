package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	files := []string{
		"../examples/morning-receipt.json",
		"../examples/simple-receipt.json",
	}

	var ids []string

    // Send POST requests
    for _, file := range files {
        data, err := os.ReadFile(file)
        if err != nil {
            fmt.Printf("Failed to read file %s: %v\n", file, err)
            continue
        }

        resp, err := http.Post("http://localhost:8080/receipts/process", "application/json", bytes.NewBuffer(data))
        if err != nil {
            fmt.Printf("Failed to send request for file %s: %v\n", file, err)
            continue
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("Failed to read response for file %s: %v\n", file, err)
            continue
        }

        var result map[string]string
        if err := json.Unmarshal(body, &result); err != nil {
            fmt.Printf("Failed to parse response for file %s: %v\n", file, err)
            continue
        }

        if id, exists := result["id"]; exists {
            ids = append(ids, id)
            fmt.Printf("Received ID for file %s: %s\n", filepath.Base(file), id)
        } else {
            fmt.Printf("ID not found in response for file %s\n", file)
        }
    }

    // Send GET requests
    for _, id := range ids {
        url := fmt.Sprintf("http://localhost:8080/receipts/%s/points", id)
        resp, err := http.Get(url)
        if err != nil {
            fmt.Printf("Failed to get points for ID %s: %v\n", id, err)
            continue
        }
        defer resp.Body.Close()

        body, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("Failed to read response for ID %s: %v\n", id, err)
            continue
        }

        var result map[string]int
        if err := json.Unmarshal(body, &result); err != nil {
            fmt.Printf("Failed to parse response for ID %s: %v\n", id, err)
            continue
        }

        if points, exists := result["points"]; exists {
            fmt.Printf("Points for ID %s: %d\n", id, points)
        } else {
            fmt.Printf("Points not found in response for ID %s\n", id)
        }
    }
}

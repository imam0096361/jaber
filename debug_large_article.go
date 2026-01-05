package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	url := "http://localhost:9999/api/add-article"
	
	// Create a large content string (approx 2MB)
	largeContent := strings.Repeat("This is a test sentence to generate volume. ", 50000)
	fmt.Printf("Generated content size: %d bytes\n", len(largeContent))

	payload := map[string]interface{}{
		"title":    "Large Article Test",
		"content":  largeContent,
		"category": "প্রযুক্তি",
		"author":   "Debug Bot",
		"featured": false,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body Length: %d\n", len(body))
	// fmt.Printf("Body: %s\n", string(body))
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// 1. Upload Image
	// Create a dummy image file for testing
	// In a real scenario we would read a file, here we mock the multipart request... 
	// actually for simplicity let's just send an article with a hardcoded image path first 
	// to see if "image": "string" works. 
    // If the user says "news dite parteci na", they might fail at image upload too.

	url := "http://localhost:9999/api/add-article"
	
	imagePath := "/frontend/uploads/debug_image.jpg"
	
	payload := map[string]interface{}{
		"title":    "Debug Article With Image",
		"content":  "Testing with image pointer.",
		"category": "Technology",
		"author":   "Debug Bot",
		"featured": true,
		"image":    imagePath,
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
	fmt.Printf("Body: %s\n", string(body))
}

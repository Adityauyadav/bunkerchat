package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ChatHistoryResponse struct {
	Messages []struct {
		From      string    `json:"from"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"messages"`
}

func downloadChat(username string, recipient string, token string) error {
	req, err := http.NewRequest("GET", ServerUrl+"/chat/"+recipient, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch chat history: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get chat history: status %d, %s", resp.StatusCode, string(body))
	}

	var chatResp ChatHistoryResponse
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	if err != nil {
		return fmt.Errorf("failed to decode chat response: %w", err)
	}

	filename := fmt.Sprintf("chat_%s_%d.txt", recipient, time.Now().Unix())
	content := fmt.Sprintf("=== Chat with %s ===\n", recipient)
	content += fmt.Sprintf("Downloaded: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	content += "---\n\n"

	for _, msg := range chatResp.Messages {
		timeStr := msg.CreatedAt.Format("15:04:05")
		content += fmt.Sprintf("[%s] %s: %s\n", timeStr, msg.From, msg.Content)
	}

	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write chat file: %w", err)
	}

	fmt.Printf("Chat saved to %s\n", filename)
	return nil
}

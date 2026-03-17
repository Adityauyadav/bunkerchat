package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func login(username string, password string) (string, error) {
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login payload: %w", err)
	}

	resp, err := http.Post(ServerUrl+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed: invalid credentials")
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	token, exists := response["token"]
	if !exists {
		return "", fmt.Errorf("no token in login response")
	}

	return token, nil
}

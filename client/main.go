package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

const Version = "1.0.0"
const ServerUrl = "http://localhost:8080"
const WSServerUrl = "ws://localhost:8080"

func main() {
	fmt.Println("Welcome to Bunker Chat v" + Version)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your username: ")
	scanner.Scan()
	username := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}


	fmt.Print("Enter your password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	password := string(passwordBytes)
	fmt.Println()

	token, err := login(username, password)
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Printf("✓ Logged in successfully!\n")
	fmt.Printf("Welcome back, %s!\n", username)
	fmt.Println()


	fmt.Print("Who do you want to chat with? (username): ")
	scanner.Scan()
	recipient := scanner.Text()
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
		os.Exit(1)
	}

	if recipient == "" {
		fmt.Println("Recipient cannot be empty")
		os.Exit(1)
	}

	err = wsConnect(token, recipient, username)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/adityauyadav/bunkerchat/models"
	"github.com/gorilla/websocket"
)

func wsConnect(token string, recipient string, username string) error {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/ws",
	}

	q := u.Query()
	q.Set("token", token)
	q.Set("recipient", recipient)
	u.RawQuery = q.Encode()

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	fmt.Printf("\n--- Connected to %s ---\n", recipient)
	fmt.Println("Type your message and press Enter to send. Press Ctrl+C to exit.")
	fmt.Println()

	done := make(chan struct{})
	var closeOnce sync.Once
	closeDone := func() {
		closeOnce.Do(func() {
			close(done)
		})
	}

	go readFromServer(ws, done, recipient, closeDone)
	go sendFromTerminal(ws, done, recipient, closeDone)

	go handleSignals(done, closeDone)

	<-done
	ws.Close()

	fmt.Println("\n--- Chat ended ---")
	fmt.Println("Session ended. Goodbye.")
	return nil
}

func readFromServer(ws *websocket.Conn, done chan struct{}, recipient string, closeDone func()) {
	for {
		select {
		case <-done:
			return
		default:
		}

		var msg models.MessagePacket
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("\nConnection closed: %v\n", err)
			}
			closeDone()
			return
		}

		fmt.Printf("%s: %s\n", recipient, msg.Content)
	}
}

func sendFromTerminal(ws *websocket.Conn, done chan struct{}, recipient string, closeDone func()) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		select {
		case <-done:
			return
		default:
		}

		message := scanner.Text()
		if message == "" {
			continue
		}

		packet := models.MessagePacket{
			To:      recipient,
			Content: message,
		}

		err := ws.WriteJSON(packet)
		if err != nil {
			fmt.Printf("\nFailed to send message: %v\n", err)
			closeDone()
			return
		}
	}

	closeDone()
}

func handleSignals(done chan struct{}, closeDone func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n\nReceived interrupt signal, disconnecting...")
	closeDone()
}

func promptDownload() bool {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Download chat log? (y/n): ")
		scanner.Scan()
		response := scanner.Text()

		if response == "y" || response == "Y" {
			return true
		} else if response == "n" || response == "N" {
			return false
		}

		fmt.Println("Please enter 'y' or 'n'")
	}
}

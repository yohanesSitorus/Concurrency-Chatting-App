package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	conn     net.Conn
	username string
}

var (
	clients = make(map[net.Conn]*Client)
	mtx     sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Server is listening on port 9090...")
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	client := &Client{
		conn: conn,
	}

	reader := bufio.NewReader(conn)

	// Loop sampai username yang dimasukkan client unik
	for {
		// Membaca username terlebih dahulu
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read username: %v\n", err)
			return
		}

		parts := strings.SplitN(strings.TrimSpace(message), ":", 2)
		if len(parts) == 2 && parts[0] == "USERNAME" {
			//Cek apakah sudah ada client yang memiliki username tersebut
			username := parts[1]

			mtx.Lock()

			isTaken := false
			for _, exClient := range clients {
				if exClient.username == username {
					isTaken = true
					break
				}
			}

			mtx.Unlock()

			if isTaken {
				conn.Write([]byte("Username " + username + " already taken. Please try again.\n"))
			} else {
				//Jika username sudah unik
				client.username = parts[1]

				mtx.Lock()
				clients[conn] = client
				mtx.Unlock()

				conn.Write([]byte("OK\n"))
				fmt.Printf("Client %s connected as %s\n", conn.RemoteAddr(), client.username)
				broadcastMessage(conn, fmt.Sprintf("📢 %s has joined the chat!\n", client.username))
				break
			}
		}
	}

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			mtx.Lock()
			delete(clients, conn)
			mtx.Unlock()
			broadcastMessage(conn, fmt.Sprintf("📢 %s has left the chat!\n", client.username))
			return
		}

		parts := strings.SplitN(strings.TrimSpace(message), ":", 3)
		if len(parts) == 3 && parts[0] == "MSG" {
			username := parts[1]
			content := parts[2]

			fmt.Printf("[%s] %s: %s\n", conn.RemoteAddr(), username, content)
			// Broadcast pesan hanya ke client lain
			broadcastMessage(conn, fmt.Sprintf("%s: %s\n", username, content))
		}
	}
}

func broadcastMessage(sender net.Conn, message string) {
	mtx.Lock()
	defer mtx.Unlock()

	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}

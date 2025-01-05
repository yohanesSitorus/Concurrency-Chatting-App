package main

import (
	"Chat-App/server/room"
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
	room     *room.Room
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
			// Cek apakah sudah ada client yang memiliki username tersebut
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
				// Jika username sudah unik
				client.username = username

				mtx.Lock()
				clients[conn] = client
				mtx.Unlock()

				conn.Write([]byte("OK\n"))
				fmt.Printf("Client %s connected as %s\n", conn.RemoteAddr(), client.username)
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
			if client.room != nil {
				client.room.Broadcast(client.conn, fmt.Sprintf("📢 %s has left the room.\n", client.username))
				client.room.RemoveClient(client.conn)
			}
			return
		}

		parts := strings.SplitN(strings.TrimSpace(message), ":", 2)

		if len(parts) == 2 && parts[0] == "JOIN" {
			// Command untuk join ke room
			roomName := parts[1]

			if client.room != nil {
				client.room.RemoveClient(client.conn)
			}

			newRoom := room.GetOrCreateRoom(roomName)
			newRoom.AddClient(client.conn)
			client.room = newRoom

			conn.Write([]byte("Joined room " + roomName + "\n"))
			fmt.Printf("Client %s has joined room %s\n", client.username, roomName)
			newRoom.Broadcast(conn, fmt.Sprintf("📢 %s has joined the room.\n", client.username))
		} else if len(parts) == 2 && parts[0] == "MSG" {
			// Command untuk mengirim pesan ke room
			if client.room != nil {
				client.room.Broadcast(conn, fmt.Sprintf("%s: %s\n", client.username, parts[1]))
			} else {
				conn.Write([]byte("You are not in a room. Join a room first using the command: JOIN:<room_name>\n"))
			}
		} else if len(parts) == 1 && parts[0] == "LEAVE" {
			// Command untuk meninggalkan room
			if client.room != nil {
				client.room.Broadcast(conn, fmt.Sprintf("📢 %s has left the room.\n", client.username))
				client.room.RemoveClient(client.conn)
				client.room = nil
				conn.Write([]byte("You have left the room.\n"))
			} else {
				conn.Write([]byte("You are not in any room to leave.\n"))
			}
		} else if len(parts) == 1 && parts[0] == "EXIT" {
			// Command untuk keluar dari server
			mtx.Lock()
			delete(clients, conn)
			mtx.Unlock()

			if client.room != nil {
				client.room.Broadcast(client.conn, fmt.Sprintf("📢 %s has left the room.\n", client.username))
				client.room.RemoveClient(client.conn)
			}

			conn.Write([]byte("Goodbye!\n"))
			fmt.Printf("Client %s disconnected.\n", client.username)
			return
		}
	}
}

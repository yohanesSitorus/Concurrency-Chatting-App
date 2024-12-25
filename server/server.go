package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string) // Menyimpan semua koneksi client
	mtx     sync.Mutex                  // Mutex untuk akses aman ke map clients
)

func main() {
	// Membuat server untuk mendengarkan di port 9000
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to listen: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Server is listening on port 9090...")
	}
	defer ln.Close()

	for {
		// Menerima koneksi baru
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}

		// Menangani koneksi client dalam goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Menambahkan client ke daftar
	mtx.Lock()
	clients[conn] = conn.RemoteAddr().String()
	mtx.Unlock()

	fmt.Printf("Client %s connected.\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	fmt.Println("Waiting for message...")
	for {
		// Membaca pesan dari client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read message: %v\n", err)
			return

		} else {
			fmt.Println("The message has been received!")
		}

		// Membersihkan dan mencetak pesan
		message = strings.TrimSpace(message)
		fmt.Printf("%s: %s\n", conn.RemoteAddr(), message)
		fmt.Println("The message has been echoed back!")

		// Menyebarkan pesan ke semua client kecuali pengirim
		broadcastMessage(conn, message)
	}
}

func broadcastMessage(sender net.Conn, message string) {
	mtx.Lock()
	defer mtx.Unlock()

	for client := range clients {
		if client != sender {
			client.Write([]byte(fmt.Sprintf("%s: %s\n", clients[sender], message)))
		}
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Menghubungkan client ke server
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	// Meminta username dari pengguna
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading username: %v\n", err)
		return
	}
	username = strings.TrimSpace(username)

	// Mengirim username ke server
	conn.Write([]byte("USERNAME:" + username + "\n"))

	fmt.Printf("Connected to server as %s.\n", username)

	// Goroutine untuk mendengarkan pesan dari server
	go func() {
		connReader := bufio.NewReader(conn)
		for {
			message, err := connReader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot read the message: %v\n", err)
				return
			}
			// Clear line saat ini, print pesan, kemudian print prompt lagi
			fmt.Print("\r\033[K") // Clear line
			fmt.Print(message)
			fmt.Print("You: ")
		}
	}()

	// Loop untuk membaca input pengguna dengan prompt
	for {
		fmt.Print("You: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break
		}

		message = strings.TrimSpace(message)
		if message == "exit" {
			fmt.Println("Exiting chat...")
			break
		}

		// Kirim pesan ke server
		conn.Write([]byte("MSG:" + username + ":" + message + "\n"))
	}
}

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

	reader := bufio.NewReader(os.Stdin)

	var username string

	// Meminta username dari pengguna
	for {
		fmt.Print("Enter your username: ")
		username, err = reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading username: %v\n", err)
			return
		}
		username = strings.TrimSpace(username)

		// Mengirim username ke server
		_, err = conn.Write([]byte("USERNAME:" + username + "\n"))
		if err != nil {
			fmt.Fprint(os.Stderr, "Error sending username: %v\n", err)
			return
		}

		// Lihat respons server
		serverResponse, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading server response: %v\n", err)
			return
		}
		serverResponse = strings.TrimSpace(serverResponse)

		if serverResponse == "OK" {
			fmt.Printf("Connected to server as %s.\n", username)
			break
		} else {
			fmt.Println(serverResponse)
		}
	}

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

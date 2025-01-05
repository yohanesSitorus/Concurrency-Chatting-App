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
			fmt.Fprintf(os.Stderr, "Error sending username: %v\n", err)
			return
		}

		// Lihat respons server
		serverResponse, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading server response: %v\n", err)
			return
		}
		serverResponse = strings.TrimSpace(serverResponse)

		if serverResponse == "OK" {
			fmt.Printf("Welcome to the server, %s.\n", username)
			fmt.Println("Type HELP for list of commands")
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

	var inRoom bool = false // Flag untuk track apakah seorang user berada dalam room atau tidak

	// Loop untuk membaca input pengguna dengan prompt
	for {
		fmt.Print("You: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			break
		}

		message = strings.TrimSpace(message)

		// Untuk EXIT dari sebuah room
		if message == "EXIT" {
			conn.Write([]byte("EXIT\n"))
			fmt.Println("Exiting chat...")
			break
		}

		// Untuk LEAVE dari sebuah room
		if message == "LEAVE" {
			if inRoom {
				conn.Write([]byte("LEAVE\n"))
				inRoom = false
				fmt.Println("You are not in a room. Join a room first using the command: JOIN:<room_name>")
			} else {
				fmt.Println("You are not in a room.")
			}
			continue
		}

		// Untuk mendapatkan list command yang ada
		if message == "HELP" {
			help()
			continue
		}

		// Untuk JOIN sebuah room
		if strings.HasPrefix(message, "JOIN:") {
			conn.Write([]byte(message + "\n"))
			inRoom = true
			continue
		}

		// Send message: hanya untuk user yang berada dalam room
		if inRoom {
			conn.Write([]byte("MSG:" + message + "\n"))
		} else {
			fmt.Println("You are not in a room. Join a room first using the command: JOIN:<room_name>")
		}
	}
}

// Fungsi untuk print sebuah command yang ada
func help(){
	fmt.Println("============================================")
	fmt.Println("JOIN:<room_name>  		--- 	Join new or existing room")
	fmt.Println("LEAVE 				--- 	Leave current room")
	fmt.Println("EXIT 				--- 	Disconnect from current server")
	fmt.Println("HELP				--- 	List of All Commands")
	fmt.Println("============================================")
}
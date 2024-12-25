package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Menghubungkan client ke server
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to server: %v\n", err)
	} else {
		fmt.Println("Connected to server. Type your message and press Enter to send.")
	}

	defer conn.Close()

	// Goroutine untuk mendengarkan pesan dari server
	go func() {
		connReader := bufio.NewReader(conn)
		for {
			message, err := connReader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot read the message: %v\n", err)
				return
			} else {
				fmt.Print(">> " + message)
			}
		}
	}()

	// Loop untuk membaca input pengguna dan mengirimkannya ke server
	localReader := bufio.NewScanner(os.Stdin)
	for localReader.Scan() {
		message := localReader.Text()
		if message == "exit" {
			fmt.Println("Exiting chat...")
			break
		}
		conn.Write([]byte(message + "\n"))
	}
}

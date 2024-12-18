package main

import (
	// "bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp", ":9090")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to listen!")
        os.Exit(1)
    } else {
        fmt.Println("Listening...")
    }

    fmt.Println(ln)
    // conn, err := ln.Accept()
    // if err != nil {
    //     fmt.Fprintf(os.Stderr, "Failed to accept connection!")
    //     os.Exit(1)
    // } else {
    //     fmt.Println("New connection accepted!")
    // }
}
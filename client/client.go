package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	address := "localhost:8080"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("❌ Erro ao conectar: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("🎮 Conectado ao LuthiBOX em %s\n", address)

	// Handle server messages
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Send user input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text() + "\n"
		conn.Write([]byte(input))
	}
}

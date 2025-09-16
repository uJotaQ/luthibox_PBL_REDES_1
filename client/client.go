package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// Host padrão: nome do serviço no docker-compose
	address := "server:8080"
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	// Tentar conectar várias vezes caso o server ainda não esteja pronto
	var conn net.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			break
		}
		fmt.Println("⏳ Tentando conectar ao server...")
		time.Sleep(1 * time.Second)
	}

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
			msg := scanner.Text()
			fmt.Println(msg)
		}
	}()

	// Send user input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		conn.Write([]byte(input + "\n"))
	}
}

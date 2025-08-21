package main

import (
	"bufio"
	"fmt"
	"net"
)

var clients = make(map[string]net.Conn)

func main() {
	// Escuta na porta 8080 no localhost
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Servidor TCP rodando em localhost:8080")

	for {
		// Aceita novas conexões
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		fmt.Println("Novo cliente conectado:", conn.RemoteAddr())

		conn.Write([]byte("Qual seu nome?\n"))

		// Roda cada cliente em uma goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		// Lê mensagem do cliente
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectado:", conn.RemoteAddr())
			return
		}
		fmt.Print("Mensagem recebida: ", message)

		// Responde ao cliente
		response := "Servidor recebeu: " + message
		conn.Write([]byte(response))
	}
}

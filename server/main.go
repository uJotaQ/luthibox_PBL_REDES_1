package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string) // conexão -> nome de usuário
	clientsMu sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Servidor rodando na porta 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Perguntar o nome
	conn.Write([]byte("Digite seu nome:\n"))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Guardar o cliente no map
	clientsMu.Lock()
	clients[conn] = username
	clientsMu.Unlock()

	// Avisar no chat que entrou
	broadcast(fmt.Sprintf("%s entrou no chat\n", username), conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		// Enviar para todos com o nome do autor
		broadcast(fmt.Sprintf("[%s]: %s\n", username, msg), conn)
	}

	// Remover cliente quando sair
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	broadcast(fmt.Sprintf("%s saiu do chat\n", username), conn)
}

func broadcast(message string, sender net.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for conn := range clients {
		if conn != sender { // não mandar de volta para quem enviou
			conn.Write([]byte(message))
		}
	}
}

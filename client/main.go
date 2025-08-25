package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Conecta no servidor
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Erro ao conectar:", err)
		return
	}
	defer conn.Close()

	// Goroutine para ler mensagens do servidor
	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Servidor desconectado")
				os.Exit(0)
			}
			fmt.Print(msg) // imprime imediatamente a linha do servidor
		}
	}()

	// Lê do teclado e envia para o servidor
	inputReader := bufio.NewReader(os.Stdin)
	for {
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		conn.Write([]byte(text + "\n"))
	}
}

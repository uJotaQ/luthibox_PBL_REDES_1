package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	fmt.Println("🧪 Teste de Stress de Conexão")
	fmt.Println("🎯 Múltiplos clientes conectando, registrando e desconectando simultaneamente")

	const numClients = 10000
	const serverAddr = "localhost:8080"

	var wg sync.WaitGroup

	fmt.Printf("🚀 Iniciando %d clientes simultaneamente...\n", numClients)

	// Iniciar todos os clientes ao mesmo tempo
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()
			runClientTest(clientID, serverAddr)
		}(i)
	}

	// Esperar todos terminarem
	wg.Wait()

	fmt.Println("\n✅ Teste de Stress de Conexão Concluído!")
}

func runClientTest(clientID int, serverAddr string) {
	nickname := fmt.Sprintf("stress_%d", clientID)
	password := "pass123"

	// Conectar ao servidor
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("❌ Cliente %d: Erro ao conectar - %v\n", clientID, err)
		return
	}

	fmt.Printf("👤 Cliente %d: Conectado\n", clientID)

	// Enviar registro
	registerCmd := fmt.Sprintf("/register %s %s\n", nickname, password)
	conn.Write([]byte(registerCmd))

	fmt.Printf("📝 Cliente %d: Registro enviado\n", clientID)

	// Esperar um tempo aleatório (simulando uso)
	//waitTime := time.Duration(1 + clientID*100) // 1-2.5 segundos
	//time.Sleep(waitTime * time.Millisecond)

	// Desconectar
	conn.Close()
	fmt.Printf("👋 Cliente %d: Desconectado\n", clientID)
}

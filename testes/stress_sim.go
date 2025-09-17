package main

import (
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("🧪 Teste de Stress de Conexão Melhorado")
	fmt.Println("🎯 Múltiplos clientes conectando, registrando e desconectando simultaneamente")

	const numClients = 40000
	const serverAddr = "localhost:8080"
	const maxConcurrent = 250 // Limite de conexões simultâneas

	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

	// Contadores de sucesso e falha
	var connectedCount int32
	var registeredCount int32
	var failedConnections int32

	fmt.Printf("🚀 Iniciando %d clientes simultaneamente (máx %d concorrentes)...\n", numClients, maxConcurrent)

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			// Controle de concorrência
			sem <- struct{}{}
			defer func() { <-sem }()

			// Pequeno delay aleatório para espalhar conexões
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))

			if runClientTest(clientID, serverAddr) {
				atomic.AddInt32(&connectedCount, 1)
				atomic.AddInt32(&registeredCount, 1)
			} else {
				atomic.AddInt32(&failedConnections, 1)
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("\n✅ Teste de Stress de Conexão Concluído!")
	fmt.Printf("📊 Relatório:\n")
	fmt.Printf("   Clientes conectados com sucesso: %d\n", connectedCount)
	fmt.Printf("   Clientes registrados com sucesso: %d\n", registeredCount)
	fmt.Printf("   Clientes que falharam ao conectar: %d\n", failedConnections)
}

func runClientTest(clientID int, serverAddr string) bool {
	nickname := fmt.Sprintf("stress_%d", clientID)
	password := "pass123"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Printf("❌ Cliente %d: Erro ao conectar - %v\n", clientID, err)
		return false
	}
	defer conn.Close()

	fmt.Printf("👤 Cliente %d: Conectado\n", clientID)

	// Enviar registro
	registerCmd := fmt.Sprintf("/register %s %s\n", nickname, password)
	_, err = conn.Write([]byte(registerCmd))
	if err != nil {
		fmt.Printf("❌ Cliente %d: Erro ao enviar registro - %v\n", clientID, err)
		return false
	}

	fmt.Printf("📝 Cliente %d: Registro enviado\n", clientID)
	return true
}

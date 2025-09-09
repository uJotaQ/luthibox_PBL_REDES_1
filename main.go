package main

import (
	"fmt"
	"luthibox/network"
	"os"
)

func main() {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	server := network.NewServer(port)

	fmt.Println("🚀 Iniciando LuthiBOX Server...")
	if err := server.Start(); err != nil {
		fmt.Printf("❌ Erro fatal: %v\n", err)
		os.Exit(1)
	}
}

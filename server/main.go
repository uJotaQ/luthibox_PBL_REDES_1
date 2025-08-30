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

type Ataque struct {
	Nome      string
	Sequencia []string // sequência de notas (ex: {"A", "B", "G", "D"})
}

type Instrumento struct {
	Nome     string
	Raridade string
	Ataques  [3]Ataque // cada instrumento terá 3 ataques fixos
}

// Banco de dados interno
var BancoInstrumentos []Instrumento

// Função que cria um instrumento
func CriarInstrumento(nome, raridade string, ataques [3]Ataque) Instrumento {
	return Instrumento{
		Nome:     nome,
		Raridade: raridade,
		Ataques:  ataques,
	}
}

// Função que popula o "banco"
func CarregarInstrumentos() {
	violino := CriarInstrumento("Violino", "Raro", [3]Ataque{
		{"Vibrato", []string{"A", "B", "G"}},
		{"Pizzicato", []string{"C", "E", "G"}},
		{"Spiccato", []string{"F", "A", "C"}},
	})

	guitarra := CriarInstrumento("Guitarra", "Épico", [3]Ataque{
		{"Power Chord", []string{"E", "A", "D"}},
		{"Solo Rápido", []string{"G", "B", "E"}},
		{"Palm Mute", []string{"E", "E", "A"}},
	})

	flauta := CriarInstrumento("Flauta", "Comum", [3]Ataque{
		{"Trinado", []string{"C", "D", "E"}},
		{"Sopro Forte", []string{"F", "G", "C"}},
		{"Melodia Doce", []string{"A", "B", "C"}},
	})

	piano := CriarInstrumento("Piano", "Épico", [3]Ataque{
		{"Arpejo", []string{"C", "E", "G"}},
		{"Glissando", []string{"D", "F", "A"}},
		{"Acorde Pesado", []string{"E", "G", "B"}},
	})

	trompete := CriarInstrumento("Trompete", "Raro", [3]Ataque{
		{"Susto Sonoro", []string{"C", "G", "C"}},
		{"Escalada", []string{"D", "F", "A"}},
		{"Sopro Metálico", []string{"E", "A", "B"}},
	})

	bateria := CriarInstrumento("Bateria", "Épico", [3]Ataque{
		{"Rufar", []string{"C", "C", "C"}},
		{"Virada", []string{"E", "G", "A"}},
		{"Prato Explosivo", []string{"F", "B", "E"}},
	})

	harpa := CriarInstrumento("Harpa", "Lendário", [3]Ataque{
		{"Arpejo Celestial", []string{"C", "F", "A"}},
		{"Glissando Suave", []string{"D", "G", "B"}},
		{"Eco Etéreo", []string{"E", "A", "C"}},
	})

	saxofone := CriarInstrumento("Saxofone", "Raro", [3]Ataque{
		{"Improviso Jazz", []string{"E", "G", "B"}},
		{"Sopro Grave", []string{"C", "E", "A"}},
		{"Nota Sustentada", []string{"D", "F", "G"}},
	})

	oboé := CriarInstrumento("Oboé", "Comum", [3]Ataque{
		{"Melodia Suave", []string{"C", "D", "E"}},
		{"Sopro Agudo", []string{"F", "A", "B"}},
		{"Eco Vibrante", []string{"G", "B", "C"}},
	})

	contrabaixo := CriarInstrumento("Contrabaixo", "Raro", [3]Ataque{
		{"Slap", []string{"E", "A", "D"}},
		{"Walking Bass", []string{"G", "B", "E"}},
		{"Grave Profundo", []string{"C", "E", "A"}},
	})

	banjo := CriarInstrumento("Banjo", "Comum", [3]Ataque{
		{"Dedilhado", []string{"C", "G", "B"}},
		{"Rasgueado", []string{"D", "A", "E"}},
		{"Bluegrass", []string{"F", "G", "C"}},
	})

	tambor := CriarInstrumento("Tambor", "Comum", [3]Ataque{
		{"Batida Rítmica", []string{"C", "C", "G"}},
		{"Toque Seco", []string{"D", "F", "A"}},
		{"Rufada Tribal", []string{"E", "A", "C"}},
	})

	acordeon := CriarInstrumento("Acordeon", "Raro", [3]Ataque{
		{"Expansão de Fole", []string{"C", "E", "A"}},
		{"Contra-Melodia", []string{"D", "G", "B"}},
		{"Dança Gaúcha", []string{"F", "A", "C"}},
	})

	ukulele := CriarInstrumento("Ukulele", "Comum", [3]Ataque{
		{"Acorde Leve", []string{"C", "E", "G"}},
		{"Ritmo Havaiano", []string{"D", "G", "A"}},
		{"Ponte Alegre", []string{"E", "A", "C"}},
	})

	// Adiciona todos ao banco
	BancoInstrumentos = append(
		BancoInstrumentos,
		violino, guitarra, flauta,
		piano, trompete, bateria, harpa,
		saxofone, oboé, contrabaixo, banjo,
		tambor, acordeon, ukulele,
	)
}

func main() {

	CarregarInstrumentos()

	for _, inst := range BancoInstrumentos {
		fmt.Println("Instrumento:", inst.Nome, "| Raridade:", inst.Raridade)
		for _, atk := range inst.Ataques {
			fmt.Println("  - Ataque:", atk.Nome, "| Notas:", atk.Sequencia)
		}
	}

	// ln, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	fmt.Println("Erro ao iniciar servidor:", err)
	// 	return
	// }
	// defer ln.Close()

	// fmt.Println("Servidor rodando na porta 8080...")

	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		fmt.Println("Erro ao aceitar conexão:", err)
	// 		continue
	// 	}

	// 	go handle(conn)
	// }
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

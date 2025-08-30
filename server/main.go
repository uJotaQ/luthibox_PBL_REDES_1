package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"math/rand"
    "time"
)

type Player struct {
    Nickname string
    Password string
    Conn     net.Conn
	Instrumentos []Instrumento
	Tokens int
}

var (
    players   = make(map[string]Player) // banco de "contas"
    playersMu sync.Mutex
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





type Pacote struct {
    ID        string
    Raridade  string
    Instrumento Instrumento
}

var BancoPacotes []Pacote
var pacotesMu sync.Mutex


func abrirPacotesFlow(player Player, raridade string) string {
    pacotesMu.Lock()
    defer pacotesMu.Unlock()

    // Filtra pacotes disponíveis da raridade desejada
    var pacotesDisponiveis []Pacote
    for _, p := range BancoPacotes {
        if p.Raridade == raridade {
            pacotesDisponiveis = append(pacotesDisponiveis, p)
        }
    }

    if len(pacotesDisponiveis) == 0 {
        return "Não há pacotes disponíveis dessa raridade.\n"
    }

    // Mostra até 10 pacotes com ID e preço
    resp := "Pacotes disponíveis:\n"
    max := 10
    if len(pacotesDisponiveis) < 10 {
        max = len(pacotesDisponiveis)
    }

    for i := 0; i < max; i++ {
        resp += fmt.Sprintf("ID: %s | Preço: %d tokens\n", pacotesDisponiveis[i].ID, pacotesDisponiveis[i].Instrumento.Preco)
    }

    resp += "Digite 'PEGAR_PACOTE <ID>' para abrir um pacote.\n"

    return resp
}

func GerarPacote(raridade string) Pacote {
    // Escolhe um instrumento aleatório da raridade desejada
    inst := pegarInstrumentoAleatorio(raridade)

    // Gera um ID único para o pacote (ex: A12, B45)
    letras := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    numeros := "0123456789"
    id := string(letras[rand.Intn(len(letras))]) +
        string(numeros[rand.Intn(len(numeros))]) +
        string(numeros[rand.Intn(len(numeros))])

    pacote := Pacote{
        ID:          id,
        Raridade:    raridade,
        Instrumento: inst,
    }

    return pacote
}



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

func pegarInstrumentoAleatorio(raridade string) Instrumento {
    rand.Seed(time.Now().UnixNano())

    // Filtra instrumentos da raridade escolhida
    var filtrados []Instrumento
    for _, inst := range BancoInstrumentos {
        if inst.Raridade == raridade {
            filtrados = append(filtrados, inst)
        }
    }

    // Caso não haja instrumentos da raridade, retorna um instrumento vazio
    if len(filtrados) == 0 {
        return Instrumento{}
    }

    // Sorteia aleatoriamente
    indice := rand.Intn(len(filtrados))
    return filtrados[indice]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	CarregarInstrumentos()
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
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	var currentPlayer Player

	// --- Loop de autenticação ---
	for {
		conn.Write([]byte("Digite '/login nick senha' ou '/register nick senha':\n"))
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectou:", conn.RemoteAddr())
			return
		}

		msg = strings.TrimSpace(msg)
		args := strings.Split(msg, " ")

		if len(args) != 3 {
			conn.Write([]byte("Comando inválido! Use '/login nick senha' ou '/register nick senha'\n"))
			continue
		}

		command, nick, pass := args[0], args[1], args[2]

		switch command {
		case "/login":
			playersMu.Lock()
			player, exists := players[nick]
			playersMu.Unlock()
			if !exists || player.Password != pass {
				conn.Write([]byte("Nickname ou senha incorretos\n"))
				continue
			}

			player.Conn = conn
			currentPlayer = player
			playersMu.Lock()
			players[nick] = player
			playersMu.Unlock()

			conn.Write([]byte("LOGIN_OK\n"))
			goto MENU

		case "/register":
			playersMu.Lock()
			_, exists := players[nick]
			if exists {
				playersMu.Unlock()
				conn.Write([]byte("Nickname já existe!\n"))
				continue
			}

			newPlayer := Player{
				Nickname:     nick,
				Password:     pass,
				Conn:         conn,
				Instrumentos: []Instrumento{},
				Tokens:       0,
			}
			players[nick] = newPlayer
			playersMu.Unlock()

			currentPlayer = newPlayer
			conn.Write([]byte("REGISTER_OK\n"))
			goto MENU

		default:
			conn.Write([]byte("Comando inválido! Use '/login' ou '/register'\n"))
		}
	}

MENU:
	// --- Loop de comandos do jogo ---
	for {
		cmdLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectou:", conn.RemoteAddr())
			return
		}
		cmdLine = strings.TrimSpace(cmdLine)
		args := strings.Split(cmdLine, " ")

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "JOGAR":
			// Aqui você chama sua função de iniciar partida
			conn.Write([]byte("JOGO_INICIADO\n"))

		case "ABRIR_PACOTES":
			if len(args) < 2 {
				conn.Write([]byte("ERRO: Informe a raridade do pacote\n"))
				continue
			}
			raridade := args[1]
			resp := abrirPacotesFlow(currentPlayer, raridade)
			conn.Write([]byte(resp))

		case "MEUS_INSTRUMENTOS":
			//resp := listarInstrumentos(currentPlayer)
			//conn.Write([]byte(resp))

		case "PEGAR_PACOTE":
			if len(args) < 2 {
				conn.Write([]byte("ERRO: Informe o ID do pacote\n"))
				continue
			}
			id := args[1]

			pacotesMu.Lock()
			var selecionado *Pacote
			for i, p := range BancoPacotes {
				if p.ID == id {
					selecionado = &p
					// Remove do estoque
					BancoPacotes = append(BancoPacotes[:i], BancoPacotes[i+1:]...)
					break
				}
			}
			pacotesMu.Unlock()

			if selecionado == nil {
				conn.Write([]byte("Pacote não encontrado ou já aberto!\n"))
				continue
			}

			// Adiciona instrumento ao player
			currentPlayer.Instrumentos = append(currentPlayer.Instrumentos, selecionado.Instrumento)

			// Repor novo pacote da mesma raridade
			novoPacote := GerarPacote(selecionado.Raridade)
			pacotesMu.Lock()
			BancoPacotes = append(BancoPacotes, novoPacote)
			pacotesMu.Unlock()

			conn.Write([]byte(fmt.Sprintf("Você abriu o pacote %s e recebeu o instrumento %s!\n", id, selecionado.Instrumento.Nome)))
		default:
			conn.Write([]byte("COMANDO_INVALIDO\n"))
		}
	}
}



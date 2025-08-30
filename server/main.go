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
    ID          string
    Raridade    string
    Instrumento Instrumento
}

var (
    pacotesDisponiveis []Pacote
    pacotesMu          sync.Mutex
)

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

func gerarID() string {
    letras := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
    numeros := []rune("0123456789")
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%c%c%c", letras[rand.Intn(len(letras))], numeros[rand.Intn(len(numeros))], numeros[rand.Intn(len(numeros))])
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

func InicializarPacotes() {
    raridades := []string{"Comum", "Raro", "Épico", "Lendário"}

    pacotesMu.Lock()
    defer pacotesMu.Unlock()

    for _, rar := range raridades {
        for i := 0; i < 5; i++ {
            instrumento := pegarInstrumentoAleatorio(rar) // função que sorteia um instrumento da raridade
            pacote := Pacote{
                ID:          gerarID(),
                Raridade:    rar,
                Instrumento: instrumento,
            }
            pacotesDisponiveis = append(pacotesDisponiveis, pacote)
        }
    }
}


func AbrirPacote(conn net.Conn, raridade string) {
    pacotesMu.Lock()
    defer pacotesMu.Unlock()

    // Filtra pacotes disponíveis da raridade
    var pacotesDaRaridade []Pacote
    for _, p := range pacotesDisponiveis {
        if p.Raridade == raridade {
            pacotesDaRaridade = append(pacotesDaRaridade, p)
        }
    }

    if len(pacotesDaRaridade) == 0 {
        conn.Write([]byte("Nenhum pacote disponível desta raridade.\n"))
        return
    }

    // Envia IDs dos pacotes para o client
    conn.Write([]byte("Pacotes disponíveis:\n"))
    for _, p := range pacotesDaRaridade {
        conn.Write([]byte(p.ID + "\n"))
    }

    // Espera escolha do client
    reader := bufio.NewReader(conn)
    escolhaID, _ := reader.ReadString('\n')
    escolhaID = strings.TrimSpace(escolhaID)

    // Procura pacote escolhido
    for i, p := range pacotesDisponiveis {
        if p.ID == escolhaID && p.Raridade == raridade {
            // Entrega instrumento
            msg := fmt.Sprintf("Você recebeu o instrumento: %s (Raridade: %s)\n", p.Instrumento.Nome, p.Instrumento.Raridade)
            conn.Write([]byte(msg))

            // Remove pacote e gera novo
            pacotesDisponiveis = append(pacotesDisponiveis[:i], pacotesDisponiveis[i+1:]...)
            novo := Pacote{
                ID:          gerarID(),
                Raridade:    raridade,
                Instrumento: pegarInstrumentoAleatorio(raridade),
            }
            pacotesDisponiveis = append(pacotesDisponiveis, novo)
            return
        }
    }

    conn.Write([]byte("Pacote inválido ou já foi aberto. Tente novamente.\n"))
}


func main() {
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

	// Loop de autenticação
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

		command := args[0]
		nick := args[1]
		pass := args[2]

		switch command {
		case "/login":
			playersMu.Lock()
			player, exists := players[nick]
			playersMu.Unlock()

			if !exists {
				conn.Write([]byte("Nickname não encontrado! Digite '/register nick senha' para criar conta\n"))
				continue
			}

			if player.Password != pass {
				conn.Write([]byte("Senha incorreta!\n"))
				continue
			}

			// Login bem-sucedido
			conn.Write([]byte(fmt.Sprintf("Login bem-sucedido! Bem-vindo, %s\n", nick)))
			player.Conn = conn
			playersMu.Lock()
			players[nick] = player
			playersMu.Unlock()
			fmt.Println(nick, "entrou no jogo.")
			return // sai do loop de login e vai pro menu do jogo

		case "/register":
			playersMu.Lock()
			_, exists := players[nick]
			if exists {
				conn.Write([]byte("Nickname já existe! Escolha outro.\n"))
				playersMu.Unlock()
				continue
			}

			// Cria novo player
			players[nick] = Player{
				Nickname: nick,
				Password: pass,
				Conn:     conn,
			}
			playersMu.Unlock()

			conn.Write([]byte(fmt.Sprintf("Registro bem-sucedido! Bem-vindo, %s\n", nick)))
			fmt.Println(nick, "se registrou e entrou no jogo.")
			return // sai do loop e vai pro menu do jogo

		default:
			conn.Write([]byte("Comando inválido! Use '/login' ou '/register'\n"))
		}
	}
}

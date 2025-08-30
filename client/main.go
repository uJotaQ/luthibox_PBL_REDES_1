package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func autenticar(conn net.Conn) {
    reader := bufio.NewReader(os.Stdin)
    serverReader := bufio.NewReader(conn)

    for {
        fmt.Println("Digite '/login nickname senha' ou '/register nickname senha':")
        input, _ := reader.ReadString('\n')
        input = strings.TrimSpace(input)

        // envia para o servidor
        conn.Write([]byte(input + "\n"))

        // recebe a resposta do server
        resposta, _ := serverReader.ReadString('\n')
        resposta = strings.TrimSpace(resposta)
        fmt.Println(resposta)

        // se login/register OK, sai do loop
        if resposta == "Login bem-sucedido" || resposta == "Registro bem-sucedido" {
            break
        }
    }
}

func abrirPacotesFlow(conn net.Conn) {
    reader := bufio.NewReader(os.Stdin)
    serverReader := bufio.NewReader(conn)

    // Pergunta a raridade
    fmt.Println("Escolha a raridade do pacote (Comum, Raro, Épico, Lendário):")
    raridade, _ := reader.ReadString('\n')
    raridade = strings.TrimSpace(raridade)

    // Envia raridade para o server
    conn.Write([]byte(raridade + "\n"))

    // Recebe lista de pacotes do server
    fmt.Println("Pacotes disponíveis:")
    for i := 0; i < 5; i++ { // supondo 5 pacotes por raridade
        linha, _ := serverReader.ReadString('\n')
        fmt.Print(linha)
    }

    // Pergunta qual pacote abrir
    fmt.Println("Digite o ID do pacote que deseja abrir:")
    escolhaID, _ := reader.ReadString('\n')
    escolhaID = strings.TrimSpace(escolhaID)
    conn.Write([]byte(escolhaID + "\n"))

    // Recebe o resultado do pacote
    resposta, _ := serverReader.ReadString('\n')
    fmt.Println(resposta)
}


func mostrarMenu(conn net.Conn) {
    for {
        fmt.Println("=== LuthiBOX ===")
        fmt.Println("1) Jogar")
        fmt.Println("2) Abrir Pacotes")
        fmt.Println("3) Meus Instrumentos")
        fmt.Print("Escolha uma opção: ")

        var opcao string
        fmt.Scanln(&opcao)
		opcao = strings.TrimSpace(opcao)

        switch opcao {
		case "0":
			fmt.Println("Saindo do jogo...")
			return
        case "1":
            conn.Write([]byte("JOGAR\n"))
        case "2":
            conn.Write([]byte("ABRIR_PACOTES\n"))
            abrirPacotesFlow(conn)
        case "3":
            conn.Write([]byte("MEUS_INSTRUMENTOS\n"))
        default:
            fmt.Println("Opção inválida")
        }
    }
}



func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Erro ao conectar no servidor:", err)
        return
    }
    defer conn.Close()

    // primeiro autenticação
    autenticar(conn)

    // só depois mostra o menu
    mostrarMenu(conn)
}

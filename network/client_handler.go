package network

import (
	"bufio"
	"fmt"
	"luthibox/game"
	"net"
	"strings"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("👤 Novo cliente conectado: %s\n", conn.RemoteAddr())

	// Authentication first
	player, err := authenticateClient(conn)
	if err != nil {
		fmt.Printf("Falha na autenticação: %v\n", err)
		return
	}

	fmt.Printf("✅ Jogador %s autenticado\n", player.Nickname)
	conn.Write([]byte(fmt.Sprintf("Bem-vindo ao LuthiBOX, %s!\n", player.Nickname)))

	// Main menu
	showMainMenu(player, conn)
}

func authenticateClient(conn net.Conn) (*game.Player, error) {
	reader := bufio.NewReader(conn)

	for {
		conn.Write([]byte("\n=== LUTHIBOX - AUTENTICAÇÃO ===\n"))
		conn.Write([]byte("Digite: /login <nickname> <senha>\n"))
		conn.Write([]byte("Ou:     /register <nickname> <senha>\n"))
		conn.Write([]byte("> "))

		msg, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		msg = strings.TrimSpace(msg)
		args := strings.Split(msg, " ")

		if len(args) != 3 {
			conn.Write([]byte("❌ Comando inválido!\n"))
			continue
		}

		command, nickname, password := args[0], args[1], args[2]

		switch command {
		case "/login":
			player, err := game.AuthenticatePlayer(nickname, password)
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("❌ %v\n", err)))
				continue
			}
			player.Conn = conn // Update connection
			return player, nil

		case "/register":
			err := game.RegisterPlayer(nickname, password, conn)
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("❌ %v\n", err)))
				continue
			}

			player, _ := game.GetPlayer(nickname)
			return player, nil

		default:
			conn.Write([]byte("❌ Comando inválido!\n"))
		}
	}
}

func showMainMenu(player *game.Player, conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		menu := `
🎮 === LUTHIBOX - MENU PRINCIPAL ===
1) 🎲 Jogar (Batalha 1v1)
2) 🎁 Abrir Pacotes
3) 🎵 Meus Instrumentos
4) 💰 Meus Tokens
0) 🚪 Sair

Escolha uma opção: `

		conn.Write([]byte(menu))

		opcao, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Cliente desconectado: %v\n", err)
			return
		}

		opcao = strings.TrimSpace(opcao)

		switch opcao {
		case "0":
			conn.Write([]byte("👋 Até logo!\n"))
			return

		case "1":
			startBattle(player, conn)

		case "2":
			openPackets(player, conn, reader)

		case "3":
			showInstruments(player, conn)

		case "4":
			showTokens(player, conn)

		default:
			conn.Write([]byte("❌ Opção inválida!\n"))
		}
	}
}

func startBattle(player *game.Player, conn net.Conn) {
	conn.Write([]byte("🔍 Procurando oponente...\n"))
	game.AddPlayerToBattleQueue(player)
	// The battle will be handled by the matchmaking system
}

func openPackets(player *game.Player, conn net.Conn, reader *bufio.Reader) {
	conn.Write([]byte("\n🎯 Escolha a raridade:\n"))
	conn.Write([]byte("1) Comum\n"))
	conn.Write([]byte("2) Raro\n"))
	conn.Write([]byte("3) Épico\n"))
	conn.Write([]byte("4) Lendário\n"))
	conn.Write([]byte("0) Voltar\n"))
	conn.Write([]byte("Raridade: "))

	raridade, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	raridade = strings.TrimSpace(raridade)
	rarityMap := map[string]string{
		"1": "Comum",
		"2": "Raro",
		"3": "Épico",
		"4": "Lendário",
	}

	if raridade == "0" {
		return
	}

	rarity, exists := rarityMap[raridade]
	if !exists {
		conn.Write([]byte("❌ Raridade inválida!\n"))
		return
	}

	// Check if player has enough tokens
	packetCost := map[string]int{
		"Comum":    10,
		"Raro":     25,
		"Épico":    50,
		"Lendário": 100,
	}

	cost := packetCost[rarity]
	if player.GetTokens() < cost {
		conn.Write([]byte(fmt.Sprintf("❌ Tokens insuficientes! Custa %d tokens.\n", cost)))
		return
	}

	// Show available packets
	packets := game.GetAvailablePacketsByRarity(rarity)
	if len(packets) == 0 {
		conn.Write([]byte("❌ Não há pacotes disponíveis dessa raridade.\n"))
		return
	}

	conn.Write([]byte(fmt.Sprintf("\n📦 Pacotes %s disponíveis:\n", rarity)))
	for i, packet := range packets {
		conn.Write([]byte(fmt.Sprintf("%d) ID: %s\n", i+1, packet.ID)))
	}

	conn.Write([]byte("\nDigite o número do pacote para abrir (0 para cancelar): "))
	escolha, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	escolha = strings.TrimSpace(escolha)
	if escolha == "0" {
		return
	}

	// Parse choice
	var packetIndex int
	fmt.Sscanf(escolha, "%d", &packetIndex)
	if packetIndex < 1 || packetIndex > len(packets) {
		conn.Write([]byte("❌ Opção inválida!\n"))
		return
	}

	selectedPacket := packets[packetIndex-1]

	// Remove tokens
	if !player.RemoveTokens(cost) {
		conn.Write([]byte("❌ Erro ao processar pagamento!\n"))
		return
	}

	// Open packet (thread-safe)
	openedPacket, err := game.OpenPacket(selectedPacket.ID)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("❌ %v\n", err)))
		// Refund tokens
		player.AddTokens(cost)
		return
	}

	// Add instrument to player
	player.AddInstrument(openedPacket.Instrument)

	conn.Write([]byte(fmt.Sprintf("\n🎉 VOCÊ ABRIU O PACOTE!\n")))
	conn.Write([]byte(fmt.Sprintf("📦 ID: %s\n", openedPacket.ID)))
	conn.Write([]byte(fmt.Sprintf("🎸 Instrumento: %s (%s)\n",
		openedPacket.Instrument.Name, openedPacket.Instrument.Rarity)))
	conn.Write([]byte(fmt.Sprintf("💰 %d tokens gastos\n", cost)))
}

func showInstruments(player *game.Player, conn net.Conn) {
	instruments := player.GetInstruments()

	if len(instruments) == 0 {
		conn.Write([]byte("\n📭 Você não tem instrumentos ainda!\n"))
		conn.Write([]byte("🎁 Abra pacotes para conseguir instrumentos.\n"))
		return
	}

	conn.Write([]byte(fmt.Sprintf("\n🎵 Seus Instrumentos (%d):\n", len(instruments))))
	for i, inst := range instruments {
		conn.Write([]byte(fmt.Sprintf("\n%d) %s (%s)\n", i+1, inst.Name, inst.Rarity)))
		conn.Write([]byte("   Ataques:\n"))
		for j, attack := range inst.Attacks {
			sequence := strings.Join(attack.Sequence, "-")
			conn.Write([]byte(fmt.Sprintf("   %d. %s: %s\n", j+1, attack.Name, sequence)))
		}
	}
}

func showTokens(player *game.Player, conn net.Conn) {
	tokens := player.GetTokens()
	conn.Write([]byte(fmt.Sprintf("\n💰 Seus Tokens: %d\n", tokens)))
}

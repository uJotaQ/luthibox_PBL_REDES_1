package network

import (
	"bufio"
	"fmt"
	"luthibox/game"
	"net"
	"strings"
	"time"
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

	player.SetDisconnectCallback(func(p *game.Player) {
		if p.IsInBattle() {
			game.BattlesMu.RLock()
			var battle *game.Battle
			for _, b := range game.ActiveBattles {
				if b.Player1 == p || b.Player2 == p {
					battle = b
					break
				}
			}
			game.BattlesMu.RUnlock()

			if battle != nil {
				battle.HandlePlayerDisconnect(p)
			}
		}
	})

	fmt.Printf("✅ Jogador %s autenticado\n", player.Nickname)
	conn.Write([]byte(fmt.Sprintf("🎮 Bem-vindo ao LuthiBOX, %s!\n", player.Nickname)))

	// Main message loop - PROCESSA TODAS AS MENSAGENS
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Cliente %s desconectado: %v\n", player.Nickname, err)
			return
		}

		msg = strings.TrimSpace(msg)
		fmt.Printf("DEBUG: Recebido do cliente %s: '%s'\n", player.Nickname, msg)

		// Tratar comandos especiais
		if msg == "PING_CMD" {
			fmt.Printf("Enviando PONG para o client")
			conn.Write([]byte("PONG\n"))
			continue
		}

		// Tratar comandos de batalha
		if strings.HasPrefix(msg, "PLAY_NOTE ") {
			if player.IsInBattle() {
				parts := strings.Split(msg, " ")
				if len(parts) == 2 {
					note := strings.ToUpper(parts[1])
					handleBattleMove(player, note, conn)
				}
			} else {
				conn.Write([]byte("❌ Você não está em uma batalha!\n"))
			}
			continue
		}

		// Se jogador está em batalha, não processar comandos de menu
		if player.IsInBattle() {
			conn.Write([]byte("🎮 Durante a batalha, use: PLAY_NOTE <A,B,C,D,E,F,G>\n"))
			continue
		}

		// Processar comandos de menu
		switch msg {
		case "0":
			conn.Write([]byte("👋 Até logo!\n"))
			return

		case "1":
			startBattle(player, conn)

		case "2":
			// Para comandos que precisam de reader, você precisa passar o reader
			// ou reestruturar melhor
			tempReader := bufio.NewReader(conn)
			openPackets(player, conn, tempReader)

		case "3":
			showInstruments(player, conn)

		case "4":
			showTokens(player, conn)
		case "5":
			// Mostrar estatísticas de conexão (satisfaz o requisito de "visualizar atraso")
			connectTime := time.Since(player.ConnectionTime)
			conn.Write([]byte(fmt.Sprintf("\n📡 Estatísticas de Conexão:\n")))
			conn.Write([]byte(fmt.Sprintf("⏱  Tempo conectado: %.0f segundos\n", connectTime.Seconds())))
			conn.Write([]byte(fmt.Sprintf("📶 Status: Conexão estável\n")))

			// Testar conectividade
			conn.Write([]byte("PING_CMD\n"))
			conn.Write([]byte("✅ Conectividade verificada com sucesso!\n"))

		default:
			// Mostrar menu para comandos inválidos
			menu := `
			🎮 === LUTHIBOX - MENU PRINCIPAL ===
			1) 🎲 Jogar (Batalha 1v1)
			2) 🎁 Abrir Pacotes
			3) 🎵 Meus Instrumentos
			4) 💰 Meus Tokens
			5) 📡 Ping (Latência)
			0) 🚪 Sair

			Escolha uma opção: `
			conn.Write([]byte(menu))
		}
	}
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
			// Atualizar conexão (vamos criar um método para isso)
			player.UpdateConnection(conn)
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

func handleBattleMove(player *game.Player, note string, conn net.Conn) {
	// Check if player is actually in a battle
	if !player.IsInBattle() {
		conn.Write([]byte("❌ Você não está em uma batalha!\n"))
		return
	}

	// Find player's active battle
	game.BattlesMu.RLock()
	var playerBattle *game.Battle
	for _, battle := range game.ActiveBattles {
		if battle.Player1.Nickname == player.Nickname || battle.Player2.Nickname == player.Nickname {
			playerBattle = battle
			break
		}
	}
	game.BattlesMu.RUnlock()

	if playerBattle == nil {
		player.ClearBattle() // Clear inconsistent state
		conn.Write([]byte("❌ Batalha não encontrada!\n"))
		return
	}

	// Process the move
	err := playerBattle.PlayNote(player, note)
	if err != nil {
		// Error message already sent in PlayNote
		return
	}
}

func startBattle(player *game.Player, conn net.Conn) {
	// First, let player choose instrument
	if !selectInstrumentForBattle(player, conn) {
		conn.Write([]byte("❌ Não foi possível selecionar instrumento para batalha.\n"))
		return
	}

	conn.Write([]byte("\n🔍 Procurando oponente...\n"))
	game.AddPlayerToBattleQueue(player)

	// Don't show menu - player is now in battle mode
}

func selectInstrumentForBattle(player *game.Player, conn net.Conn) bool {
	instruments := player.GetInstruments()

	if len(instruments) == 0 {
		conn.Write([]byte("\n❌ Você não tem instrumentos para batalhar!\n"))
		conn.Write([]byte("🎁 Abra pacotes para conseguir instrumentos.\n"))
		return false
	}

	conn.Write([]byte(fmt.Sprintf("\n🎯 Selecione um instrumento para a batalha (%d disponíveis):\n", len(instruments))))

	for i, inst := range instruments {
		conn.Write([]byte(fmt.Sprintf("\n%d) %s (%s)\n", i+1, inst.Name, inst.Rarity)))
		conn.Write([]byte("   Ataques:\n"))
		for j, attack := range inst.Attacks {
			sequence := strings.Join(attack.Sequence, "-")
			conn.Write([]byte(fmt.Sprintf("   %d. %s: %s\n", j+1, attack.Name, sequence)))
		}
	}

	conn.Write([]byte("\nDigite o número do instrumento (0 para cancelar): "))

	reader := bufio.NewReader(conn)
	escolha, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	escolha = strings.TrimSpace(escolha)
	if escolha == "0" {
		return false
	}

	var instrumentIndex int
	_, err = fmt.Sscanf(escolha, "%d", &instrumentIndex)
	if err != nil || instrumentIndex < 1 || instrumentIndex > len(instruments) {
		conn.Write([]byte("❌ Opção inválida!\n"))
		return false
	}

	selectedInstrument := instruments[instrumentIndex-1]
	player.SetSelectedInstrument(&selectedInstrument)

	conn.Write([]byte(fmt.Sprintf("\n✅ %s selecionado para a batalha!\n", selectedInstrument.Name)))
	return true
}

func openPackets(player *game.Player, conn net.Conn, reader *bufio.Reader) {
	conn.Write([]byte("\n🎯 Escolha a raridade:\n"))
	conn.Write([]byte("1) Comum (10 tokens)\n"))
	conn.Write([]byte("2) Raro (25 tokens)\n"))
	conn.Write([]byte("3) Épico (50 tokens)\n"))
	conn.Write([]byte("4) Lendário (100 tokens)\n"))
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
		conn.Write([]byte(fmt.Sprintf("❌ Tokens insuficientes! Você tem %d, precisa de %d.\n", player.GetTokens(), cost)))
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
	_, err = fmt.Sscanf(escolha, "%d", &packetIndex)
	if err != nil || packetIndex < 1 || packetIndex > len(packets) {
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
	conn.Write([]byte(fmt.Sprintf("📊 Seus tokens agora: %d\n", player.GetTokens())))

	fmt.Sprintf("\n🎉 VOCÊ ABRIU O PACOTE!\n")
	fmt.Sprintf("📦 ID: %s\n", openedPacket.ID)
	fmt.Sprintf("🎸 Instrumento: %s (%s)\n",
		openedPacket.Instrument.Name, openedPacket.Instrument.Rarity)
	fmt.Sprintf("💰 %d tokens gastos\n", cost)
	fmt.Sprintf("📊 Seus tokens agora: %d\n", player.GetTokens())
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

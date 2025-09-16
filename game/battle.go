package game

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Battle struct {
	ID           string
	Player1      *Player
	Player2      *Player
	PlayedNotes  []string
	Player1Score int
	Player2Score int
	CurrentTurn  int // 1 for Player1, 2 for Player2
	Active       bool
	mu           sync.RWMutex
}

var (
	ActiveBattles = make(map[string]*Battle)
	BattleQueue   = make(chan *Player, 100)
	BattlesMu     sync.RWMutex
)

func init() {
	go matchmakingSystem()
}

// Start matchmaking system
func matchmakingSystem() {
	for {
		player1 := <-BattleQueue
		player2 := <-BattleQueue

		// Create battle
		battle := &Battle{
			ID:          fmt.Sprintf("BATTLE_%d", len(ActiveBattles)),
			Player1:     player1,
			Player2:     player2,
			CurrentTurn: 1,
			Active:      true,
		}

		BattlesMu.Lock()
		ActiveBattles[battle.ID] = battle
		BattlesMu.Unlock()

		// Notify players
		notifyBattleStart(battle)
	}
}

// Add player to battle queue
func AddPlayerToBattleQueue(player *Player) {
	BattleQueue <- player
}

// Process note play in battle
// Atualizar a função PlayNote para mostrar mais informações:
func (b *Battle) PlayNote(player *Player, note string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check if it's player's turn
	currentPlayerNum := 1
	if player.Nickname == b.Player2.Nickname {
		currentPlayerNum = 2
	}

	if currentPlayerNum != b.CurrentTurn {
		player.Conn.Write([]byte("❌ Não é sua vez!\n"))
		return fmt.Errorf("não é sua vez")
	}

	// Validate note (A-G)
	validNotes := []string{"A", "B", "C", "D", "E", "F", "G"}
	isValid := false
	for _, validNote := range validNotes {
		if note == validNote {
			isValid = true
			break
		}
	}

	if !isValid {
		player.Conn.Write([]byte("❌ Nota inválida! Use A, B, C, D, E, F ou G\n"))
		return fmt.Errorf("nota inválida")
	}

	// Add note to sequence
	b.PlayedNotes = append(b.PlayedNotes, note)

	// Notify both players about the played note
	noteMessage := fmt.Sprintf("🎵 %s jogou nota: %s\n", player.Nickname, note)
	b.Player1.Conn.Write([]byte(noteMessage))
	b.Player2.Conn.Write([]byte(noteMessage))

	// Show current sequence to both players
	sequence := strings.Join(b.PlayedNotes, "-")
	sequenceMessage := fmt.Sprintf("📝 Sequência atual: %s\n", sequence)
	b.Player1.Conn.Write([]byte(sequenceMessage))
	b.Player2.Conn.Write([]byte(sequenceMessage))

	// 🔥 CRÍTICO: Verificar ataques de AMBOS os jogadores
	attackTriggered := false
	var attackPlayer *Player
	var attackName string

	// Verificar se algum jogador completou um ataque
	// Verificar ataques do Player1
	attackName1 := b.checkAttackCompletion(b.Player1)
	if attackName1 != "" {
		attackTriggered = true
		attackPlayer = b.Player1
		attackName = attackName1
	}

	// Verificar ataques do Player2 (se o Player1 não tiver acertado)
	if !attackTriggered {
		attackName2 := b.checkAttackCompletion(b.Player2)
		if attackName2 != "" {
			attackTriggered = true
			attackPlayer = b.Player2
			attackName = attackName2
		}
	}

	// Se algum ataque foi acertado
	if attackTriggered {
		// Incrementar score do jogador que tem o ataque
		if attackPlayer == b.Player1 {
			b.Player1Score++
			b.Player1.Conn.Write([]byte(fmt.Sprintf("🎉 ATAQUE '%s' REALIZADO! Pontos: %d\n", attackName, b.Player1Score)))
			b.Player2.Conn.Write([]byte(fmt.Sprintf("❌ Oponente realizou ataque '%s'! Pontos dele: %d\n", attackName, b.Player1Score)))
		} else {
			b.Player2Score++
			b.Player2.Conn.Write([]byte(fmt.Sprintf("🎉 ATAQUE '%s' REALIZADO! Pontos: %d\n", attackName, b.Player2Score)))
			b.Player1.Conn.Write([]byte(fmt.Sprintf("❌ Oponente realizou ataque '%s'! Pontos dele: %d\n", attackName, b.Player2Score)))
		}

		// Clear notes
		b.PlayedNotes = []string{}
		b.Player1.Conn.Write([]byte("📝 Sequência resetada\n"))
		b.Player2.Conn.Write([]byte("📝 Sequência resetada\n"))

		// Check for victory
		if b.Player1Score >= 2 || b.Player2Score >= 2 {
			b.endBattle()
			return nil
		}
	}

	// Switch turns apenas se NENHUM ataque foi acertado
	if !attackTriggered {
		if b.CurrentTurn == 1 {
			b.CurrentTurn = 2
			b.Player1.Conn.Write([]byte("⏳ Aguarde a vez do oponente...\n"))
			b.Player2.Conn.Write([]byte("🎮 Sua vez! Use PLAY_NOTE <nota>\n"))
		} else {
			b.CurrentTurn = 1
			b.Player2.Conn.Write([]byte("⏳ Aguarde a vez do oponente...\n"))
			b.Player1.Conn.Write([]byte("🎮 Sua vez! Use PLAY_NOTE <nota>\n"))
		}
	} else {
		// Se ataque foi acertado, o mesmo jogador continua
		// (ou você pode decidir quem continua - vamos manter o mesmo jogador por agora)
		currentPlayerName := b.Player1.Nickname
		if b.CurrentTurn == 2 {
			currentPlayerName = b.Player2.Nickname
		}
		playerMsg := fmt.Sprintf("🎮 %s continua! Use PLAY_NOTE <nota>\n", currentPlayerName)
		b.Player1.Conn.Write([]byte(playerMsg))
		b.Player2.Conn.Write([]byte(playerMsg))
	}

	return nil
}

func (b *Battle) checkAttackCompletion(forPlayer *Player) string {
	// Get player's selected instrument
	instrument := forPlayer.GetSelectedInstrument()

	// If no instrument selected, try to get first available
	if instrument == nil {
		instruments := forPlayer.GetInstruments()
		if len(instruments) > 0 {
			instrument = &instruments[0]
			forPlayer.SetSelectedInstrument(instrument)
		} else {
			// Use default instrument
			defaultInst := GetRandomInstrumentByRarity("Comum")
			if defaultInst != nil {
				instrument = defaultInst
			} else {
				return "" // No instruments available
			}
		}
	}

	// DebugLog("Checking attacks for %s with instrument %s", forPlayer.Nickname, instrument.Name)
	// DebugLog("Current sequence: %v (length: %d)", b.PlayedNotes, len(b.PlayedNotes))

	// Check each attack of the instrument
	for _, attack := range instrument.Attacks {
		attackSequence := attack.Sequence
		sequenceLen := len(attackSequence)
		notesLen := len(b.PlayedNotes)

		// DebugLog("Checking attack '%s' - sequence: %v (length: %d)",
		//     attack.Name, attackSequence, sequenceLen)

		// Need at least the same number of notes as the attack sequence
		if notesLen >= sequenceLen && sequenceLen > 0 {
			// 🔥 NOVA LÓGICA: Procurar a sequência em qualquer posição
			// Verificar se a sequência do ataque aparece em ordem em qualquer parte

			// Percorrer todas as posições possíveis onde a sequência pode começar
			for startPos := 0; startPos <= notesLen-sequenceLen; startPos++ {
				// Verificar se a sequência completa aparece a partir desta posição
				match := true
				for i := 0; i < sequenceLen; i++ {
					if b.PlayedNotes[startPos+i] != attackSequence[i] {
						match = false
						break
					}
				}

				if match {
					// DebugLog("Attack '%s' MATCHED at position %d!", attack.Name, startPos)
					return attack.Name
				}
			}
		}
	}

	// DebugLog("No attack matched for %s", forPlayer.Nickname)
	return ""
}

// End battle and declare winner
func (b *Battle) endBattle() {
	var winner, loser *Player
	rewardTokens := 10

	if b.Player1Score >= 2 {
		winner = b.Player1
		loser = b.Player2
	} else {
		winner = b.Player2
		loser = b.Player1
	}

	// Give reward to winner
	winner.AddTokens(rewardTokens)

	// Notify results
	winner.Conn.Write([]byte(fmt.Sprintf("\n🏆 VITÓRIA! Você ganhou %d tokens!\n", rewardTokens)))
	loser.Conn.Write([]byte("\n💀 DERROTA! Mais sorte na próxima vez.\n"))

	// Show final score
	winner.Conn.Write([]byte(fmt.Sprintf("Placar final: Você %d - %d Oponente\n",
		func() int {
			if winner == b.Player1 {
				return b.Player1Score
			}
			return b.Player2Score
		}(),
		func() int {
			if loser == b.Player1 {
				return b.Player1Score
			}
			return b.Player2Score
		}())))

	loser.Conn.Write([]byte(fmt.Sprintf("Placar final: Você %d - %d Oponente\n",
		func() int {
			if loser == b.Player1 {
				return b.Player1Score
			}
			return b.Player2Score
		}(),
		func() int {
			if winner == b.Player1 {
				return b.Player1Score
			}
			return b.Player2Score
		}())))

	// Remove battle
	BattlesMu.Lock()
	delete(ActiveBattles, b.ID)
	BattlesMu.Unlock()

	b.Player1.ClearBattle()
	b.Player2.ClearBattle()

	b.Active = false
}

// Notification functions
func notifyBattleStart(battle *Battle) {
	// Set players in battle
	battle.Player1.SetCurrentBattle(battle.ID)
	battle.Player2.SetCurrentBattle(battle.ID)

	battle.Player1.Conn.Write([]byte(fmt.Sprintf("\n⚔️ BATALHA INICIADA!\n")))
	battle.Player1.Conn.Write([]byte(fmt.Sprintf("VS %s\n", battle.Player2.Nickname)))
	battle.Player1.Conn.Write([]byte("🎮 Sua vez! Use PLAY_NOTE <A,B,C,D,E,F,G>\n"))

	time.Sleep(10 * time.Millisecond)

	battle.Player2.Conn.Write([]byte(fmt.Sprintf("\n⚔️ BATALHA INICIADA!\n")))
	battle.Player2.Conn.Write([]byte(fmt.Sprintf("VS %s\n", battle.Player1.Nickname)))
	battle.Player2.Conn.Write([]byte("⏳ Aguarde a vez do oponente...\n"))
}

// Adicionar no final de game/battle.go
func (b *Battle) DebugInfo() {
	fmt.Printf("=== DEBUG BATTLE INFO ===\n")
	fmt.Printf("Battle ID: %s\n", b.ID)
	fmt.Printf("Player1: %s (Score: %d)\n", b.Player1.Nickname, b.Player1Score)
	fmt.Printf("Player2: %s (Score: %d)\n", b.Player2.Nickname, b.Player2Score)
	fmt.Printf("Current sequence: %v\n", b.PlayedNotes)
	fmt.Printf("Turn: %d\n", b.CurrentTurn)

	// Show player instruments
	if b.Player1.GetSelectedInstrument() != nil {
		inst := b.Player1.GetSelectedInstrument()
		fmt.Printf("Player1 instrument: %s\n", inst.Name)
		for i, attack := range inst.Attacks {
			fmt.Printf("  Attack %d: %s - %v\n", i+1, attack.Name, attack.Sequence)
		}
	}

	if b.Player2.GetSelectedInstrument() != nil {
		inst := b.Player2.GetSelectedInstrument()
		fmt.Printf("Player2 instrument: %s\n", inst.Name)
		for i, attack := range inst.Attacks {
			fmt.Printf("  Attack %d: %s - %v\n", i+1, attack.Name, attack.Sequence)
		}
	}
	fmt.Printf("========================\n")
}

func (b *Battle) EndBattle() {
    // Clear battle state for players
    b.Player1.ClearBattle()
    b.Player2.ClearBattle()
    
    // Remove battle
    BattlesMu.Lock()
    delete(ActiveBattles, b.ID)
    BattlesMu.Unlock()
    
    b.Active = false
}
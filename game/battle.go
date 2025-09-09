package game

import (
    "fmt"
    "sync"
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
    activeBattles = make(map[string]*Battle)
    battleQueue   = make(chan *Player, 100)
    battlesMu     sync.RWMutex
)

func init() {
    go matchmakingSystem()
}

// Start matchmaking system
func matchmakingSystem() {
    for {
        player1 := <-battleQueue
        player2 := <-battleQueue

        // Create battle
        battle := &Battle{
            ID:          fmt.Sprintf("BATTLE_%d", len(activeBattles)),
            Player1:     player1,
            Player2:     player2,
            CurrentTurn: 1,
            Active:      true,
        }

        battlesMu.Lock()
        activeBattles[battle.ID] = battle
        battlesMu.Unlock()

        // Notify players
        notifyBattleStart(battle)
    }
}

// Add player to battle queue
func AddPlayerToBattleQueue(player *Player) {
    battleQueue <- player
}

// Process note play in battle
func (b *Battle) PlayNote(player *Player, note string) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    // Check if it's player's turn
    currentPlayerNum := 1
    if player.Nickname == b.Player2.Nickname {
        currentPlayerNum = 2
    }

    if currentPlayerNum != b.CurrentTurn {
        return fmt.Errorf("não é sua vez")
    }

    // Add note to sequence
    b.PlayedNotes = append(b.PlayedNotes, note)

    // Check for completed attacks (simplified)
    attackCompleted := b.checkAttackCompletion(player)

    if attackCompleted {
        // Increment score
        if currentPlayerNum == 1 {
            b.Player1Score++
        } else {
            b.Player2Score++
        }

        // Clear notes
        b.PlayedNotes = []string{}

        // Notify about attack
        notifyAttack(b, player)

        // Check for victory
        if b.Player1Score >= 2 || b.Player2Score >= 2 {
            b.endBattle()
            return nil
        }
    }

    // Switch turns
    if b.CurrentTurn == 1 {
        b.CurrentTurn = 2
    } else {
        b.CurrentTurn = 1
    }

    return nil
}

// Simplified attack checking
func (b *Battle) checkAttackCompletion(player *Player) bool {
    // For demo purposes, any 3 notes complete an attack
    return len(b.PlayedNotes) >= 3
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
    notifyBattleEnd(b, winner, loser, rewardTokens)

    // Remove battle
    battlesMu.Lock()
    delete(activeBattles, b.ID)
    battlesMu.Unlock()

    b.Active = false
}

// Notification functions
func notifyBattleStart(battle *Battle) {
    battle.Player1.Conn.Write([]byte(fmt.Sprintf("=== BATALHA INICIADA ===\nVocê vs %s\n", battle.Player2.Nickname)))
    battle.Player2.Conn.Write([]byte(fmt.Sprintf("=== BATALHA INICIADA ===\nVocê vs %s\n", battle.Player1.Nickname)))
    battle.Player1.Conn.Write([]byte("Use PLAY_NOTE <nota> para jogar (A,B,C,D,E,F,G)\n"))
    battle.Player2.Conn.Write([]byte("Use PLAY_NOTE <nota> para jogar (A,B,C,D,E,F,G)\n"))
}

func notifyAttack(battle *Battle, player *Player) {
    message := fmt.Sprintf("🎉 ATAQUE REALIZADO! Pontos: %d\n", 
        func() int {
            if player.Nickname == battle.Player1.Nickname {
                return battle.Player1Score
            }
            return battle.Player2Score
        }())
    
    player.Conn.Write([]byte(message))
    
    // Notify opponent
    opponent := battle.Player1
    if player.Nickname == battle.Player1.Nickname {
        opponent = battle.Player2
    }
    opponent.Conn.Write([]byte("❌ Oponente realizou um ataque!\n"))
}

func notifyBattleEnd(battle *Battle, winner, loser *Player, tokens int) {
    winner.Conn.Write([]byte(fmt.Sprintf("🏆 VITÓRIA! Você ganhou %d tokens!\n", tokens)))
    loser.Conn.Write([]byte("💀 DERROTA! Mais sorte na próxima vez.\n"))
}

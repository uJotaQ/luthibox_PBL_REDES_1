package game

import (
    "fmt"
    "sync"
)

// GameManager simples para gerenciar partidas
type GameManager struct {
    Battles map[string]*Battle
    Queue   chan *Player
    mu      sync.RWMutex
}

var defaultManager = &GameManager{
    Battles: make(map[string]*Battle),
    Queue:   make(chan *Player, 100),
}

// GetGameManager retorna a instância do gerenciador
func GetGameManager() *GameManager {
    return defaultManager
}

// Start inicia o sistema de partidas
func (gm *GameManager) Start() {
    go gm.matchmakingLoop()
}

// matchmakingLoop pareia jogadores automaticamente
func (gm *GameManager) matchmakingLoop() {
    for {
        player1 := <-gm.Queue
        player2 := <-gm.Queue

        // Criar partida
        battle := &Battle{
            ID:          fmt.Sprintf("BATTLE_%d", len(gm.Battles)),
            Player1:     player1,
            Player2:     player2,
            CurrentTurn: 1,
            Active:      true,
        }

        // Adicionar à lista de partidas
        gm.mu.Lock()
        gm.Battles[battle.ID] = battle
        gm.mu.Unlock()

        // Notificar jogadores que a partida começou
        gm.notifyBattleStart(battle)
    }
}

// AddPlayerToQueue adiciona jogador à fila de espera
func (gm *GameManager) AddPlayerToQueue(player *Player) {
    gm.Queue <- player
}

// GetBattle pega uma partida pelo ID
func (gm *GameManager) GetBattle(battleID string) (*Battle, bool) {
    gm.mu.RLock()
    defer gm.mu.RUnlock()

    battle, exists := gm.Battles[battleID]
    return battle, exists
}

// RemoveBattle remove uma partida finalizada
func (gm *GameManager) RemoveBattle(battleID string) {
    gm.mu.Lock()
    defer gm.mu.Unlock()

    delete(gm.Battles, battleID)
}

// notifyBattleStart notifica jogadores do início da partida
func (gm *GameManager) notifyBattleStart(battle *Battle) {
    battle.Player1.Conn.Write([]byte(fmt.Sprintf("⚔️ BATALHA INICIADA contra %s!\n", battle.Player2.Nickname)))
    battle.Player1.Conn.Write([]byte("Use: PLAY_NOTE <A,B,C,D,E,F,G>\n"))

    battle.Player2.Conn.Write([]byte(fmt.Sprintf("⚔️ BATALHA INICIADA contra %s!\n", battle.Player1.Nickname)))
    battle.Player2.Conn.Write([]byte("Use: PLAY_NOTE <A,B,C,D,E,F,G>\n"))
}

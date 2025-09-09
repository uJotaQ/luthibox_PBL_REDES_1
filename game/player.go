package game

import (
	"fmt"
	"net"
	"sync"
)

type Player struct {
	ID          string
	Nickname    string
	Password    string
	Conn        net.Conn
	Instruments []Instrument
	Tokens      int
	mu          sync.RWMutex
}

var (
	players   = make(map[string]*Player) // nickname -> player
	playersMu sync.RWMutex
)

// Register new player
func RegisterPlayer(nickname, password string, conn net.Conn) error {
	playersMu.Lock()
	defer playersMu.Unlock()

	if _, exists := players[nickname]; exists {
		return fmt.Errorf("nickname já existe")
	}

	player := &Player{
		ID:       fmt.Sprintf("PLAYER_%d", len(players)),
		Nickname: nickname,
		Password: password,
		Conn:     conn,
		Tokens:   100, // Starting tokens
	}

	players[nickname] = player
	return nil
}

// Authenticate player
func AuthenticatePlayer(nickname, password string) (*Player, error) {
	playersMu.RLock()
	defer playersMu.RUnlock()

	player, exists := players[nickname]
	if !exists {
		return nil, fmt.Errorf("jogador não encontrado")
	}

	if player.Password != password {
		return nil, fmt.Errorf("senha incorreta")
	}

	return player, nil
}

// Get player by nickname
func GetPlayer(nickname string) (*Player, bool) {
	playersMu.RLock()
	defer playersMu.RUnlock()

	player, exists := players[nickname]
	return player, exists
}

// Add tokens to player
func (p *Player) AddTokens(amount int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Tokens += amount
}

// Remove tokens from player
func (p *Player) RemoveTokens(amount int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.Tokens >= amount {
		p.Tokens -= amount
		return true
	}
	return false
}

// Add instrument to player
func (p *Player) AddInstrument(instrument Instrument) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Instruments = append(p.Instruments, instrument)
}

// Get player's instruments
func (p *Player) GetInstruments() []Instrument {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Instruments
}

// Get player's tokens
func (p *Player) GetTokens() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Tokens
}

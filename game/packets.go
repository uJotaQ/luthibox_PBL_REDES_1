package game

import (
    "fmt"
    "math/rand"
    "sync"
)

type Packet struct {
    ID          string
    Rarity      string
    Instrument  Instrument
    Opened      bool
    mu          sync.RWMutex
}

var (
    packetStock = make(map[string]*Packet) // Global packet stock
    stockMu     sync.RWMutex
)

func init() {
    initializePacketStock()
}

// Initialize packet stock
func initializePacketStock() {
    rarities := []string{"Comum", "Raro", "Épico", "Lendário"}
    quantities := map[string]int{
        "Comum":    20,
        "Raro":     15,
        "Épico":    10,
        "Lendário": 5,
    }

    for _, rarity := range rarities {
        for i := 0; i < quantities[rarity]; i++ {
            packet := generatePacket(rarity)
            AddPacketToStock(packet)
        }
    }
}

// Generate a packet with random instrument
func generatePacket(rarity string) *Packet {
    instrument := GetRandomInstrumentByRarity(rarity)
    if instrument == nil {
        // Fallback to any instrument
        allInstruments := GetAllInstruments()
        if len(allInstruments) > 0 {
            instrument = &allInstruments[rand.Intn(len(allInstruments))]
        }
    }

    return &Packet{
        ID:         generatePacketID(),
        Rarity:     rarity,
        Instrument: *instrument,
        Opened:     false,
    }
}

// Add packet to global stock (thread-safe)
func AddPacketToStock(packet *Packet) {
    stockMu.Lock()
    defer stockMu.Unlock()
    packetStock[packet.ID] = packet
}

// Get available packets by rarity (thread-safe)
func GetAvailablePacketsByRarity(rarity string) []*Packet {
    stockMu.RLock()
    defer stockMu.RUnlock()

    var available []*Packet
    for _, packet := range packetStock {
        packet.mu.RLock()
        if !packet.Opened && packet.Rarity == rarity {
            available = append(available, packet)
        }
        packet.mu.RUnlock()
    }

    // Limit to 10 packets
    if len(available) > 10 {
        available = available[:10]
    }

    return available
}

// Open packet (thread-safe and concurrent-safe)
func OpenPacket(packetID string) (*Packet, error) {
    stockMu.Lock()
    packet, exists := packetStock[packetID]
    if !exists {
        stockMu.Unlock()
        return nil, fmt.Errorf("pacote não encontrado")
    }

    packet.mu.Lock()
    if packet.Opened {
        packet.mu.Unlock()
        stockMu.Unlock()
        return nil, fmt.Errorf("pacote já foi aberto")
    }

    packet.Opened = true
    openedPacket := *packet // Make a copy
    packet.mu.Unlock()
    
    delete(packetStock, packetID) // Remove from stock
    stockMu.Unlock()

    // Replenish stock with new packet of same rarity
    newPacket := generatePacket(openedPacket.Rarity)
    AddPacketToStock(newPacket)

    return &openedPacket, nil
}

// Generate unique packet ID
func generatePacketID() string {
    letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    numbers := "0123456789"
    return fmt.Sprintf("%c%c%c", 
        letters[rand.Intn(len(letters))],
        numbers[rand.Intn(len(numbers))],
        numbers[rand.Intn(len(numbers))])
}

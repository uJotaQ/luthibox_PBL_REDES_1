package game

import (
	"math/rand"
	"time"
)

type Attack struct {
	Name     string
	Sequence []string // Ex: {"A", "B", "G"}
}

type Instrument struct {
	ID      string
	Name    string
	Rarity  string // Comum, Raro, Épico, Lendário
	Attacks [3]Attack
	Price   int
}

var instrumentDatabase []Instrument

func init() {
	rand.Seed(time.Now().UnixNano())
	loadInstruments()
}

func loadInstruments() {
	instrumentDatabase = []Instrument{
		{
			ID:     "VIOLIN_001",
			Name:   "Violino",
			Rarity: "Raro",
			Price:  50,
			Attacks: [3]Attack{
				{"Vibrato", []string{"A", "B", "G"}},
				{"Pizzicato", []string{"C", "E", "G"}},
				{"Spiccato", []string{"F", "A", "C"}},
			},
		},
		{
			ID:     "GUITAR_001",
			Name:   "Guitarra",
			Rarity: "Épico",
			Price:  100,
			Attacks: [3]Attack{
				{"Power Chord", []string{"E", "A", "D"}},
				{"Solo Rápido", []string{"G", "B", "E"}},
				{"Palm Mute", []string{"E", "E", "A"}},
			},
		},
		{
			ID:     "FLUTE_001",
			Name:   "Flauta",
			Rarity: "Comum",
			Price:  25,
			Attacks: [3]Attack{
				{"Trinado", []string{"C", "D", "E"}},
				{"Sopro Forte", []string{"F", "G", "C"}},
				{"Melodia Doce", []string{"A", "B", "C"}},
			},
		},
		{
			ID:     "PIANO_001",
			Name:   "Piano",
			Rarity: "Épico",
			Price:  100,
			Attacks: [3]Attack{
				{"Arpejo", []string{"C", "E", "G"}},
				{"Glissando", []string{"D", "F", "A"}},
				{"Acorde Pesado", []string{"E", "G", "B"}},
			},
		},
		{
			ID:     "DRUM_001",
			Name:   "Bateria",
			Rarity: "Épico",
			Price:  100,
			Attacks: [3]Attack{
				{"Rufar", []string{"C", "C", "C"}},
				{"Virada", []string{"E", "G", "A"}},
				{"Prato Explosivo", []string{"F", "B", "E"}},
			},
		},
	}
}

// Get random instrument by rarity
func GetRandomInstrumentByRarity(rarity string) *Instrument {
	var filtered []Instrument
	for _, inst := range instrumentDatabase {
		if inst.Rarity == rarity {
			filtered = append(filtered, inst)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	return &filtered[rand.Intn(len(filtered))]
}

// Get all instruments
func GetAllInstruments() []Instrument {
	return instrumentDatabase
}


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
	// Violino
	violino := Instrument{
		ID:     "VIOLIN_001",
		Name:   "Violino",
		Rarity: "Raro",
		Price:  50,
		Attacks: [3]Attack{
			{"Vibrato", []string{"A", "B", "G"}},
			{"Pizzicato", []string{"C", "E", "G"}},
			{"Spiccato", []string{"F", "A", "C"}},
		},
	}

	// Guitarra
	guitarra := Instrument{
		ID:     "GUITAR_001",
		Name:   "Guitarra",
		Rarity: "Épico",
		Price:  100,
		Attacks: [3]Attack{
			{"Power Chord", []string{"E", "A", "D"}},
			{"Solo Rápido", []string{"G", "B", "E"}},
			{"Palm Mute", []string{"E", "E", "A"}},
		},
	}

	// Flauta
	flauta := Instrument{
		ID:     "FLUTE_001",
		Name:   "Flauta",
		Rarity: "Comum",
		Price:  25,
		Attacks: [3]Attack{
			{"Trinado", []string{"C", "D", "E"}},
			{"Sopro Forte", []string{"F", "G", "C"}},
			{"Melodia Doce", []string{"A", "B", "C"}},
		},
	}

	// Piano
	piano := Instrument{
		ID:     "PIANO_001",
		Name:   "Piano",
		Rarity: "Épico",
		Price:  100,
		Attacks: [3]Attack{
			{"Arpejo", []string{"C", "E", "G"}},
			{"Glissando", []string{"D", "F", "A"}},
			{"Acorde Pesado", []string{"E", "G", "B"}},
		},
	}

	// Trompete
	trompete := Instrument{
		ID:     "TRUMPET_001",
		Name:   "Trompete",
		Rarity: "Raro",
		Price:  50,
		Attacks: [3]Attack{
			{"Susto Sonoro", []string{"C", "G", "C"}},
			{"Escalada", []string{"D", "F", "A"}},
			{"Sopro Metálico", []string{"E", "A", "B"}},
		},
	}

	// Bateria
	bateria := Instrument{
		ID:     "DRUM_001",
		Name:   "Bateria",
		Rarity: "Épico",
		Price:  100,
		Attacks: [3]Attack{
			{"Rufar", []string{"C", "C", "C"}},
			{"Virada", []string{"E", "G", "A"}},
			{"Prato Explosivo", []string{"F", "B", "E"}},
		},
	}

	// Harpa
	harpa := Instrument{
		ID:     "HARP_001",
		Name:   "Harpa",
		Rarity: "Lendário",
		Price:  200,
		Attacks: [3]Attack{
			{"Arpejo Celestial", []string{"C", "F", "A"}},
			{"Glissando Suave", []string{"D", "G", "B"}},
			{"Eco Etéreo", []string{"E", "A", "C"}},
		},
	}

	// Saxofone
	saxofone := Instrument{
		ID:     "SAX_001",
		Name:   "Saxofone",
		Rarity: "Raro",
		Price:  50,
		Attacks: [3]Attack{
			{"Improviso Jazz", []string{"E", "G", "B"}},
			{"Sopro Grave", []string{"C", "E", "A"}},
			{"Nota Sustentada", []string{"D", "F", "G"}},
		},
	}

	// Oboé
	oboe := Instrument{
		ID:     "OBOE_001",
		Name:   "Oboé",
		Rarity: "Comum",
		Price:  25,
		Attacks: [3]Attack{
			{"Melodia Suave", []string{"C", "D", "E"}},
			{"Sopro Agudo", []string{"F", "A", "B"}},
			{"Eco Vibrante", []string{"G", "B", "C"}},
		},
	}

	// Contrabaixo
	contrabaixo := Instrument{
		ID:     "BASS_001",
		Name:   "Contrabaixo",
		Rarity: "Raro",
		Price:  50,
		Attacks: [3]Attack{
			{"Slap", []string{"E", "A", "D"}},
			{"Walking Bass", []string{"G", "B", "E"}},
			{"Grave Profundo", []string{"C", "E", "A"}},
		},
	}

	// Banjo
	banjo := Instrument{
		ID:     "BANJO_001",
		Name:   "Banjo",
		Rarity: "Comum",
		Price:  25,
		Attacks: [3]Attack{
			{"Dedilhado", []string{"C", "G", "B"}},
			{"Rasgueado", []string{"D", "A", "E"}},
			{"Bluegrass", []string{"F", "G", "C"}},
		},
	}

	// Tambor
	tambor := Instrument{
		ID:     "TAMBO_001",
		Name:   "Tambor",
		Rarity: "Comum",
		Price:  25,
		Attacks: [3]Attack{
			{"Batida Rítmica", []string{"C", "C", "G"}},
			{"Toque Seco", []string{"D", "F", "A"}},
			{"Rufada Tribal", []string{"E", "A", "C"}},
		},
	}

	// Acordeon
	acordeon := Instrument{
		ID:     "ACCORD_001",
		Name:   "Acordeon",
		Rarity: "Raro",
		Price:  50,
		Attacks: [3]Attack{
			{"Expansão de Fole", []string{"C", "E", "A"}},
			{"Contra-Melodia", []string{"D", "G", "B"}},
			{"Dança Gaúcha", []string{"F", "A", "C"}},
		},
	}

	// Ukulele
	ukulele := Instrument{
		ID:     "UKULE_001",
		Name:   "Ukulele",
		Rarity: "Comum",
		Price:  25,
		Attacks: [3]Attack{
			{"Acorde Leve", []string{"C", "E", "G"}},
			{"Ritmo Havaiano", []string{"D", "G", "A"}},
			{"Ponte Alegre", []string{"E", "A", "C"}},
		},
	}

	// Adiciona todos ao banco
	instrumentDatabase = append(
		instrumentDatabase,
		violino, guitarra, flauta,
		piano, trompete, bateria, harpa,
		saxofone, oboe, contrabaixo, banjo,
		tambor, acordeon, ukulele,
	)
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

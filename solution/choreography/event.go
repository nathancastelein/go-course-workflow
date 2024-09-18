package choreography

import "github.com/nathancastelein/go-course-workflows/pokemon"

type EventType string
type MessageBrokerTopic chan (Event)

var (
	PokemonEncountered EventType = "PokemonEncountered"
	PokemonWeakened    EventType = "PokemonWeakened"
	PokemonParalyzed   EventType = "PokemonParalyzed"
	PokemonCaptured    EventType = "PokemonCaptured"
)

type Event struct {
	Type    EventType
	Pokemon *pokemon.Pokemon
	Trainer *pokemon.Trainer
}

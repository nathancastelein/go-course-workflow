package pokemon

import (
	"log/slog"
	"math/rand/v2"
	"strings"
)

func Sacha() *Trainer {
	return &Trainer{
		ID:       rand.IntN(10000),
		Name:     "Sacha",
		Pokemons: []*Pokemon{Pikachu()},
	}
}

type Trainer struct {
	ID       int
	Name     string
	Pokemons []*Pokemon
}

func (t Trainer) LogValue() slog.Value {
	pokemons := make([]string, len(t.Pokemons))
	for i, pokemon := range t.Pokemons {
		pokemons[i] = pokemon.Name
	}
	return slog.GroupValue(
		slog.Int("id", t.ID),
		slog.String("name", t.Name),
		slog.String("pokemons", strings.Join(pokemons, ",")),
	)
}

package mocks

import pokemon "github.com/nathancastelein/go-course-workflows/pokemon"

func Blue() *pokemon.Trainer {
	return &pokemon.Trainer{
		ID:       1,
		Name:     "Blue",
		Pokemons: []*pokemon.Pokemon{Squirtle()},
	}
}

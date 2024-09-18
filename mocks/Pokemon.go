package mocks

import (
	pokemon "github.com/nathancastelein/go-course-workflows/pokemon"
)

func Squirtle() *pokemon.Pokemon {
	return &pokemon.Pokemon{
		ID:            42,
		Name:          "Squirtle",
		Level:         5,
		CurrentHealth: 20,
		MaxHealth:     20,
		Status:        pokemon.StatusHealthy,
	}
}

func Rattata() *pokemon.Pokemon {
	return &pokemon.Pokemon{
		ID:            19,
		Name:          "Rattata",
		Level:         3,
		CurrentHealth: 15,
		MaxHealth:     15,
		Status:        pokemon.StatusHealthy,
	}
}

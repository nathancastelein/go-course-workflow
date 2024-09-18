package capture

import (
	"context"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

func (w *Worker) ParalyzeActivity(ctx context.Context, pokemon *pokemon.Pokemon) (*pokemon.Pokemon, error) {
	return pokemon, w.status.Paralyze(pokemon)
}

func (w *Worker) AttackActivity(ctx context.Context, pokemon *pokemon.Pokemon) (*pokemon.Pokemon, error) {
	return pokemon, w.combat.Attack(pokemon)
}

type ThrowPokeballOutput struct {
	Pokemon *pokemon.Pokemon
	Trainer *pokemon.Trainer
}

func (w *Worker) ThrowPokeballActivity(ctx context.Context, trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) (ThrowPokeballOutput, error) {
	return ThrowPokeballOutput{Pokemon: pokemon, Trainer: trainer}, w.pokeball.Throw(trainer, pokemon)
}

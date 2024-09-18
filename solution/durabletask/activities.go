package durabletask

import (
	"github.com/microsoft/durabletask-go/task"
	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type ActivityInput struct {
	Trainer *pokemon.Trainer
	Pokemon *pokemon.Pokemon
}

func NewActivityInput(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) ActivityInput {
	return ActivityInput{
		Trainer: trainer,
		Pokemon: pokemon,
	}
}

func (w *Worker) Paralyze(ctx task.ActivityContext) (any, error) {
	var input ActivityInput
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}

	err := w.status.Paralyze(input.Pokemon)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (w *Worker) Attack(ctx task.ActivityContext) (any, error) {
	var input ActivityInput
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}

	err := w.combat.Attack(input.Pokemon)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (w *Worker) ThrowPokeball(ctx task.ActivityContext) (any, error) {
	var input ActivityInput
	if err := ctx.GetInput(&input); err != nil {
		return nil, err
	}

	err := w.pokeball.Throw(input.Trainer, input.Pokemon)
	if err != nil {
		return nil, err
	}

	return input, nil
}

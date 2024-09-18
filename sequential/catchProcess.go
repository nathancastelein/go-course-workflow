package sequential

import (
	"context"
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type CatchProcess struct {
	status   pokemon.StatusService
	combat   pokemon.CombatService
	pokeball pokemon.PokeballService
}

func NewCatchProcess(
	status pokemon.StatusService,
	combat pokemon.CombatService,
	pokeball pokemon.PokeballService,
) *CatchProcess {
	return &CatchProcess{
		status:   status,
		combat:   combat,
		pokeball: pokeball,
	}
}

func (g *CatchProcess) Do(ctx context.Context, trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
	log := func() *slog.Logger {
		return slog.With(
			slog.Any("trainer", trainer),
			slog.Any("pokemon", pokemon),
		)
	}

	log().Info("trying to catch pokemon")

	log().Info("paralyze pokemon")
	err := g.status.Paralyze(pokemon)
	if err != nil {
		return err
	}

	log().Info("attack pokemon")
	err = g.combat.Attack(pokemon)
	if err != nil {
		return err
	}

	log().Info("throwing a pokeball")
	err = g.pokeball.Throw(trainer, pokemon)
	if err != nil {
		return err
	}

	log().Info("pokemon caught")
	return nil
}

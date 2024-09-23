package sequential

import (
	"context"
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type CaptureProcess struct {
	status   pokemon.StatusService
	combat   pokemon.CombatService
	pokeball pokemon.PokeballService
}

func NewCaptureProcess(
	status pokemon.StatusService,
	combat pokemon.CombatService,
	pokeball pokemon.PokeballService,
) *CaptureProcess {
	return &CaptureProcess{
		status:   status,
		combat:   combat,
		pokeball: pokeball,
	}
}

func (g *CaptureProcess) Do(ctx context.Context, trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error {
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

	log().Info("pokemon captured")
	return nil
}

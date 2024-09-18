package choreography

import (
	"context"
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type CombatWorker struct {
	combat        pokemon.CombatService
	combatTopic   MessageBrokerTopic
	pokeballTopic MessageBrokerTopic
}

func NewCombatWorker(
	combat pokemon.CombatService,
	combatTopic MessageBrokerTopic,
	pokeballTopic MessageBrokerTopic,
) CombatWorker {
	return CombatWorker{
		combat:        combat,
		combatTopic:   combatTopic,
		pokeballTopic: pokeballTopic,
	}
}

func (c *CombatWorker) Run(ctx context.Context) {
	slog.Info("starting combat worker")
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c.combatTopic:
			switch event.Type {
			case PokemonParalyzed:
				slog.Info("attack pokemon", slog.Any("pokemon", event.Pokemon))
				err := c.combat.Attack(event.Pokemon)
				if err != nil {
					slog.Error("fail to attack pokemon", slog.Any("error", err))
					continue
				}

				c.pokeballTopic <- Event{
					Type:    PokemonWeakened,
					Pokemon: event.Pokemon,
					Trainer: event.Trainer,
				}
			}
		}
	}
}

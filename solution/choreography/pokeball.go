package choreography

import (
	"context"
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type PokeballWorker struct {
	pokeball      pokemon.PokeballService
	pokeballTopic MessageBrokerTopic
	captureTopic  MessageBrokerTopic
}

func NewPokeballWorker(
	pokeball pokemon.PokeballService,
	pokeballTopic MessageBrokerTopic,
	captureTopic MessageBrokerTopic,
) PokeballWorker {
	return PokeballWorker{
		pokeball:      pokeball,
		pokeballTopic: pokeballTopic,
		captureTopic:  captureTopic,
	}
}

func (p *PokeballWorker) Run(ctx context.Context) {
	slog.Info("starting pokeball worker")
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-p.pokeballTopic:
			switch event.Type {
			case PokemonWeakened:
				slog.Info("throw pokeball", slog.Any("pokemon", event.Pokemon))
				err := p.pokeball.Throw(event.Trainer, event.Pokemon)
				if err != nil {
					slog.Error("fail to throw pokeball", slog.Any("error", err))
					continue
				}

				slog.Info("pokemon Captured", slog.Any("pokemon", event.Pokemon), slog.Any("trainer", event.Trainer))

				p.captureTopic <- Event{
					Type:    PokemonCaptured,
					Pokemon: event.Pokemon,
					Trainer: event.Trainer,
				}
			}
		}
	}
}

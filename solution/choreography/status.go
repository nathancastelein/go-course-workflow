package choreography

import (
	"context"
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type StatusWorker struct {
	status      pokemon.StatusService
	statusTopic MessageBrokerTopic
	combatTopic MessageBrokerTopic
}

func NewStatusWorker(
	status pokemon.StatusService,
	statusTopic MessageBrokerTopic,
	combatTopic MessageBrokerTopic,
) StatusWorker {
	return StatusWorker{
		status:      status,
		statusTopic: statusTopic,
		combatTopic: combatTopic,
	}
}

func (s *StatusWorker) Run(ctx context.Context) {
	slog.Info("starting status worker")
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.statusTopic:
			switch event.Type {
			case PokemonEncountered:
				slog.Info("paralyze pokemon", slog.Any("pokemon", event.Pokemon))

				err := s.status.Paralyze(event.Pokemon)
				if err != nil {
					slog.Error("fail to paralyze pokemon", slog.Any("error", err))
					continue
				}

				s.combatTopic <- Event{
					Type:    PokemonParalyzed,
					Pokemon: event.Pokemon,
					Trainer: event.Trainer,
				}
			}
		}
	}
}

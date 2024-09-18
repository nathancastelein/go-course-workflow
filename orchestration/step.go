package orchestration

import (
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type Step struct {
	Name string
	Do   func(trainer *pokemon.Trainer, pokemon *pokemon.Pokemon) error
}

func (w Step) LogValue() slog.Value {
	return slog.StringValue(w.Name)
}

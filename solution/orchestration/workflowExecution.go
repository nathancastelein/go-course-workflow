package orchestration

import (
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
)

type WorkflowExecution struct {
	ID             string
	WorkflowName   string
	PerformedSteps []string
	Trainer        *pokemon.Trainer
	Pokemon        *pokemon.Pokemon
}

func (w WorkflowExecution) LogValue() slog.Value {
	return slog.StringValue(w.ID)
}

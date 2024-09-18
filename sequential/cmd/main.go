package main

import (
	"context"
	"log/slog"

	"github.com/nathancastelein/go-course-workflows/pokemon"
	"github.com/nathancastelein/go-course-workflows/sequential"
)

var ()

func main() {
	catchProcess := sequential.NewCatchProcess(
		pokemon.NewStatusService(),
		pokemon.NewCombatService(),
		pokemon.NewPokeballService(),
	)

	err := catchProcess.Do(context.Background(), pokemon.Sacha(), pokemon.Mewtwo())
	if err != nil {
		slog.Error("fail to catch pokemon", slog.Any("error", err))
		return
	}
}
